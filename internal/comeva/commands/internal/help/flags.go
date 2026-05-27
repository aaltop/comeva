package help

import (
	"errors"
	"flag"
	"strings"
)

var FlagSet = flag.NewFlagSet("", flag.ContinueOnError)

type FlagString string

var flagNames = struct {
	Path FlagString
}{
	Path: "path",
}

type args struct {
	Path []string
}

func NewArgs() (arg *args, e error) {
	arg = &args{}
	arg.Path = pathFlag
	return
}

var pathFlag = make(path, 0)

func init() {
	FlagSet.Var(&pathFlag, string(flagNames.Path), "Path of help to show. Slash-separated without start or end slash.")
}

// path represents a slash-separated path.
type path []string

func (p *path) String() string {
	return strings.Join(*p, "/")
}

// Set takes a flag that is assumed to be a path starting with a slash,
// with slash separated parts of the path, and potentially ending with a slash.
// It sets in `p` the path split by the forward slash. Note that this leaves
// an empty string at the start and end if these have a slash in the flag,
// and for a single slash as the flag value, the result is a slice of two empty
// strings.
func (p *path) Set(flg string) (e error) {

	if !(len(flg) > 0 && flg[0] == '/') {
		return errors.New("Path should start with a forward slash (/).")
	}

	*p = strings.Split(flg, "/")
	return
}
