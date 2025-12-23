package common

type Config struct {
	Log struct {
		Level string `koanf:"level"`
	} `koanf:"log"`
}
