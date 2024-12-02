package config

type Config struct {
	Mongo struct {
		DB       string
		Url      string
		User     string
		PASSWORD string
	}
}
