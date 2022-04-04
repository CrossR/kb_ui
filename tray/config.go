package tray

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/adrg/xdg"
	"github.com/skratchdot/open-golang/open"
)

type LayerConfig struct {
	Key     string `json:"key"`
	Mods    string `json:"mods"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	AltIcon string `json:"alt_icon"`
}

type Config struct {
	LayerInfo []LayerConfig `json:"layers"`
	InfoMods  string        `json:"infoMods"`
	InfoKey   string        `json:"infoKey"`
	AltMods   string        `json:"altMods"`
	AltKey    string        `json:"altKey"`
}

func LoadConfiguration() (Config, error) {
	cfg := Config{}

	configFile, _ := xdg.ConfigFile("/kb_ui/config.json")

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		initConfig()
	}

	file, err := ioutil.ReadFile(configFile)

	if err != nil {
		return Config{}, err
	}

	json.Unmarshal(file, &cfg)

	return cfg, nil
}

func initConfig() {

	defaultBind := []LayerConfig{
		{"1", "ctrl-shift-win-alt", "Gaming", "kb_light", "kb_dark"},
	}
	defaultConfig := Config{defaultBind, "ctrl-shift-win-alt", "0", "ctrl-shift-win-alt", "-"}

	json, err := json.MarshalIndent(defaultConfig, "", "    ")

	if err != nil {
		return
	}

	configFile, _ := xdg.ConfigFile("/kb_ui/config.json")
	err = ioutil.WriteFile(configFile, json, 0644)

	if err != nil {
		return
	}

}

// Open the configuration file in the default reader for the JSON file type.
func OpenConfig() {
	configFile, _ := xdg.ConfigFile("/kb_ui/config.json")
	open.Start(configFile)
}
