package tray

import (
	"github.com/getlantern/systray"
)

func Start() {

	trayState := GetInitialState()

	onReady := func() {
		appStart(&trayState)
	}
	onExit := func() {
		appEnd(&trayState)
	}

	systray.Run(onReady, onExit)
}

// On exit, save the current state of the application, un-register any keybindings.
func appEnd(state *TrayState) {

	state.quitting = true

	for _, hk := range *state.keybinds {
		err := hk.bind.Unregister()

		if err != nil {
			state.logger.Println("Failed to unregister keybind:", err.Error())
		}

		hk.bind = nil
	}

	state.logger.Printf("Final state was %+v\n", state)

	state.SaveCurrentState()
}

// On ready, load the user configuration, setup the keybindings, then just wait
// and react to them.
func appStart(state *TrayState) {

	SetupInitialTrayState(state)

	// Load the user config, but stop if there is nothing defined.
	config, err := LoadConfiguration()

	if err != nil {
		state.logger.Printf("Error loading configuration: %s\n", err.Error())
		systray.Quit()
		return
	} else if len(config.LayerInfo) == 0 {
		state.logger.Println("No layers defined, exiting.")
		systray.Quit()
		return
	}

	// Load the actual user disconnect icon.
	*state.disconnect_icon, err = ParseIcon(config.DisconnectIcon)

	if err != nil {
		state.logger.Printf("Error parsing disconnect icon: %s\n", err.Error())
		*state.disconnect_icon, _ = ParseIcon("disconnected")
	}

	// Parse the actual layer bindings out.
	for i, binding := range config.LayerInfo {

		keybind := MakeKeybinding(state, binding, i)
		err = keybind.SetupKeybinding(state)

		if err != nil {
			state.logger.Printf("Error setting up keybind %d: %s\n", i, err.Error())
			continue
		}

		// Store the binding, so we can unregister it later.
		*state.keybinds = append(*state.keybinds, keybind)
	}

	// Set the initial state of the application, if there is one.
	state.LoadPreviousState()

	connectToggleBinding, err := SetupConnectKeybind(state, &config)
	if err == nil {
		*state.keybinds = append(*state.keybinds, connectToggleBinding)
	} else {
		state.logger.Printf("Failed to create connect toggle keybind: %s\n", err.Error())
	}
}
