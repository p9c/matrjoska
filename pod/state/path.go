package state

import (
	"os"
	"path/filepath"

	"github.com/p9c/matrjoska/pod/podcfgs"
	"github.com/p9c/matrjoska/version"
)

// BlockDb returns the path to the block database given a database type.
func BlockDb(cx *State, dbType string, namePrefix string) string {
	// The database name is based on the database type.
	dbName := namePrefix + "_" + dbType
	if dbType == "sqlite" {
		dbName += ".db"
	}
	dbPath := filepath.Join(
		filepath.Join(
			cx.Config.DataDir.V(),
			cx.ActiveNet.Name,
		), dbName,
	)
	return dbPath
}


func Main() {
	// log.SetLogLevel("debug")
	T.Ln(version.Get())
	var cx *State
	var e error
	if cx, e = GetNew(podcfgs.GetDefaultConfig()); E.Chk(e) {
		fail()
	}

	// if e = debugConfig(cx.Config); E.Chk(e) {
	// 	fail()
	// }

	D.Ln("running command", cx.Config.RunningCommand.Name)
	if e = cx.Config.RunningCommand.Entrypoint(cx); E.Chk(e) {
		fail()
	}
}

// func debugConfig(c *podopts.Config) (e error) {
// 	c.ShowAll = true
// 	defer func() { c.ShowAll = false }()
// 	var j []byte
// 	if j, e = c.MarshalJSON(); E.Chk(e) {
// 		return
// 	}
// 	var b []byte
// 	jj := bytes.NewBuffer(b)
// 	if e = json.Indent(jj, j, "", "\t"); E.Chk(e) {
// 		return
// 	}
// 	I.Ln("\n"+jj.String())
// 	return
// }

func fail() {
	os.Exit(1)
}
