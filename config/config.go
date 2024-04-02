package config

import (
	"log"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var Raw = koanf.New(".")

func Setup() {
	err := Raw.Load(file.Provider(".env"), dotenv.Parser())
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	err = Raw.Load(env.Provider("", ".", nil), nil)
	if err != nil {
		log.Fatalf("error reading env: %v", err)
	}
}
