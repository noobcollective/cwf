package entities

// Typedef for the yaml config object.
type Config_t struct {
	MotherShipIP           string `yaml:"motherShipIP"`
	MotherShipPort         string `yaml:"motherShipPort"`
	MotherShipCWFDirectory string `yaml:"motherShipCWFDirectory"`
}

// Global variable used as shared variable between main,serve and client
var MotherShip Config_t
