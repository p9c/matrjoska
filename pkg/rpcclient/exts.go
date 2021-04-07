package rpcclient

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	js "encoding/json"
	"fmt"
	"github.com/p9c/monorepo/pkg/btcaddr"
	
	"github.com/p9c/monorepo/pkg/btcjson"
	"github.com/p9c/monorepo/pkg/chainhash"
	"github.com/p9c/monorepo/pkg/wire"
)

// FutureDebugLevelResult is a future promise to deliver the result of a DebugLevelAsync RPC invocation (or an
// applicable error).
type FutureDebugLevelResult chan *response

// Receive waits for the response promised by the future and returns the result of setting the debug logging level to
// the passed level specification or the list of of the available subsystems for the special keyword 'show'.
func (r FutureDebugLevelResult) Receive() (string, error) {
	res, e := receiveFuture(r)
	if e != nil {
		return "", e
	}
	// Unmashal the result as a string.
	var result string
	e = js.Unmarshal(res, &result)
	if e != nil {
		return "", e
	}
	return result, nil
}

// DebugLevelAsync returns an instance of a type that can be used to get the result of the RPC at some future time by
// invoking the Receive function on the returned instance. See DebugLevel for the blocking version and more details.
// NOTE: This is a pod extension.
func (c *Client) DebugLevelAsync(levelSpec string) FutureDebugLevelResult {
	cmd := btcjson.NewDebugLevelCmd(levelSpec)
	return c.sendCmd(cmd)
}

// DebugLevel dynamically sets the debug logging level to the passed level specification. The levelspec can be either a
// debug level or of the form:
//
// 	<subsystem>=<level>,<subsystem2>=<level2>,...
//
// Additionally, the special keyword 'show' can be used to get a list of the available subsystems.
//
// NOTE: This is a pod extension.
func (c *Client) DebugLevel(levelSpec string) (string, error) {
	return c.DebugLevelAsync(levelSpec).Receive()
}

// FutureCreateEncryptedWalletResult is a future promise to deliver the error result of a CreateEncryptedWalletAsync RPC
// invocation.
type FutureCreateEncryptedWalletResult chan *response

// Receive waits for and returns the error response promised by the future.
func (r FutureCreateEncryptedWalletResult) Receive() (e error) {
	_, e = receiveFuture(r)
	return e
}

// CreateEncryptedWalletAsync returns an instance of a type that can be used to get the result of the RPC at some future
// time by invoking the Receive function on the returned instance. See CreateEncryptedWallet for the blocking version
// and more details. NOTE: This is a btcwallet extension.
func (c *Client) CreateEncryptedWalletAsync(passphrase string) FutureCreateEncryptedWalletResult {
	cmd := btcjson.NewCreateEncryptedWalletCmd(passphrase)
	return c.sendCmd(cmd)
}

// CreateEncryptedWallet requests the creation of an encrypted wallet. Wallets managed by btcwallet are only written to
// disk with encrypted private keys, and generating wallets on the fly is impossible as it requires user input for the
// encryption passphrase.
//
// This RPC specifies the passphrase and instructs the wallet creation. This may error if a
// wallet is already opened, or the new wallet cannot be written to disk. NOTE: This is a btcwallet extension.
func (c *Client) CreateEncryptedWallet(passphrase string) (e error) {
	return c.CreateEncryptedWalletAsync(passphrase).Receive()
}

// FutureListAddressTransactionsResult is a future promise to deliver the result of a ListAddressTransactionsAsync RPC
// invocation (or an applicable error).
type FutureListAddressTransactionsResult chan *response

// Receive waits for the response promised by the future and returns information about all transactions associated with
// the provided addresses.
func (r FutureListAddressTransactionsResult) Receive() ([]btcjson.ListTransactionsResult, error) {
	res, e := receiveFuture(r)
	if e != nil {
		return nil, e
	}
	// Unmarshal the result as an array of listtransactions objects.
	var transactions []btcjson.ListTransactionsResult
	e = js.Unmarshal(res, &transactions)
	if e != nil {
		return nil, e
	}
	return transactions, nil
}

// ListAddressTransactionsAsync returns an instance of a type that can be used get the result of the RPC at some future
// time by invoking the Receive function on the returned instance. See ListAddressTransactions for the blocking version
// and more details. NOTE: This is a pod extension.
func (c *Client) ListAddressTransactionsAsync(
	addresses []btcaddr.Address,
	account string,
) FutureListAddressTransactionsResult {
	// Convert addresses to strings.
	addrs := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		addrs = append(addrs, addr.EncodeAddress())
	}
	cmd := btcjson.NewListAddressTransactionsCmd(addrs, &account)
	return c.sendCmd(cmd)
}

