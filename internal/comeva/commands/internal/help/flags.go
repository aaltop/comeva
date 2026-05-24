package help

import (
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
func (p *path) Set(flg string) (e error) {
	*p = strings.Split(flg, "/")

	return
}
