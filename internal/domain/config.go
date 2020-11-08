package domain

type ConfigRepository interface {
	ReadConfig() (*Config, error)
}

type Config struct {
	Team  string
	Token string
}
