//+build tools

package libsql

// Tools written in Go can be installed using the go-get command.
//
// This technique is described in more detail in:
// https://github.com/go-modules-by-example/index/tree/23a56e1095864bf596f2ce3aec296ecc89ab71b9/010_tools
// and https://github.com/golang/go/issues/25922#issuecomment-451123151.
import (
	_ "github.com/gojuno/minimock/v3"
)
