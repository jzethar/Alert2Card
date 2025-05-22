package types

import (
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Host     string   `yaml:"host"`
	Ports    Ports    `yaml:"ports"`
	Youtrack Youtrack `yaml:"youtrack"`
}

type Youtrack struct {
	Host      string `yaml:"host"`
	Key       string `yaml:"key"`
	ProjectID string `yaml:"projectID"`
}

type Ports struct {
	Rpc   string `yaml:"rpc"`
	Debug string `yaml:"debug"`
}

func (dc *Config) LoadConfig(configFilePath string) error {

	var err error
	var file *os.File
	var data []byte

	if file, err = os.Open(configFilePath); err != nil {
		log.Fatal().Msg("can't open file")
	}
	defer file.Close()

	if data, err = io.ReadAll(file); err != nil {
		log.Fatal().Msg("can't read file")
	}
	if err = yaml.Unmarshal(data, &dc); err != nil {
		log.Fatal().Msg("can't unmarshal file")
	}
	return nil
}
