package help

import (
	flagUtils "comeva/internal/utils/flag"
	"errors"
	"flag"
	"fmt"
	"strings"
)

var FlagSet *flag.FlagSet
var OptionalFlagSet = flag.NewFlagSet("", flag.ContinueOnError)
var OptionalBoolFlagSet = flag.NewFlagSet("", flag.ContinueOnError)

type FlagString string

var flagNames = struct {
	Path, CreateDocs FlagString
}{
	Path:       "path",
	CreateDocs: "create-docs",
}

type args struct {
	Path       []string
	CreateDocs string
}

func NewArgs() (arg *args, e error) {
	arg = &args{}
	arg.Path = pathFlag
	arg.CreateDocs = string(createDocsFlag)
	return
}

var createDocsFlag docsRoot = docsRoot("./docs/")

var pathFlag = make(path, 0)

func init() {
	OptionalBoolFlagSet.Var(&createDocsFlag, string(flagNames.CreateDocs),
		`Recreate the docs under the current directory in the passed `+"`path`"+` sub-directory.`,
	)

	OptionalFlagSet.Var(&pathFlag, string(flagNames.Path),
		`Path of help to show. Slash-separated with a slash at the beginning. Without a
slash at the end, get a file; with a slash at the end, get a directory.`)

	FlagSet = flagUtils.Combine(OptionalFlagSet, OptionalBoolFlagSet)
}

type docsRoot string

func (docs *docsRoot) String() string {
	return fmt.Sprintf(`"%v"`, string(*docs))
}

func (docs *docsRoot) Set(flg string) (e error) {
	if flg == "true" {
		return
	}
	*docs = docsRoot(flg)
	return
}

func (docs *docsRoot) IsBoolFlag() bool {
	return true
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
