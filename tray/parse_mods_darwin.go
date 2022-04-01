//go:build darwin

package tray

import (
	"strings"

	"golang.design/x/hotkey"
)

func ParseModifiers(modifiers string) []hotkey.Modifier {

	lower_modifiers := strings.ToLower(modifiers)
	mods := []hotkey.Modifier{}

	if strings.Contains(lower_modifiers, "ctrl") {
		mods = append(mods, hotkey.ModCtrl)
	}

	if strings.Contains(lower_modifiers, "alt") {
		mods = append(mods, hotkey.ModOption)
	}

	if strings.Contains(lower_modifiers, "shift") {
		mods = append(mods, hotkey.ModShift)
	}

	if strings.Contains(lower_modifiers, "win") {
		mods = append(mods, hotkey.ModCmd)
	}

	return mods
}
