package tray

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/adrg/xdg"
)

type LayerConfig struct {
	Key  string `json:"key"`
	Mods string `json:"mods"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Config struct {
	LayerInfo []LayerConfig
}

func LoadConfiguration() (Config, error) {
	cfg := Config{}

	configFile := xdg.ConfigHome + "/kb_ui/config.json"

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		initConfig()
	}

	file, err := ioutil.ReadFile(configFile)

	if err != nil {
		return Config{}, err
	}

	json.Unmarshal(file, &cfg.LayerInfo)

	return cfg, nil
}

func initConfig() {

	defaultBind := []LayerConfig{
		{"1", "ctrl-shift-win-alt", "Gaming", "kb_gaming"},
	}
	defaultConfig := Config{defaultBind}

	json, err := json.MarshalIndent(defaultConfig, "", "    ")

	if err != nil {
		return
	}

	configFile := xdg.ConfigHome + "/kb_ui/config.json"
	err = ioutil.WriteFile(configFile, json, 0644)

	if err != nil {
		return
	}

}
