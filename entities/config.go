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

type ServerConfig_t struct {
	MotherShipIP   string            `toml:"motherShipIP"`
	MotherShipPort string            `toml:"motherShipPort"`
	MotherShipSSL  bool              `toml:"motherShipSSL"`
	Accounts       []ServerAccount_t `toml:"accounts"`
}

type ServerAccount_t struct {
	Name string `toml:"userName"`
	// FIXME: change to "ID"
	Nonce    string `toml:"nonce"`
	Registered bool   `toml:"registered"`
}

type ServerToml_t struct {
	Server ServerConfig_t `toml:"server"`
}

type ClientConfig_t struct {
	Mothership Mothership `toml:"mothership"`
	Client     Client     `toml:"client"`
}

// Global variable used as shared variable server
// FIXME: Move to server package
var ServerConfig ServerToml_t
