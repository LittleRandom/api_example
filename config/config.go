package config

type Config struct {
	Host     string
	Port     int
	DataPath string
}

func NewConfig() *Config {
	return &Config{
		Host:     HOST,
		Port:     PORT,
		DataPath: DATAPATH,
	}
}
