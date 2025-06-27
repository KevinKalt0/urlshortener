package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

// TODO: Définir une structure Config avec les sections suivantes:
// - Server (BaseURL string, Port int)
// - Analytics (BufferSize int)
// - Database (DSN string)
// Utiliser les tags mapstructure appropriés
type Config struct {
	Server struct {
		BaseURL string `mapstructure:"base_url"`
		Port    int    `mapstructure:"port"`
	} `mapstructure:"server"`

	Analytics struct {
		BufferSize int `mapstructure:"buffer_size"`
	} `mapstructure:"analytics"`

	Database struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"database"`
}

var (
	config *Config
	once   sync.Once
)

// TODO: Implémenter GetConfig() qui retourne une instance singleton de Config
// Utiliser sync.Once pour s'assurer qu'elle n'est initialisée qu'une seule fois
func GetConfig() *Config {
	once.Do(func() {
		config = loadConfig()
	})
	return config
}

// TODO: Implémenter loadConfig() qui:
// 1. Configure Viper pour lire depuis config.yaml
// 2. Définit des valeurs par défaut raisonnables
// 3. Permet la surcharge par variables d'environnement
// 4. Gère les erreurs de lecture gracieusement
func loadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Valeurs par défaut
	viper.SetDefault("server.base_url", "http://localhost:8080")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("analytics.buffer_size", 1000)
	viper.SetDefault("database.dsn", "urlshortener.db")

	// Variables d'environnement
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v. Using defaults and environment variables.", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	return &cfg
}
