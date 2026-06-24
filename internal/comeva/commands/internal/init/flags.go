package init

import "flag"

var FlagSet = flag.NewFlagSet("", flag.ContinueOnError)

type FlagString string

var flagNames = struct {
	Global FlagString
}{
	Global: "global",
}

var global = FlagSet.Bool(string(flagNames.Global), false,
	"Initialise the configuration globally instead, e.g. in $HOME/.config/comeva/ on Linux.",
)

type args struct {
	Global bool
}

func NewArgs() (arg *args) {
	arg = &args{
		Global: *global,
	}
	return
}

type passedArgs struct {
	Global bool
}

func getPassedArgs(passedMap map[string]bool) (passed *passedArgs) {
	passed = &passedArgs{}
	passed.Global = passedMap[string(flagNames.Global)]
	return
}
