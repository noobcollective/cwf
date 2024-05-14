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

// Load Toml file
// Returns byte slice of content
func LoadConfig() ([]byte, error) {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		zap.L().Error("Could not retrieve home directory!")
		return nil, err
	}

	file, err := os.ReadFile(usrHome + "/.config/cwf/serverConfig.toml")
	if err != nil {
		zap.L().Error("No config file found! Check README for config example! Error " + err.Error())
		return nil, err
	}

	return file, nil
}

// TODO dont hardcode paths
// check filemode
// Function to write to file content
func WriteConfig(content []byte) error {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		zap.L().Error("Could not retrieve home directory!")
		return err
	}

	err = os.WriteFile(usrHome+"/.config/cwf/serverConfig.toml", content, 0644)
	if err != nil {
		zap.L().Info("Failed writing to toml file. Err: " + err.Error())
		return err
	}

	return nil
}
