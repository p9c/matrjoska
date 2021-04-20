// Package podopts is a configuration system to fit with the all-in-one philosophy guiding the design of the parallelcoin
// pod.
//
// The configuration is stored by each component of the connected applications, so all data is stored in concurrent-safe
// atomics, and there is a facility to invoke a function in response to a new value written into a field by other
// threads.
//
// There is a custom JSON marshal/unmarshal for each field type and for the whole configuration that only saves values
// that differ from the defaults, similar to 'omitempty' in struct tags but where 'empty' is the default value instead
// of the default zero created by Go's memory allocator. This enables easy compositing of multiple sources.
//
package config

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/btcsuite/go-socks/socks"

	"github.com/p9c/matrjoska/pkg/apputil"
	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pkg/constant"
	"github.com/p9c/opts/binary"
	"github.com/p9c/opts/cmds"
	"github.com/p9c/opts/duration"
	"github.com/p9c/opts/float"
	"github.com/p9c/opts/integer"
	"github.com/p9c/opts/list"
	"github.com/p9c/opts/opt"
	"github.com/p9c/opts/text"
)

// Call uses settings in the context to call the method with the given parameters and returns the raw json bytes
func Call(cx *Config, wallet bool, method string, params ...interface{}) (result []byte, e error) {
	// Ensure the specified method identifies a valid registered command and is one of the usable types.
	var usageFlags btcjson.UsageFlag
	usageFlags, e = btcjson.MethodUsageFlags(method)
	if e != nil {
		e = errors.New("Unrecognized command '" + method + "' : " + e.Error())
		// HelpPrint()
		return
	}
	if usageFlags&btcjson.UnusableFlags != 0 {
		E.F("The '%s' command can only be used via websockets\n", method)
		// HelpPrint()
		return
	}
	// Attempt to create the appropriate command using the arguments provided by the user.
	var cmd interface{}
	cmd, e = btcjson.NewCmd(method, params...)
	if e != nil {
		// Show the error along with its error code when it's a json. BTCJSONError as it realistically will always be
		// since the NewCmd function is only supposed to return errors of that type.
		if jerr, ok := e.(btcjson.GeneralError); ok {
			errText := fmt.Sprintf("%s command: %v (code: %s)\n", method, e, jerr.ErrorCode)
			e = errors.New(errText)
			// CommandUsage(method)
			return
		}
		// The error is not a json.BTCJSONError and this really should not happen. Nevertheless fall back to just
		// showing the error if it should happen due to a bug in the package.
		errText := fmt.Sprintf("%s command: %v\n", method, e)
		e = errors.New(errText)
		// CommandUsage(method)
		return
	}
	// Marshal the command into a JSON-RPC byte slice in preparation for sending it to the RPC server.
	var marshalledJSON []byte
	marshalledJSON, e = btcjson.MarshalCmd(1, cmd)
	if e != nil {
		return
	}
	// Send the JSON-RPC request to the server using the user-specified connection configuration.
	result, e = sendPostRequest(marshalledJSON, cx, wallet)
	if e != nil {
		return
	}
	return
}

// newHTTPClient returns a new HTTP client that is configured according to the proxy and TLS settings in the associated
// connection configuration.
func newHTTPClient(cfg *Config) (*http.Client, func(), error) {
	var dial func(ctx context.Context, network string, addr string) (net.Conn, error)
	ctx, cancel := context.WithCancel(context.Background())
	// Configure proxy if needed.
	if cfg.ProxyAddress.V() != "" {
		proxy := &socks.Proxy{
			Addr:     cfg.ProxyAddress.V(),
			Username: cfg.ProxyUser.V(),
			Password: cfg.ProxyPass.V(),
		}
		dial = func(_ context.Context, network string, addr string) (net.Conn, error) {
			c, e := proxy.Dial(network, addr)
			if e != nil {
				return nil, e
			}
			go func() {
			out:
				for {
					select {
					case <-ctx.Done():
						if e := c.Close(); E.Chk(e) {
						}
						break out
					}
				}
			}()
			return c, nil
		}
	}
	// Configure TLS if needed.
	var tlsConfig *tls.Config
	if cfg.ClientTLS.True() && cfg.RPCCert.V() != "" {
		pem, e := ioutil.ReadFile(cfg.RPCCert.V())
		if e != nil {
			cancel()
			return nil, nil, e
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pem)
		tlsConfig = &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: cfg.TLSSkipVerify.True(),
		}
	}
	// Create and return the new HTTP client potentially configured with a proxy and TLS.
	client := http.Client{
		Transport: &http.Transport{
			Proxy:                  nil,
			DialContext:            dial,
			TLSClientConfig:        tlsConfig,
			TLSHandshakeTimeout:    0,
			DisableKeepAlives:      false,
			DisableCompression:     false,
			MaxIdleConns:           0,
			MaxIdleConnsPerHost:    0,
			MaxConnsPerHost:        0,
			IdleConnTimeout:        0,
			ResponseHeaderTimeout:  0,
			ExpectContinueTimeout:  0,
			TLSNextProto:           nil,
			ProxyConnectHeader:     nil,
			MaxResponseHeaderBytes: 0,
			WriteBufferSize:        0,
			ReadBufferSize:         0,
			ForceAttemptHTTP2:      false,
		},
	}
	return &client, cancel, nil
}

