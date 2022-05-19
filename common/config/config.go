package config

type Config = ConfigV1

func Read() (*Config, error) {
	return ReadV1()
}

func Write(config *Config) error {
	return WriteV1(config)
}