// ListAddressTransactions returns information about all transactions associated with the provided addresses. NOTE: This
// is a btcwallet extension.
func (c *Client) ListAddressTransactions(addresses []btcaddr.Address, account string) (
	[]btcjson.ListTransactionsResult,
	error,
) {
	return c.ListAddressTransactionsAsync(addresses, account).Receive()
}

// FutureGetBestBlockResult is a future promise to deliver the result of a GetBestBlockAsync RPC invocation (or an
// applicable error).
type FutureGetBestBlockResult chan *response

// Receive waits for the response promised by the future and returns the hash and height of the block in the longest
// (best) chain.
func (r FutureGetBestBlockResult) Receive() (*chainhash.Hash, int32, error) {
	res, e := receiveFuture(r)
	if e != nil {
		return nil, 0, e
	}
	// Unmarshal result as a getbestblock result object.
	var bestBlock btcjson.GetBestBlockResult
	e = js.Unmarshal(res, &bestBlock)
	if e != nil {
		return nil, 0, e
	}
	// Convert to hash from string.
	hash, e := chainhash.NewHashFromStr(bestBlock.Hash)
	if e != nil {
		return nil, 0, e
	}
	return hash, bestBlock.Height, nil
}

// GetBestBlockAsync returns an instance of a type that can be used to get the result of the RPC at some future time by
// invoking the Receive function on the returned instance. See GetBestBlock for the blocking version and more details.
//
// NOTE: This is a pod extension.
func (c *Client) GetBestBlockAsync() FutureGetBestBlockResult {
	cmd := btcjson.NewGetBestBlockCmd()
	return c.sendCmd(cmd)
}

// GetBestBlock returns the hash and height of the block in the longest (best) chain.
//
// NOTE: This is a pod extension.
func (c *Client) GetBestBlock() (*chainhash.Hash, int32, error) {
	return c.GetBestBlockAsync().Receive()
}

// FutureGetCurrentNetResult is a future promise to deliver the result of a GetCurrentNetAsync RPC invocation (or an
// applicable error).
type FutureGetCurrentNetResult chan *response

// Receive waits for the response promised by the future and returns the network the server is running on.
func (r FutureGetCurrentNetResult) Receive() (wire.BitcoinNet, error) {
	res, e := receiveFuture(r)
	if e != nil {
		return 0, e
	}
	// Unmarshal result as an int64.
	var net int64
	e = js.Unmarshal(res, &net)
	if e != nil {
		return 0, e
	}
	return wire.BitcoinNet(net), nil
}

// GetCurrentNetAsync returns an instance of a type that can be used to get the result of the RPC at some future time by
// invoking the Receive function on the returned instance. See GetCurrentNet for the blocking version and more details.
//
// NOTE: This is a pod extension.
func (c *Client) GetCurrentNetAsync() FutureGetCurrentNetResult {
	cmd := btcjson.NewGetCurrentNetCmd()
	return c.sendCmd(cmd)
}

// GetCurrentNet returns the network the server is running on.
//
// NOTE: This is a pod extension.
func (c *Client) GetCurrentNet() (wire.BitcoinNet, error) {
	return c.GetCurrentNetAsync().Receive()
}

// FutureGetHeadersResult is a future promise to deliver the result of a getheaders RPC invocation (or an applicable
// error).
//
// NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
type FutureGetHeadersResult chan *response

// Receive waits for the response promised by the future and returns the getheaders result.
//
// NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (r FutureGetHeadersResult) Receive() ([]wire.BlockHeader, error) {
	res, e := receiveFuture(r)
	if e != nil {
		return nil, e
	}
	// Unmarshal result as a slice of strings.
	var result []string
	e = js.Unmarshal(res, &result)
	if e != nil {
		return nil, e
	}
	// Deserialize the []string into []wire.BlockHeader.
	headers := make([]wire.BlockHeader, len(result))
	for i, headerHex := range result {
		serialized, e := hex.DecodeString(headerHex)
		if e != nil {
			return nil, e
		}
		e = headers[i].Deserialize(bytes.NewReader(serialized))
		if e != nil {
			return nil, e
		}
	}
	return headers, nil
}

// GetHeadersAsync returns an instance of a type that can be used to get the result of the RPC at some future time by
// invoking the Receive function on the returned instance. See GetHeaders for the blocking version and more details.
//
// NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (c *Client) GetHeadersAsync(blockLocators []chainhash.Hash, hashStop *chainhash.Hash) FutureGetHeadersResult {
	locators := make([]string, len(blockLocators))
	for i := range blockLocators {
		locators[i] = blockLocators[i].String()
	}
	hash := ""
	if hashStop != nil {
		hash = hashStop.String()
	}
	cmd := btcjson.NewGetHeadersCmd(locators, hash)
	return c.sendCmd(cmd)
}

// GetHeaders mimics the wire protocol getheaders and headers messages by returning all headers on the main chain after
// the first known block in the locators, up until a block hash matches hashStop.
//
// NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (c *Client) GetHeaders(blockLocators []chainhash.Hash, hashStop *chainhash.Hash) ([]wire.BlockHeader, error) {
	return c.GetHeadersAsync(blockLocators, hashStop).Receive()
}

