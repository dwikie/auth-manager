package config

import "github.com/spf13/viper"

type AppConfig struct {
	Mode        string `env:"MODE" envDefault:"development"`
	Host        string `env:"HOST" envDefault:"localhost"`
	Port        string `env:"PORT" envDefault:"8080"`
	DatabaseUrl string `env:"DB_URL"`
	Secrets     Secrets
}

type Secrets struct {
	AccessToken  string `env:"ACCESS_TOKEN"`
	RefreshToken string `env:"REFRESH_TOKEN"`
}

func NewAppConfig() *AppConfig {
	// Create a new AppConfig instance and load data from internal/secret/config.toml file
	viperConfig := viper.New()
	viperConfig.SetConfigFile("internal/secret/config.dev.toml")
	viperConfig.ReadInConfig()

	return &AppConfig{
		Mode:        viperConfig.GetString("app.mode"),
		Host:        viperConfig.GetString("app.host"),
		Port:        viperConfig.GetString("app.port"),
		DatabaseUrl: viperConfig.GetString("database.url"),
		Secrets: Secrets{
			AccessToken:  viperConfig.GetString("secrets.access_token"),
			RefreshToken: viperConfig.GetString("secrets.refresh_token"),
		},
	}
}
