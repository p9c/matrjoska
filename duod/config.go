package main

import (
	"fmt"
	"strconv"
	"strings"
	
	"github.com/p9c/duod/pkg/chaincfg"
	"github.com/p9c/duod/pkg/chainhash"
	"github.com/p9c/duod/pkg/database"
	
	// This ensures the database drivers get registered
	_ "github.com/p9c/duod/pkg/database/ffldb"
)

var (
	// defaultHomeDir is the default home directory location (
	// this should be centralised)
	
	// KnownDbTypes stores the currently supported database drivers
	KnownDbTypes = database.SupportedDrivers()
	// runServiceCommand is only set to a real function on Windows.
	// It is used to parse and execute service commands specified via the -s flag.
	runServiceCommand func(string) error
)

// newCheckpointFromStr parses checkpoints in the '<height>:<hash>' format.
func newCheckpointFromStr(checkpoint string) (chaincfg.Checkpoint, error) {
	parts := strings.Split(checkpoint, ":")
	if len(parts) != 2 {
		return chaincfg.Checkpoint{}, fmt.Errorf(
			"unable to parse "+
				"checkpoint %q -- use the syntax <height>:<hash>",
			checkpoint,
		)
	}
	height, e := strconv.ParseInt(parts[0], 10, 32)
	if e != nil {
		return chaincfg.Checkpoint{}, fmt.Errorf(
			"unable to parse "+
				"checkpoint %q due to malformed height", checkpoint,
		)
	}
	if len(parts[1]) == 0 {
		return chaincfg.Checkpoint{}, fmt.Errorf(
			"unable to parse "+
				"checkpoint %q due to missing hash", checkpoint,
		)
	}
	hash, e := chainhash.NewHashFromStr(parts[1])
	if e != nil {
		return chaincfg.Checkpoint{}, fmt.Errorf(
			"unable to parse "+
				"checkpoint %q due to malformed hash", checkpoint,
		)
	}
	return chaincfg.Checkpoint{
			Height: int32(height),
			Hash:   hash,
		},
		nil
}

// ParseCheckpoints checks the checkpoint strings for valid syntax ( '<height>:<hash>') and parses them to
// chaincfg.Checkpoint instances.
func ParseCheckpoints(checkpointStrings []string) ([]chaincfg.Checkpoint, error) {
	if len(checkpointStrings) == 0 {
		return nil, nil
	}
	checkpoints := make([]chaincfg.Checkpoint, len(checkpointStrings))
	for i, cpString := range checkpointStrings {
		checkpoint, e := newCheckpointFromStr(cpString)
		if e != nil {
			return nil, e
		}
		checkpoints[i] = checkpoint
	}
	return checkpoints, nil
}

// ValidDbType returns whether or not dbType is a supported database type.
func ValidDbType(dbType string) bool {
	for _, knownType := range KnownDbTypes {
		if dbType == knownType {
			return true
		}
	}
	return false
}

// validLogLevel returns whether or not logLevel is a valid debug log level.
func validLogLevel(logLevel string) bool {
	switch logLevel {
	case "trace":
		fallthrough
	case "debug":
		fallthrough
	case "info":
		fallthrough
	case "warn":
		fallthrough
	case "error":
		fallthrough
	case "critical":
		return true
	}
	return false
}