// FutureExportWatchingWalletResult is a future promise to deliver the result of an ExportWatchingWalletAsync RPC
// invocation (or an applicable error).
type FutureExportWatchingWalletResult chan *response

// Receive waits for the response promised by the future and returns the exported wallet.
func (r FutureExportWatchingWalletResult) Receive() ([]byte, []byte, error) {
	res, e := receiveFuture(r)
	if e != nil {
		return nil, nil, e
	}
	// Unmarshal result as a JSON object.
	var obj map[string]interface{}
	e = js.Unmarshal(res, &obj)
	if e != nil {
		return nil, nil, e
	}
	// Chk for the wallet and tx string fields in the object.
	base64Wallet, ok := obj["wallet"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"unexpected response type for exportwatchingwallet 'wallet' field: %T\n",
			obj["wallet"],
		)
	}
	base64TxStore, ok := obj["tx"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"unexpected response type for exportwatchingwallet 'tx' field: %T\n",
			obj["tx"],
		)
	}
	walletBytes, e := base64.StdEncoding.DecodeString(base64Wallet)
	if e != nil {
		return nil, nil, e
	}
	txStoreBytes, e := base64.StdEncoding.DecodeString(base64TxStore)
	if e != nil {
		return nil, nil, e
	}
	return walletBytes, txStoreBytes, nil
}

// ExportWatchingWalletAsync returns an instance of a type that can be used to get the result of the RPC at some future
// time by invoking the Receive function on the returned instance. See ExportWatchingWallet for the blocking version and
// more details.
//
// NOTE: This is a btcwallet extension.
func (c *Client) ExportWatchingWalletAsync(account string) FutureExportWatchingWalletResult {
	cmd := btcjson.NewExportWatchingWalletCmd(&account, btcjson.Bool(true))
	return c.sendCmd(cmd)
}

// ExportWatchingWallet returns the raw bytes for a watching-only version of wallet.bin and tx.bin, respectively, for
// the specified account that can be used by btcwallet to enable a wallet which does not have the private keys necessary
// to spend funds.
//
// NOTE: This is a btcwallet extension.
func (c *Client) ExportWatchingWallet(account string) ([]byte, []byte, error) {
	return c.ExportWatchingWalletAsync(account).Receive()
}

// FutureSessionResult is a future promise to deliver the result of a SessionAsync RPC invocation (or an applicable
// error).
type FutureSessionResult chan *response

// Receive waits for the response promised by the future and returns the session result.
func (r FutureSessionResult) Receive() (*btcjson.SessionResult, error) {
	res, e := receiveFuture(r)
	if e != nil {
		return nil, e
	}
	// Unmarshal result as a session result object.
	var session btcjson.SessionResult
	e = js.Unmarshal(res, &session)
	if e != nil {
		return nil, e
	}
	return &session, nil
}

// SessionAsync returns an instance of a type that can be used to get the result of the RPC at some future time by
// invoking the Receive function on the returned instance. See Session for the blocking version and more details.
//
// NOTE: This is a btcsuite extension.
func (c *Client) SessionAsync() FutureSessionResult {
	// Not supported in HTTP POST mode.
	if c.config.HTTPPostMode {
		return newFutureError(ErrWebsocketsRequired)
	}
	cmd := btcjson.NewSessionCmd()
	return c.sendCmd(cmd)
}

// Session returns details regarding a websocket client's current connection. This RPC requires the client to be running
// in websocket mode.
//
// NOTE: This is a btcsuite extension.
func (c *Client) Session() (*btcjson.SessionResult, error) {
	return c.SessionAsync().Receive()
}

// FutureVersionResult is a future promise to deliver the result of a version RPC invocation (or an applicable error).
//
// NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
type FutureVersionResult chan *response

// Receive waits for the response promised by the future and returns the version result.
//
// NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (r FutureVersionResult) Receive() (
	map[string]btcjson.VersionResult,
	error,
) {
	res, e := receiveFuture(r)
	if e != nil {
		return nil, e
	}
	// Unmarshal result as a version result object.
	var vr map[string]btcjson.VersionResult
	e = js.Unmarshal(res, &vr)
	if e != nil {
		return nil, e
	}
	return vr, nil
}

// VersionAsync returns an instance of a type that can be used to get the result of the RPC at some future time by
// invoking the Receive function on the returned instance. See Version for the blocking version and more details.
//
// NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (c *Client) VersionAsync() FutureVersionResult {
	cmd := btcjson.NewVersionCmd()
	return c.sendCmd(cmd)
}

// Version returns information about the server's JSON-RPC API versions.
//
// NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (c *Client) Version() (map[string]btcjson.VersionResult, error) {
	return c.VersionAsync().Receive()
}