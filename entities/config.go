package entities

// Typedef for the toml config objects.
type Mothership struct {
	IP    string `toml:"ip"`
	Port  string `toml:"port"`
	SSL   bool   `toml:"ssl"`
}

type Client struct {
	User  string `toml:"user_name"`
	ID    string `toml:"user_id"`
}

type Server struct {
	Port     string `toml:"port"`
	SSL      bool   `toml:"ssl"`
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
