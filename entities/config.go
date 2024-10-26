package entities

import (
	"errors"
	"fmt"
	"os"
	"strings"

	//"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"
)

// Typedef for the toml config objects.
type Mothership struct {
	IP   string `toml:"ip"`
	Port string `toml:"port"`
	SSL  *bool  `toml:"ssl"`
}

type Client struct {
	User string `toml:"user_name"`
	ID   string `toml:"user_id"`
}

type Server struct {
	Port     string `toml:"port"`
	SSL      *bool  `toml:"ssl"`
	FilesDir string `toml:"files_dir"`
	CertsDir string `toml:"certs_dir"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}

type Account_t struct {
	Name       string `toml:"user_name"`
	ID         string `toml:"id"`
	Registered bool   `toml:"registered"`
}

type ClientConfig_t struct {
	Mothership Mothership `toml:"mothership"`
	Client     Client     `toml:"client"`
}

type ServerConfig_t struct {
	General  Server      `toml:"general"`
	Accounts []Account_t `toml:"accounts"`
}

func (config *ServerConfig_t) InitConfig(filePath string, users map[string]Account_t) error {
	// Reading config file, filtewatcher can be used after we check if read was ok
	file, err := config.LoadConfig(filePath)
	if err != nil {
		return errors.New("Failed loading Config")
	}

	zap.L().Info("Reading allowed users from config")
	err = toml.Unmarshal(file, &config)
	if err != nil {
		zap.L().Error("Error deconding toml err: " + err.Error())
		return err
	}

	if emptyValues, ok := config.validateConfig(); !ok {
		fmt.Fprintf(os.Stderr, "Missing values in config: %s are empty!\n", strings.Join(emptyValues, ", "))
		return errors.New("Missing values in config file")
	}

	filesDir := config.General.FilesDir
	if _, err := os.Stat(filesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(filesDir, 0777); err != nil {
			zap.L().Error(err.Error())
			return err
		}
	}

	if len(config.Accounts) == 0 {
		fmt.Fprintf(os.Stderr, "Aborting server initialization, no accounts provided")
		return errors.New("Aborting server initialization, no accounts provided")
	}

	zap.L().Info("Generating UUID's for Users")
	for i := range config.Accounts {
		user := &config.Accounts[i]
		id := uuid.New()
		if config.Accounts[i].Registered {
			users[user.Name] = *user
			continue
		}

		config.Accounts[i].ID = id.String()
		users[user.Name] = *user
	}

	tomlContent, err := toml.Marshal(config)
	if err != nil {
		zap.L().Error("Failed to parse config into string.")
		return err
	}

	err = config.WriteConfig(filePath, tomlContent)
	if err != nil {
		zap.L().Error("Failed writing to config")
		return err
	}

	return nil
}

// Load Toml file
// Returns byte slice of content
func (config *ServerConfig_t) LoadConfig(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		zap.L().Error("No config file found! Check README for config example! Error " + err.Error())
		return nil, err
	}

	return file, nil
}

func (config *Server) String() string {
	sslStatus := "nil"
	if config.SSL != nil {
		sslStatus = fmt.Sprintf("%v", *config.SSL)
	}
	return fmt.Sprintf("Port: %s, SSL: %s, FilesDir: %s, CertsDir: %s, CertFile: %s, KeyFile: %s",
		config.Port, sslStatus, config.FilesDir, config.CertsDir, config.CertFile, config.KeyFile)
}

func (sc *ServerConfig_t) String() string {
	accountsStr := ""
	for _, account := range sc.Accounts {
		accountsStr += fmt.Sprintf("Name: %s, ID: %s, Registered: %v\n", account.Name, account.ID, account.Registered)
	}
	return fmt.Sprintf("General: {%s}, Accounts: [%s]", sc.General.String(), accountsStr)
}

// TODO dont hardcode paths
// check filemode
// Function to write to file content
func (config *ServerConfig_t) WriteConfig(path string, content []byte) error {
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		zap.L().Info("Failed writing to toml file. Err: " + err.Error())
		return err
	}

	return nil
}

// Checks if there are mising values in config file.
// Returns empty fields and bool to check if config is valid.
func (config *ServerConfig_t) validateConfig() ([]string, bool) {
	var emptyValues []string

	if config.General.Port == "" {
		emptyValues = append(emptyValues, "Port")
	}

	if config.General.FilesDir == "" {
		emptyValues = append(emptyValues, "FilesDir")
	}

	if config.General.SSL == nil {
		emptyValues = append(emptyValues, "SSL")
	} else if *config.General.SSL {
		if config.General.CertsDir == "" {
			emptyValues = append(emptyValues, "CertsDir")
		}

		if config.General.CertFile == "" {
			emptyValues = append(emptyValues, "CertFile")
		}

		if config.General.KeyFile == "" {
			emptyValues = append(emptyValues, "Keyfile")
		}
	}

	return emptyValues, len(emptyValues) == 0
}
