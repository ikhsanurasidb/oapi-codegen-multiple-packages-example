package config

import (
    "github.com/kelseyhightower/envconfig"
    "sync"
)

type Config struct {
    MySQL MySQLConfig `envconfig:"DB"`
}

type MySQLConfig struct {
    Host     string `envconfig:"HOST" default:"localhost"`
    Port     int    `envconfig:"PORT" default:"3306"`
    User     string `envconfig:"USER" default:"root"`
    Password string `envconfig:"PASSWORD" default:""`
    Database string `envconfig:"DATABASE" default:"petstore"`
}

var (
    config Config
    once   sync.Once
)

func Get() Config {
    once.Do(func() {
        if err := envconfig.Process("STORE", &config); err != nil {
            panic("Failed to load store configuration: " + err.Error())
        }
    })

    return config
}