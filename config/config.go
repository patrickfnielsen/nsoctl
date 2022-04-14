/*
Copyright © 2022 Patrick Falk Nielsen <git@patricknielsen.dk>
*/
package config

type NsoConfig struct {
	InsecureSkipVerify bool
	ServerFqdn         string
	Username           string
	Password           string
}

type Config struct {
	Nso NsoConfig `mapstructure:"nso"`
}
