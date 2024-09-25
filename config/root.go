package config

type Config struct {
	Mongo struct {
		DB       string
		Url      string
		User     string
		PASSWORD string
	}
}

// not use toml in lambda
//func NewConfig(path string) *Config {
//	cfg := new(Config)
//	if open, err := os.Open(path); err != nil {
//		panic(err)
//	} else if err := toml.NewDecoder(open).Decode(cfg); err != nil {
//		panic(err)
//	} else {
//		return cfg
//	}
//}