// sendPostRequest sends the marshalled JSON-RPC command using HTTP-POST mode to the server described in the passed
// config struct. It also attempts to unmarshal the response as a JSON-RPC response and returns either the result field
// or the error field depending on whether or not there is an error.
func sendPostRequest(marshalledJSON []byte, cx *Config, wallet bool) ([]byte, error) {
	// Generate a request to the configured RPC server.
	protocol := "http"
	if cx.ClientTLS.True() {
		protocol = "https"
	}
	serverAddr := cx.RPCConnect.V()
	if wallet {
		serverAddr = cx.WalletServer.V()
		_, _ = fmt.Fprintln(os.Stderr, "using wallet server", serverAddr)
	}
	url := protocol + "://" + serverAddr
	bodyReader := bytes.NewReader(marshalledJSON)
	httpRequest, e := http.NewRequest("POST", url, bodyReader)
	if e != nil {
		return nil, e
	}
	httpRequest.Close = true
	httpRequest.Header.Set("Content-Type", "application/json")
	// Configure basic access authorization.
	httpRequest.SetBasicAuth(cx.Username.V(), cx.Password.V())
	T.Ln(cx.Username.V(), cx.Password.V())
	// Create the new HTTP client that is configured according to the user - specified options and submit the request.
	var httpClient *http.Client
	var cancel func()
	httpClient, cancel, e = newHTTPClient(cx)
	if e != nil {
		return nil, e
	}
	httpResponse, e := httpClient.Do(httpRequest)
	if e != nil {
		return nil, e
	}
	// close connection
	cancel()
	// Read the raw bytes and close the response.
	respBytes, e := ioutil.ReadAll(httpResponse.Body)
	if e := httpResponse.Body.Close(); E.Chk(e) {
		e = fmt.Errorf("error reading json reply: %v", e)
		return nil, e
	}
	// Handle unsuccessful HTTP responses
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		// Generate a standard error to return if the server body is empty. This should not happen very often, but it's
		// better than showing nothing in case the target server has a poor implementation.
		if len(respBytes) == 0 {
			return nil, fmt.Errorf("%d %s", httpResponse.StatusCode,
				http.StatusText(httpResponse.StatusCode),
			)
		}
		return nil, fmt.Errorf("%s", respBytes)
	}
	// Unmarshal the response.
	var resp btcjson.Response
	if e := json.Unmarshal(respBytes, &resp); E.Chk(e) {
		return nil, e
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

// ListCommands categorizes and lists all of the usable commands along with their one-line usage.
func ListCommands() (s string) {
	const (
		categoryChain uint8 = iota
		categoryWallet
		numCategories
	)
	// Get a list of registered commands and categorize and filter them.
	cmdMethods := btcjson.RegisteredCmdMethods()
	categorized := make([][]string, numCategories)
	for _, method := range cmdMethods {
		var e error
		var flags btcjson.UsageFlag
		if flags, e = btcjson.MethodUsageFlags(method); E.Chk(e) {
			continue
		}
		// Skip the commands that aren't usable from this utility.
		if flags&btcjson.UnusableFlags != 0 {
			continue
		}
		var usage string
		if usage, e = btcjson.MethodUsageText(method); E.Chk(e) {
			continue
		}
		// Categorize the command based on the usage flags.
		category := categoryChain
		if flags&btcjson.UFWalletOnly != 0 {
			category = categoryWallet
		}
		categorized[category] = append(categorized[category], usage)
	}
	// Display the command according to their categories.
	categoryTitles := make([]string, numCategories)
	categoryTitles[categoryChain] = "Chain Server Commands:"
	categoryTitles[categoryWallet] = "Wallet Server Commands (--wallet):"
	for category := uint8(0); category < numCategories; category++ {
		s += categoryTitles[category]
		s += "\n"
		for _, usage := range categorized[category] {
			s += "\t" + usage + "\n"
		}
		s += "\n"
	}
	return
}

// HelpPrint is the uninitialized help print function
var HelpPrint = func() {
	fmt.Println("help has not been overridden")
}

// CtlMain is the entry point for the pod.Ctl component
func CtlMain(cx *Config) {
	args := cx.ExtraArgs
	if len(args) < 1 {
		ListCommands()
		os.Exit(1)
	}
	// Ensure the specified method identifies a valid registered command and is one of the usable types.
	method := args[0]
	var usageFlags btcjson.UsageFlag
	var e error
	if usageFlags, e = btcjson.MethodUsageFlags(method); E.Chk(e) {
		_, _ = fmt.Fprintf(os.Stderr, "Unrecognized command '%s'\n", method)
		HelpPrint()
		os.Exit(1)
	}
	if usageFlags&btcjson.UnusableFlags != 0 {
		_, _ = fmt.Fprintf(os.Stderr, "The '%s' command can only be used via websockets\n", method)
		HelpPrint()
		os.Exit(1)
	}
	// Convert remaining command line args to a slice of interface values to be passed along as parameters to new
	// command creation function. Since some commands, such as submitblock, can involve data which is too large for the
	// Operating System to allow as a normal command line parameter, support using '-' as an argument to allow the
	// argument to be read from a stdin pipe.
	bio := bufio.NewReader(os.Stdin)
	params := make([]interface{}, 0, len(args[1:]))
	for _, arg := range args[1:] {
		if arg == "-" {
			var param string
			if param, e = bio.ReadString('\n'); E.Chk(e) && e != io.EOF {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to read data from stdin: %v\n", e)
				os.Exit(1)
			}
			if e == io.EOF && len(param) == 0 {
				_, _ = fmt.Fprintln(os.Stderr, "Not enough lines provided on stdin")
				os.Exit(1)
			}
			param = strings.TrimRight(param, "\r\n")
			params = append(params, param)
			continue
		}
		params = append(params, arg)
	}
	var result []byte
	if result, e = Call(cx, cx.UseWallet.True(), method, params...); E.Chk(e) {
		return
	}
	// // Attempt to create the appropriate command using the arguments provided by the user.
	// cmd, e := btcjson.NewCmd(method, params...)
	// if e != nil  {
	// 	E.Ln(e)
	// 	// Show the error along with its error code when it's a json. BTCJSONError as it realistically will always be
	// 	// since the NewCmd function is only supposed to return errors of that type.
	// 	if jerr, ok := err.(btcjson.BTCJSONError); ok {
	// 		fmt.Fprintf(os.Stderr, "%s command: %v (code: %s)\n", method, e, jerr.ErrorCode)
	// 		CommandUsage(method)
	// 		os.Exit(1)
	// 	}
	// 	// The error is not a json.BTCJSONError and this really should not happen. Nevertheless fall back to just
	// 	// showing the error if it should happen due to a bug in the package.
	// 	fmt.Fprintf(os.Stderr, "%s command: %v\n", method, e)
	// 	CommandUsage(method)
	// 	os.Exit(1)
	// }
	// // Marshal the command into a JSON-RPC byte slice in preparation for sending it to the RPC server.
	// marshalledJSON, e := btcjson.MarshalCmd(1, cmd)
	// if e != nil  {
	// 	E.Ln(e)
	// 	fmt.Println(e)
	// 	os.Exit(1)
	// }
	// // Send the JSON-RPC request to the server using the user-specified connection configuration.
	// result, e := sendPostRequest(marshalledJSON, cx)
	// if e != nil  {
	// 	E.Ln(e)
	// 	os.Exit(1)
	// }
	// Choose how to display the result based on its type.
	strResult := string(result)
	switch {
	case strings.HasPrefix(strResult, "{") || strings.HasPrefix(strResult, "["):
		var dst bytes.Buffer
		if e = json.Indent(&dst, result, "", "  "); E.Chk(e) {
			fmt.Printf("Failed to format result: %v", e)
			os.Exit(1)
		}
		fmt.Println(dst.String())
	case strings.HasPrefix(strResult, `"`):
		var str string
		if e = json.Unmarshal(result, &str); E.Chk(e) {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to unmarshal result: %v", e)
			os.Exit(1)
		}
		fmt.Println(str)
	case strResult != "null":
		fmt.Println(strResult)
	}
}

// CommandUsage display the usage for a specific command.
func CommandUsage(method string) {
	var usage string
	var e error
	if usage, e = btcjson.MethodUsageText(method); E.Chk(e) {
		// This should never happen since the method was already checked before calling this function, but be safe.
		fmt.Println("Failed to obtain command usage:", e)
		return
	}
	fmt.Println("Usage:")
	fmt.Printf("  %s\n", usage)
}

// Configs is the source location for the Config items, which is used to generate the Config struct
type Configs map[string]opt.Option
type ConfigSliceElement struct {
	Opt  opt.Option
	Name string
}
type ConfigSlice []ConfigSliceElement

func (c ConfigSlice) Len() int           { return len(c) }
func (c ConfigSlice) Less(i, j int) bool { return c[i].Name < c[j].Name }
func (c ConfigSlice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

// Initialize loads in configuration from disk and from environment on top of the default base
func (c *Config) Initialize() (e error) {
	// the several places configuration is sourced from are overlaid in the following order:
	// default -> config file -> environment variables -> commandline flags
	T.Ln("initializing configuration...")
	// first lint the configuration
	var aos map[string][]string
	if aos, e = getAllOptionStrings(c); E.Chk(e) {
		return
	}
	// this function will panic if there is potential for ambiguity in the commandline configuration args.
	T.Ln("linting configuration items")
	if _, e = findConflictingItems(aos); E.Chk(e) {
	}
	// generate and add the help commands to the help tree
	GetHelp(c)
	// process the commandline
	T.Ln("processing commandline arguments", os.Args[1:])
	var cm *cmds.Command
	var options []opt.Option
	var optVals []string
	if c.ExtraArgs, cm, options, optVals, e = c.processCommandlineArgs(os.Args[1:]); E.Chk(e) {
		return
	}

	if cm != nil {
		c.RunningCommand = *cm
		// I.S(c.RunningCommand)
		// } else {
		// 	c.RunningCommand = c.Commands[0]
	}
	// if the user sets the configfile directly, or the datadir on the commandline we need to load it from that path
	T.Ln("checking from where to load the configuration file")
	datadir := c.DataDir.V()
	var configPath string
	for i := range options {
		if options[i].Name() == "configfile" {
			if _, e = options[i].ReadInput(optVals[i]); E.Chk(e) {
				configPath = optVals[i]
			}
		}
		if options[i].Name() == "datadir" {
			I.Ln("datadir was set", optVals[i])
			if _, e = options[i].ReadInput(optVals[i]); !E.Chk(e) {
				datadir = options[i].Type().(*text.Opt).V()
				// I.Ln(datadir)
				// if _, e = c.DataDir.ReadInput(datadir); E.Chk(e) {
				// }
				// D.Ln(c.DataDir.V(), c.RPCKey.V(), c.RPCKey.Def, c.RPCCert.V(), c.RPCCert.Def, c.RPCKey.V(), c.RPCKey.Def)

				// reset all defaults that base on the datadir to apply hereafter
				// if the value is default, update it to the new datadir, and update the default field, otherwise assume
				// it has been set in the commandline args, and if it is different in environment or config file
				// it will be loaded in next, with the command line option value overriding at the end.
				if c.CAFile.V() == c.CAFile.Def {
					e = c.CAFile.Set(filepath.Join(datadir, "ca.cert"))
				}
				c.CAFile.Def = filepath.Join(datadir, "ca.cert")
				if c.ConfigFile.V() == c.ConfigFile.Def {
					e = c.ConfigFile.Set(filepath.Join(datadir, "pod.json"))
				}
				c.ConfigFile.Def = filepath.Join(datadir, constant.PodConfigFilename)
				if c.RPCKey.V() == c.RPCKey.Def {
					e = c.RPCKey.Set(filepath.Join(datadir, "rpc.key"))
				}
				c.RPCKey.Def = filepath.Join(datadir, "rpc.key")
				if c.RPCCert.V() == c.RPCCert.Def {
					e = c.RPCCert.Set(filepath.Join(datadir, "rpc.cert"))
				}
				c.RPCCert.Def = filepath.Join(datadir, "rpc.cert")
			}
		}
	}
	// I.Ln(c.RPCKey.V(), c.RPCKey.Def, c.RPCCert.V(), c.RPCCert.Def, c.RPCKey.V(), c.RPCKey.Def)
	// D.Ln(c.WalletFile.V(), c.WalletFile.Def, c.LogDir.V(), c.LogDir.Def)
	for i := range options {
		if options[i].Name() == "network" {
			I.Ln("network was set", optVals[i])
			if c.WalletFile.V() == c.WalletFile.Def {
				_, e = c.WalletFile.ReadInput(filepath.Join(datadir, optVals[i], "wallet.db"))
			}
			c.WalletFile.Def = filepath.Join(datadir, optVals[i], constant.DbName)
			if c.LogDir.V() == c.LogDir.Def {
				_, e = c.LogDir.ReadInput(filepath.Join(datadir, optVals[i]))
			}
			c.LogDir.Def = filepath.Join(datadir, optVals[i])
		}
	}
	// D.Ln(c.WalletFile.V(), c.WalletFile.Def, c.LogDir.V(), c.LogDir.Def)
	// load the configuration file into the config
	resolvedConfigPath := c.ConfigFile.V()
	if configPath != "" {
		T.Ln("loading config from", configPath)
		resolvedConfigPath = configPath
	} else {
		if datadir != "" {
			if strings.HasPrefix(datadir, "~") {
				var homeDir string
				var usr *user.User
				var e error
				if usr, e = user.Current(); e == nil {
					homeDir = usr.HomeDir
				}
				// Fall back to standard HOME environment variable that works for most POSIX OSes if the directory from the Go
				// standard lib failed.
				if e != nil || homeDir == "" {
					homeDir = os.Getenv("HOME")
				}

				datadir = strings.Replace(datadir, "~", homeDir, 1)
			}

			if resolvedConfigPath, e = filepath.Abs(filepath.Clean(filepath.Join(datadir, constant.PodConfigFilename,
			),
			),
			); E.Chk(e) {
			}
			T.Ln("loading config from", resolvedConfigPath)
		}
	}
	var configExists bool
	if e = c.loadConfig(resolvedConfigPath); !D.Chk(e) {
		configExists = true
	}
	// read the environment variables into the config
	if e = c.loadEnvironment(); D.Chk(e) {
	}
	// read in the commandline options over top as they have highest priority
	for i := range options {
		if _, e = options[i].ReadInput(optVals[i]); E.Chk(e) {
		}
	}
	if !configExists || c.Save.True() {
		c.Save.F()
		// save the configuration file
		var j []byte
		// c.ShowAll=true
		if j, e = json.MarshalIndent(c, "", "    "); !E.Chk(e) {
			I.F("saving config\n%s\n", string(j))
			apputil.EnsureDir(resolvedConfigPath)
			if e = ioutil.WriteFile(resolvedConfigPath, j, 0660); E.Chk(e) {
				panic(e)
			}
		}

	}
	return
}

// loadEnvironment scans the environment variables for values relevant to pod
func (c Config) loadEnvironment() (e error) {
	env := os.Environ()
	c.ForEach(func(o opt.Option) bool {
		varName := "POD_" + strings.ToUpper(o.Name())
		for i := range env {
			if strings.HasPrefix(env[i], varName) {
				envVal := strings.Split(env[i], varName)[1]
				if _, e = o.LoadInput(envVal); D.Chk(e) {
				}
			}
		}
		return true
	},
	)

	return
}

// loadConfig loads the config from a file and unmarshals it into the config
func (c *Config) loadConfig(path string) (e error) {
	e = fmt.Errorf("no config found at %s", path)
	var cf []byte
	if !apputil.FileExists(path) {
		return
	} else if cf, e = ioutil.ReadFile(path); !D.Chk(e) {
		if e = json.Unmarshal(cf, c); D.Chk(e) {
		}
	}
	return
}

// WriteToFile writes the current config to a file as json
func (c Config) WriteToFile(filename string) (e error) {
	// always load first and ensure written changes propagated or if one is manually running components independently
	if e = c.loadConfig(filename); E.Chk(e) {
		return
	}
	var j []byte
	if j, e = json.MarshalIndent(c, "", "  "); E.Chk(e) {
		return
	}
	if e = ioutil.WriteFile(filename, j, 0660); E.Chk(e) {
	}
	return
}

// ForEach iterates the options in defined order with a closure that takes an opt.Option
func (c Config) ForEach(fn func(ifc opt.Option) bool) bool {
	t := reflect.ValueOf(c)
	// t = t.Elem()
	for i := 0; i < t.NumField(); i++ {
		// asserting to an Option ensures we skip the ancillary fields
		if iff, ok := t.Field(i).Interface().(opt.Option); ok {
			if !fn(iff) {
				return false
			}
		}
	}
	return true
}

// GetOption searches for a match amongst the podopts
func (c Config) GetOption(input string) (op opt.Option, value string, e error) {
	T.Ln("checking arg for opt:", input)
	found := false
	if c.ForEach(func(ifc opt.Option) bool {
		aos := ifc.GetAllOptionStrings()
		for i := range aos {
			if strings.HasPrefix(input, aos[i]) {
				value = input[len(aos[i]):]
				found = true
				op = ifc
				return false
			}
		}
		return true
	},
	) {
	}
	if !found {
		e = fmt.Errorf("opt not found")
	}
	return
}

// MarshalJSON implements the json marshaller for the config. It only stores non-default values so can be composited.
func (c Config) MarshalJSON() (b []byte, e error) {
	outMap := make(map[string]interface{})
	c.ForEach(
		func(ifc opt.Option) bool {
			switch ii := ifc.(type) {
			case *binary.Opt:
				if ii.True() == ii.Def && ii.Data.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = ii.True()
			case *list.Opt:
				v := ii.S()
				if len(v) == len(ii.Def) && ii.Data.OmitEmpty && !c.ShowAll {
					foundMismatch := false
					for i := range v {
						if v[i] != ii.Def[i] {
							foundMismatch = true
							break
						}
					}
					if !foundMismatch {
						return true
					}
				}
				outMap[ii.Option] = v
			case *float.Opt:
				if ii.Value.Load() == ii.Def && ii.Data.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = ii.Value.Load()
			case *integer.Opt:
				if ii.Value.Load() == ii.Def && ii.Data.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = ii.Value.Load()
			case *text.Opt:
				v := string(ii.Value.Load().([]byte))
				// fmt.Printf("def: '%s'", v)
				// spew.Dump(ii.def)
				if v == ii.Def && ii.Data.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = v
			case *duration.Opt:
				if ii.Value.Load() == ii.Def && ii.Data.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = fmt.Sprint(ii.Value.Load())
			default:
			}
			return true
		},
	)
	return json.Marshal(&outMap)
}

// UnmarshalJSON implements the Unmarshaller interface so it only writes to fields with those non-default values set.
func (c Config) UnmarshalJSON(data []byte) (e error) {
	ifc := make(map[string]interface{})
	if e = json.Unmarshal(data, &ifc); E.Chk(e) {
		return
	}
	// I.S(ifc)
	c.ForEach(func(iii opt.Option) bool {
		switch ii := iii.(type) {
		case *binary.Opt:
			if i, ok := ifc[ii.Option]; ok {
				var ir bool
				if ir, ok = i.(bool); ir != ii.Def {
					// I.Ln(ii.Option+":", i.(binary), "default:", ii.def, "prev:", c.Map[ii.Option].(*Opt).True())
					ii.Set(ir)
				}
			}
		case *list.Opt:
			matched := true
			if d, ok := ifc[ii.Option]; ok {
				if ds, ok2 := d.([]interface{}); ok2 {
					for i := range ds {
						if len(ii.Def) >= len(ds) {
							if ds[i] != ii.Def[i] {
								matched = false
								break
							}
						} else {
							matched = false
						}
					}
					if matched {
						return true
					}
					// I.Ln(ii.Option+":", ds, "default:", ii.def, "prev:", c.Map[ii.Option].(*Opt).S())
					ii.Set(ifcToStrings(ds))
				}
			}
		case *float.Opt:
			if d, ok := ifc[ii.Option]; ok {
				// I.Ln(ii.Option+":", d.(float64), "default:", ii.def, "prev:", c.Map[ii.Option].(*Opt).V())
				ii.Set(d.(float64))
			}
		case *integer.Opt:
			if d, ok := ifc[ii.Option]; ok {
				// I.Ln(ii.Option+":", int64(d.(float64)), "default:", ii.def, "prev:", c.Map[ii.Option].(*Opt).V())
				ii.Set(int(d.(float64)))
			}
		case *text.Opt:
			if d, ok := ifc[ii.Option]; ok {
				if ds, ok2 := d.(string); ok2 {
					if ds != ii.Def {
						// I.Ln(ii.Option+":", d.(string), "default:", ii.def, "prev:", c.Map[ii.Option].(*Opt).V())
						ii.Set(d.(string))
					}
				}
			}
		case *duration.Opt:
			if d, ok := ifc[ii.Option]; ok {
				var parsed time.Duration
				parsed, e = time.ParseDuration(d.(string))
				// I.Ln(ii.Option+":", parsed, "default:", ii.Opt(), "prev:", c.Map[ii.Option].(*Opt).V())
				ii.Set(parsed)
			}
		default:
		}
		return true
	},
	)
	return
}

func (c Config) processCommandlineArgs(args []string) (
	remArgs []string, cm *cmds.Command, op []opt.Option,
	optVals []string, e error,
) {
	// I.S(c.Commands)
	// I.S(args)
	// first we will locate all the commands specified to mark the 3 sections, opt, commands, and the remainder is
	// arbitrary for the node
	commands := make(map[int]cmds.Command)
	var commandsStart, commandsEnd int = -1, -1
	var found bool
	for i := range args {
		T.Ln("checking for commands:", args[i], commandsStart, commandsEnd, "current arg index:", i)
		var depth, dist int
		if found, depth, dist, cm, e = c.Commands.Find(args[i], depth, dist); E.Chk(e) {
			continue
		}
		if cm != nil {
			if cm.Parent != nil {
				if cm.Parent.Name == "help" {
					// I.S(commands)
					if commands[0].Name != "help" {
						found = false
					}
				}
			}
		}
		if found {
			if commandsStart == -1 {
				commandsStart = i
				commandsEnd = i + 1
			}
			if oc, ok := commands[depth]; ok {
				e = fmt.Errorf("second command found at same depth '%s' and '%s'", oc.Name, cm.Name)
				return
			}
			commandsEnd = i + 1
			T.Ln("commandStart", commandsStart, commandsEnd, args[commandsStart:commandsEnd])
			T.Ln("found command", cm.Name, "argument number", i, "at depth", depth, "distance", dist)
			commands[depth] = *cm
		} else {
			// commandsStart=i+1
			// commandsEnd=i+1
			// T.Ln("not found:", args[i], "commandStart", commandsStart, commandsEnd, args[commandsStart:commandsEnd])
			T.Ln("argument", args[i], "is not a command", commandsStart, commandsEnd)
		}
	}
	// I.S(commands, cm)
	// commandsEnd++
	cmds := []int{}
	if len(commands) == 0 {
		I.Ln("setting default command")
		commands[0] = c.Commands[0]
	} else {
		T.Ln("checking commands")
		// I.S(commands)
		for i := range commands {
			cmds = append(cmds, i)
		}
		// I.S(cmds)
		if len(cmds) > 0 {
			sort.Ints(cmds)
			var cms []string
			for i := range commands {
				cms = append(cms, commands[i].Name)
			}
			if cmds[0] != 1 {
				e = fmt.Errorf("commands must include base level item for disambiguation %v", cms)
			}
			prev := cmds[0]
			for i := range cmds {
				if i == 0 {
					continue
				}
				if cmds[i] != prev+1 {
					e = fmt.Errorf("more than one command specified, %v", cms)
					return
				}
				found = false
				for j := range commands[cmds[i-1]].Commands {
					if commands[cmds[i]].Name == commands[cmds[i-1]].Commands[j].Name {
						found = true
					}
				}
				if !found {
					e = fmt.Errorf("multiple commands are not a path on the command tree %v", cms)
					return
				}
			}
		}
		T.Ln("commands:", commandsStart, commandsEnd, args[commandsStart:commandsEnd])
		T.Ln("length of gathered commands", len(commands))
		if len(commands) == 1 {
			for _, x := range commands {
				cm = &x
			}
		}
	}
	// if there was no command the commands start and end after all the args
	if commandsStart < 0 || commandsEnd < 0 {
		commandsStart = len(args)
		commandsEnd = commandsStart
	}
	T.Ln("commands section:", commandsStart, commandsEnd)
	if commandsStart > 0 {
		T.Ln("opt found", args[:commandsStart])
		// we have opt to check
		for i := range args {
			// if i == 0 {
			// 	continue
			// }
			if i >= commandsStart {
				break
			}
			var val string
			var o opt.Option
			if o, val, e = c.GetOption(args[i]); E.Chk(e) {
				e = fmt.Errorf("argument %d: '%s' lacks a valid opt prefix", i, args[i])
				return
			}
			if _, e = o.ReadInput(val); E.Chk(e) {
				return
			}
			T.Ln("found opt:", o.String())
			op = append(op, o)
			optVals = append(optVals, val)
		}
	}
	T.S(op)
	if len(cmds) < 1 {
		cmds = []int{0}
		commands[0] = c.Commands[0]
	}
	if commandsEnd > 0 && len(args) > commandsEnd {
		remArgs = args[commandsEnd:]
	}
	D.F("args that will pass to command: %v", remArgs)
	// I.S(commands, cm)
	// I.S(commands[cmds[len(cmds)-1]], op, args[commandsEnd:])
	return
}

// ReadCAFile reads in the configured Certificate Authority for TLS connections
func (c Config) ReadCAFile() []byte {
	// Read certificate file if TLS is not disabled.
	var certs []byte
	if c.ClientTLS.True() {
		var e error
		if certs, e = ioutil.ReadFile(c.CAFile.V()); E.Chk(e) {
			// If there's an error reading the CA file, continue with nil certs and without the client connection.
			certs = nil
		}
	} else {
		I.Ln("chain server RPC TLS is disabled")
	}
	return certs
}

func ifcToStrings(ifc []interface{}) (o []string) {
	for i := range ifc {
		o = append(o, ifc[i].(string))
	}
	return
}

type details struct {
	name, option, desc, def string
	aliases                 []string
	documentation           string
}

// GetHelp walks the command tree and gathers the options and creates a set of help functions for all commands and
// options in the set
func GetHelp(c *Config) {
	cm := cmds.Command{
		Name:        "help",
		Description: "prints information about how to use pod",
		Entrypoint:  HelpFunction,
		Commands:    nil,
	}
	// first add all the options
	c.ForEach(func(ifc opt.Option) bool {
		o := fmt.Sprintf("Parallelcoin Pod All-in-One Suite\n\n")
		var dt details
		switch ii := ifc.(type) {
		case *binary.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description, fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		case *list.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description, fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		case *float.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description, fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		case *integer.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description, fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		case *text.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description, fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		case *duration.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description, fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		}
		cm.Commands = append(cm.Commands, cmds.Command{
			Name:        dt.option,
			Description: dt.desc,
			Entrypoint: func(ifc interface{}) (e error) {
				o += fmt.Sprintf("Help information about %s\n\n\toption name:\n\t\t%s\n\taliases:\n\t\t%s\n\tdescription:\n\t\t%s\n\tdefault:\n\t\t%v\n",
					dt.name, dt.option, dt.aliases, dt.desc, dt.def,
				)
				if dt.documentation != "" {
					o += "\tdocumentation:\n\t\t" + dt.documentation + "\n\n"
				}
				fmt.Fprint(os.Stderr, o)
				return
			},
			Commands: nil,
			Parent:   &cm,
		},
		)
		// for i := range
		return true
	},
	)
	// // next add all the commands
	// c.Commands.ForEach(func(cm cmds.Command) bool {
	// 	return true
	// }, 0, 0,
	// )
	c.Commands = append(c.Commands, cm)
	return
}

func HelpFunction(ifc interface{}) error {
	c := assertToPodState(ifc)
	var o string
	o += fmt.Sprintf("Parallelcoin Pod All-in-One Suite\n\n")
	o += fmt.Sprintf("Usage:\n\t%s [options] [commands] [command parameters]\n\n", os.Args[0])
	o += fmt.Sprintf("Commands:\n")
	for i := range c.Commands {
		oo := fmt.Sprintf("\t%s", c.Commands[i].Name)
		nrunes := utf8.RuneCountInString(oo)
		o += oo + fmt.Sprintf(strings.Repeat(" ", 9-nrunes)+"%s\n", c.Commands[i].Description)
	}
	o += fmt.Sprintf(
		"\nOptions:\n\tset values on options concatenated against the option keyword or separated with '='\n",
	)
	o += fmt.Sprintf("\teg: addcheckpoints=deadbeefcafe,someothercheckpoint AP127.0.0.1:11047\n")
	o += fmt.Sprintf("\tfor items that take multiple string values, you can repeat the option with further\n")
	o += fmt.Sprintf("\tinstances of the option or separate the items with (only) commas as the above example\n\n")
	// items := make(map[string][]opt.Option)
	descs := make(map[string]string)
	c.ForEach(func(ifc opt.Option) bool {
		meta := ifc.GetMetadata()
		oo := fmt.Sprintf("\t%s %v", meta.Option, meta.Aliases)
		nrunes := utf8.RuneCountInString(oo)
		var def string
		switch ii := ifc.(type) {
		case *binary.Opt:
			def = fmt.Sprint(ii.Def)
		case *list.Opt:
			def = fmt.Sprint(ii.Def)
		case *float.Opt:
			def = fmt.Sprint(ii.Def)
		case *integer.Opt:
			def = fmt.Sprint(ii.Def)
		case *text.Opt:
			def = fmt.Sprint(ii.Def)
		case *duration.Opt:
			def = fmt.Sprint(ii.Def)
		}
		descs[meta.Group] += oo + fmt.Sprintf(strings.Repeat(" ", 32-nrunes)+"%s, default: %s\n", meta.Description, def)
		return true
	},
	)
	var cats []string
	for i := range descs {
		cats = append(cats, i)
	}
	// I.S(cats)
	sort.Strings(cats)
	for i := range cats {
		if cats[i] != "" {
			o += "\n" + cats[i] + "\n"
		}
		o += descs[cats[i]]
	}
	// for i := range cats {
	// }
	o += fmt.Sprintf("\nadd the name of the command or option after 'help' or append it after "+
		"'help' in the commandline to get more detail - eg: %s help upnp\n\n", os.Args[0],
	)
	fmt.Fprintf(os.Stderr, o)
	return nil
}

func assertToPodState(ifc interface{}) (c *Config) {
	var ok bool
	if c, ok = ifc.(*Config); !ok {
		panic("wth")
	}
	return
}

func getAllOptionStrings(c *Config) (s map[string][]string, e error) {
	s = make(map[string][]string)
	if c.ForEach(func(ifc opt.Option) bool {
		md := ifc.GetMetadata()
		if _, ok := s[ifc.Name()]; ok {
			e = fmt.Errorf("conflicting opt names: %v %v", ifc.GetAllOptionStrings(), s[ifc.Name()])
			return false
		}
		s[ifc.Name()] = md.GetAllOptionStrings()
		return true
	},
	) {
	}
	// s["commandslist"] = c.Commands.GetAllCommands()
	return
}

func findConflictingItems(valOpts map[string][]string) (o []string, e error) {
	var ss, ls string
	for i := range valOpts {
		for j := range valOpts {
			if i == j {
				continue
			}
			a := valOpts[i]
			b := valOpts[j]
			for ii := range a {
				for jj := range b {
					ss, ls = shortestString(a[ii], b[jj])
					if ss == ls[:len(ss)] {
						E.F("conflict between %s and %s, ", ss, ls)
						o = append(o, ss, ls)
					}
				}
			}
		}
	}
	if len(o) > 0 {
		panic(fmt.Sprintf("conflicts found: %v", o))
	}
	return
}

func shortestString(a, b string) (s, l string) {
	switch {
	case len(a) > len(b):
		s, l = b, a
	default:
		s, l = a, b
	}
	return
}
