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
