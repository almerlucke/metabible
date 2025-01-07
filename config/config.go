package config

import (
	"encoding/json"
	"github.com/almerlucke/metabible/pkg/store"
	"io"
	"os"
)

type Config struct {
	Mongo *store.Config `json:"mongo"`
}

func MustLoad(file string) *Config {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return &config
}
