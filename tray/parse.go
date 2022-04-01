package tray

import (
	"fmt"
	"strings"

	"github.com/CrossR/kb_ui/tray/icons"
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
	}

	return hotkey.KeyA, fmt.Errorf("unknown key: %s", key)
}

func ParseIcon(icon string) ([]byte, error) {

	lower_icon := strings.ToLower(icon)

	switch lower_icon {
	case "kb_dark":
		return icons.KB_Dark_Data, nil
	case "kb_light":
		return icons.KB_Light_Data, nil
	case "gaming_dark":
		return icons.Game_Dark_Data, nil
	case "gaming_light":
		return icons.Game_Light_Data, nil
	case "mac_dark":
		return icons.Mac_Dark_Data, nil
	case "mac_light":
		return icons.Mac_Light_Data, nil
	}

	return icons.KB_Dark_Data, fmt.Errorf("unknown icon: %s", icon)
}
