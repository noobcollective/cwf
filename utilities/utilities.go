package utilities

import (
	"flag"
	"os"
	"regexp"

	"go.uber.org/zap"
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

// Loads the given toml file.
// Returns file pointer or error from reading.
func LoadConfig(configPath string) (*os.File, error) {
	file, err := os.OpenFile(configPath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Writes given content to file.
// Returns err of write operation or nil.
func WriteConfig(content []byte, file *os.File) error {
	if _, err := file.WriteAt(content, 0); err != nil {
		zap.L().Info("Failed writing to toml file. Err: " + err.Error())
		return err
	}

	return nil
}
