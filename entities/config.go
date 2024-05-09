package entities

// Typedef for the yaml config object.
type Config_t struct {
	MotherShipIP    string `yaml:"motherShipIP"`
	MotherShipPort	string `yaml:"motherShipPort"`
	MotherShipSSL   bool   `yaml:"motherShipSSL"`
}

// Global variable used as shared variable between main,serve and client
var MotherShip Config_t
