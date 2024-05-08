package utilities

import (
	"flag"
)

// Get value of a flag.
func GetFlagValue[V bool | string | int](flagName string) V {
	return flag.Lookup(flagName).Value.(flag.Getter).Get().(V)
}
