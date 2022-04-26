package tray

import (
	"errors"
	"fmt"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"golang.design/x/hotkey"
)

type Keybinding struct {
	bind      *hotkey.Hotkey
	mods      []hotkey.Modifier
	key       hotkey.Key
	id        int
	name      string
	icon      *[]byte
	dark_icon *[]byte
}

func MakeKeybinding(state *TrayState, binding LayerConfig, i int) Keybinding {

	// Parse the configurations strings into its mods / keys and icon.
	mods := ParseModifiers(binding.Mods)
	if len(mods) == 0 {
		state.logger.Printf("Error parsing mods\n")
		return Keybinding{}
	}

	key, err := ParseKey(binding.Key)
	if err != nil {
		state.logger.Printf("Error parsing key: %s\n", err.Error())
		return Keybinding{}
	}

	icon, err := ParseIcon(binding.Icon)
	if err != nil {
		state.logger.Printf("Error parsing icon: %s\n", err.Error())
	}

	dark_icon, err := ParseIcon(binding.DarkIcon)
	if err != nil {
		state.logger.Printf("Error parsing icon: %s\n", err.Error())
	}

	keybind := Keybinding{nil, mods, key, i, binding.Name, &icon, &dark_icon}

	return keybind
}

// Setup the actual keybinds to notify the user of layer changes.
func (keybind *Keybinding) SetupKeybinding(state *TrayState) error {

	hk := hotkey.New(keybind.mods, keybind.key)
	err := hk.Register()

	if err != nil {
		return err
	}

	go func() {
		layerSwapName := fmt.Sprintf("Swapped to %s layer", keybind.name)
		for hk != nil {
			<-hk.Keydown()

			if state.quitting {
				break
			}

			// If we are already in this layer, do nothing.
			if state.layer_id == keybind.id {
				continue
			}

			// Update the tray icon and title.
			state.tray.layer.SetTitle(fmt.Sprintf("%s Layer", keybind.name))
			systray.SetIcon(*keybind.GetIcon(state))

			// Make sure the app state is saved.
			state.layer_id = keybind.id
			state.layer_name = keybind.name

			// If quiet, don't alert the user.
			if state.quiet {
				continue
			}

			// Notify the user of the layer change.
			err := beeep.Notify("Layer Swapped", layerSwapName, "")

			if err != nil {
				state.logger.Printf("Error notifying user: %s\n", err.Error())
			}
		}
	}()

	keybind.bind = hk

	return nil
}

// A small helper function that just alerts the user on the current state.
func SetupInfoKeybind(state *TrayState, config *Config) (Keybinding, error) {

	mods := ParseModifiers(config.InfoMods)
	if len(mods) == 0 {
		return Keybinding{}, errors.New("info keybind declared with no modifiers")
	}

	key, err := ParseKey(config.InfoKey)
	if err != nil {
		return Keybinding{}, errors.New("failed to parse info key")
	}

	keybind := Keybinding{nil, mods, key, -1, "Info", nil, nil}
	hk := hotkey.New(keybind.mods, keybind.key)
	err = hk.Register()

	if err != nil {
		return Keybinding{}, errors.New("info keybind failed to register")
	}

	go func() {
		for hk != nil {
			<-hk.Keydown()

			if state.quitting {
				break
			}

			if state.is_connected {
				beeep.Notify(state.layer_name, "The current keybinding layer is "+state.layer_name, "")
			} else {
				beeep.Notify("Disconnected", "The keyboard is currently not connected...", "")
			}
		}
	}()

	keybind.bind = hk

	return keybind, nil
}

// A small helper function that just toggles the disconnected icon.
// I.e., when the board swaps output to another device, swap to an icon
// that shows this disconnected state.
func SetupConnectKeybind(state *TrayState, config *Config) (Keybinding, error) {

	mods := ParseModifiers(config.ConnectMods)
	if len(mods) == 0 {
		return Keybinding{}, errors.New("connect toggle keybind declared with no modifiers")
	}

	key, err := ParseKey(config.ConnectKey)
	if err != nil {
		return Keybinding{}, errors.New("failed to parse connect toggle key")
	}

	keybind := Keybinding{nil, mods, key, -1, "Connect Toggle", nil, nil}
	hk := hotkey.New(keybind.mods, keybind.key)
	err = hk.Register()

	if err != nil {
		return Keybinding{}, errors.New("connect toggle keybind failed to register")
	}

	go func() {
		for hk != nil {
			<-hk.Keydown()

			if state.quitting {
				break
			}

			// Update the tray icon and title.
			bind := (*state.keybinds)[state.layer_id]
			state.is_connected = !state.is_connected
			systray.SetIcon(*bind.GetIcon(state))
		}
	}()

	keybind.bind = hk

	return keybind, nil
}

// Get the current app icon.
func (keybind *Keybinding) GetIcon(state *TrayState) *[]byte {
	if state.is_connected {
		if state.dark_mode {
			return keybind.dark_icon
		} else {
			return keybind.icon
		}
	} else {
		return state.disconnect_icon
	}
}
