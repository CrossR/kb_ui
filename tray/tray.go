package tray

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/CrossR/kb_ui/tray/icons"

	"github.com/adrg/xdg"
	"github.com/getlantern/systray"
)

type TrayState struct {
	logger          *log.Logger
	keybinds        *[]Keybinding
	layer_id        int
	layer_name      string
	is_connected    bool
	dark_mode       bool
	quiet           bool
	disconnect_icon *[]byte
}

type SaveState struct {
	PreviousId   int    `json:"id"`
	PreviousName string `json:"name"`
	WasConnected bool   `json:"was_connected"`
	Quiet        bool   `json:"quiet"`
}

func Start() {

	trayState := getInitialState()

	onReady := func() {
		traySetup(&trayState)
	}
	onExit := func() {
		trayEnd(&trayState)
	}

	systray.Run(onReady, onExit)
}

// On exit, save the current state of the application, un-register any keybindings.
func trayEnd(state *TrayState) {
	for _, hk := range *state.keybinds {
		err := hk.bind.Unregister()

		if err != nil {
			state.logger.Println("Failed to unregister keybind:", err.Error())
		}

		hk.bind = nil
	}

	state.logger.Printf("Final state was %+v\n", state)

	dataFile, err := xdg.DataFile("kb_ui/state.json")
	if err != nil {
		state.logger.Printf("Failed to create state file: %s\n", err.Error())
		return
	}

	endState := SaveState{state.layer_id, state.layer_name, state.is_connected, state.quiet}
	json, err := json.MarshalIndent(endState, "", "    ")
	if err != nil {
		state.logger.Printf("Failed to marshall state: %s\n", err.Error())
		return
	}

	err = ioutil.WriteFile(dataFile, json, 0644)

	if err != nil {
		state.logger.Printf("Failed to save state: %s\n", err.Error())
		return
	}

	state.logger.Printf("Saved state %+v\n", endState)
}

// On ready, load the user configuration, setup the keybindings, then just wait
// and react to them.
func traySetup(state *TrayState) {

	// Setup the default parts of the system tray.
	systray.SetTemplateIcon(icons.KB_Dark_Data, icons.KB_Dark_Data)
	systray.SetTooltip("Keyboard Status")

	// Load the user config, but stop if there is nothing defined.
	config, err := LoadConfiguration()

	if len(config.LayerInfo) == 0 {
		state.logger.Println("No layers defined, exiting.")
		systray.Quit()
		return
	}

	// Set a default layer, so there is something set before checking for a
	// previous run file.
	defaultName := fmt.Sprintf("%s Layer", config.LayerInfo[0].Name)
	mCurrentLayer := systray.AddMenuItem(defaultName, defaultName)

	if err != nil {
		state.logger.Printf("Error loading configuration: %s\n", err.Error())
		systray.Quit()
		return
	}

	// Add the final entries to configure or quit the application.
	systray.AddSeparator()
	mConfigure := systray.AddMenuItem("Configure", "Configure")
	mQuiet := systray.AddMenuItem("Quiet", "Don't show notifications")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Load the previous run file, if it exists.
	// This can be used to setup the default layer.
	// Don't worry about errors, just ignore it since its only a save state of
	// the previous state.
	dataFile, err := xdg.DataFile("kb_ui/state.json")
	if err != nil {
		state.logger.Printf("Could not find state file: %s\n", err.Error())
	}
	file, _ := ioutil.ReadFile(dataFile)

	prevState := SaveState{}
	json.Unmarshal(file, &prevState)

	state.logger.Printf("Loaded previous state: %+v\n", prevState)
	state.quiet = prevState.Quiet

	// Parse the actual layer bindings out.
	for i, binding := range config.LayerInfo {

		keybind := MakeKeybinding(state, binding, i)
		err = keybind.SetupKeybinding(state, mCurrentLayer)

		if err != nil {
			state.logger.Printf("Error setting up keybind %d: %s\n", i, err.Error())
			continue
		}

		// Store the binding, so we can unregister it later.
		*state.keybinds = append(*state.keybinds, keybind)

		// If this is the state we left off in last time, set it.
		if binding.Name == prevState.PreviousName {
			mCurrentLayer.SetTitle(fmt.Sprintf("%s Layer", keybind.name))
			state.layer_id = i
			state.layer_name = keybind.name
			state.is_connected = prevState.WasConnected
			systray.SetIcon(*keybind.GetIcon(state))
		}
	}

	// On tray quit event.
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	// On configure event.
	go func() {
		for mConfigure != nil {
			<-mConfigure.ClickedCh
			OpenConfig()
		}
	}()

	// Toggle the quiet mode.
	go func() {
		for mQuiet != nil {
			<-mQuiet.ClickedCh
			state.quiet = !state.quiet

			if state.quiet {
				mQuiet.SetTitle("Notify")
			} else {
				mQuiet.SetTitle("Quiet")
			}
		}
	}()

	// Finally, hook up the auxillary bindings.
	infoBind, err := SetupInfoKeybind(state, &config)
	if err == nil {
		*state.keybinds = append(*state.keybinds, infoBind)
	} else {
		state.logger.Printf("Failed to create info keybind: %s\n", err.Error())
	}

	connectToggleBinding, err := SetupConnectKeybind(state, &config)
	if err == nil {
		*state.keybinds = append(*state.keybinds, connectToggleBinding)
	} else {
		state.logger.Printf("Failed to create connect toggle keybind: %s\n", err.Error())
	}
}

// Get the initial application state.
// Mostly just sets up the logger.
func getInitialState() TrayState {

	var keybinds []Keybinding

	log_file_path, _ := xdg.DataFile("kb_ui/kb_ui.log")
	f, _ := os.OpenFile(log_file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	logger := log.New(f, "", log.LstdFlags)

	disconnected_icon, _ := ParseIcon("disconnected")

	return TrayState{logger, &keybinds, 0, "", true, false, false, &disconnected_icon}

}
