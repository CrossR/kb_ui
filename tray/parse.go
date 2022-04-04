package tray

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/CrossR/kb_ui/tray/icons"
	"github.com/adrg/xdg"
	"golang.design/x/hotkey"
)

func ParseKey(key string) (hotkey.Key, error) {

	// TODO: Expand this to support more key names.
	// This is good enough for now.
	switch strings.ToLower(key) {
	case "0":
		return hotkey.Key0, nil
	case "1":
		return hotkey.Key1, nil
	case "2":
		return hotkey.Key2, nil
	case "3":
		return hotkey.Key3, nil
	case "4":
		return hotkey.Key4, nil
	case "5":
		return hotkey.Key5, nil
	case "6":
		return hotkey.Key6, nil
	case "7":
		return hotkey.Key7, nil
	case "8":
		return hotkey.Key8, nil
	case "9":
		return hotkey.Key9, nil
	case "f1":
		return hotkey.Key(0x70), nil
	case "f2":
		return hotkey.Key(0x71), nil
	}

	return hotkey.KeyA, fmt.Errorf("unknown key: %s", key)
}

func loadIconFile(icon_path string) ([]byte, error) {

	full_path, err := xdg.ConfigFile(fmt.Sprintf("kb_ui/%s", icon_path))

	if err != nil {
		return icons.KB_Light_Data, err
	}

	if _, err := os.Stat(full_path); errors.Is(err, os.ErrNotExist) {
		return icons.KB_Light_Data, fmt.Errorf("icon file not found: %s", full_path)
	}

	file, err := ioutil.ReadFile(full_path)

	if err != nil {
		return icons.KB_Light_Data, err
	}

	return file, nil
}

func ParseIcon(icon string) ([]byte, error) {

	lower_icon := strings.ToLower(icon)

	switch lower_icon {
	case "kb_dark":
		return icons.KB_Dark_Data, nil
	case "kb_light":
		return icons.KB_Light_Data, nil
	default:
		return loadIconFile(icon)
	}
}
