package tray

import (
	"encoding/json"
	"io/ioutil"

	"github.com/adrg/xdg"
)

type LayerConfig struct {
	Key  string `json:"key"`
	Mods string `json:"mods"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Config struct {
	layerInfo []LayerConfig
}

// TODO: If there is no config file, create one.
func LoadConfiguration() (Config, error) {
	cfg := Config{}

	configFile := xdg.ConfigHome + "/kb_ui/config.json"
	file, err := ioutil.ReadFile(configFile)

	if err != nil {
		return Config{}, err
	}

	json.Unmarshal(file, &cfg.layerInfo)

	return cfg, nil
}
