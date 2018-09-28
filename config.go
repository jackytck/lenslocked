package main

import "fmt"

// # models/users.go
// const userPwPepper = "P4P]tV6$LZc;,bu5"
// const hmacSecretKey = "E4j!STJ$??cc]UhQ"
//
// # models/services.go
// db, err := gorm.Open("postgres", connectionInfo)
// db.LogMode(true)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "jacky",
		Password: "natnat",
		Name:     "lenslocked_dev",
	}
}

type Config struct {
	Port int
	Env  string
}

func (c Config) isProd() bool {
	return c.Env == "prod"
}

func DefaultConfig() Config {
	return Config{
		Port: 3000,
		Env:  "dev",
	}
}
