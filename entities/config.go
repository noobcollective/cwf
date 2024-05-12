package entities

// Typedef for the yaml config object.
type Config_t struct {
	MotherShipIP   string
	MotherShipPort string
	MotherShipSSL  bool
}

type ServerConfig_t struct {
	MotherShipIP   string            `toml:"motherShipIP"`
	MotherShipPort string            `toml:"motherShipPort"`
	MotherShipSSL  bool              `toml:"motherShipSSL"`
	Accounts       []ServerAccount_t `toml:"accounts"`
}

type ServerAccount_t struct {
	UserName string `toml:"userName"`
	Nonce    string `toml:"nonce"`
	Registed bool   `toml:"registered"`
}

type ServerToml_t struct {
	Server ServerConfig_t `toml:"server"`
}

// Global variable used as shared variable between main,serve and client
var MotherShip Config_t
var ServerConfig ServerToml_t
