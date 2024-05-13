package utilities

import (
	"flag"
	"regexp"
)

// Get value of a flag.
func GetFlagValue[V bool | string | int](flagName string) V {
	return flag.Lookup(flagName).Value.(flag.Getter).Get().(V)
}

// Function to check if gived uuid is valid
// Regex not by me i just trust people from the internet
// Could use a library :)
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
