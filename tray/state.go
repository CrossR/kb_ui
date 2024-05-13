package tray

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/getlantern/systray"
	"golang.org/x/exp/slices"
)

type TrayState struct {
	tray            *TrayItems
	logger          *log.Logger
	keybinds        *[]Keybinding
	layer_id        int
	layer_name      string
	is_connected    bool
	dark_mode       bool
	disconnect_icon *[]byte
	quitting        bool
}

type SaveState struct {
	LayerId     int    `json:"id"`
	LayerName   string `json:"name"`
	IsConnected bool   `json:"is_connected"`
}

// Get the initial application state.
// Mostly just sets up the logger.
func GetInitialState() TrayState {

	var keybinds []Keybinding

	log_file_path, _ := xdg.DataFile("kb_ui/kb_ui.log")
	f, _ := os.OpenFile(log_file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	logger := log.New(f, "", log.LstdFlags)

	disconnected_icon, _ := ParseIcon("disconnected")
	quitting := false

	return TrayState{nil, logger, &keybinds, 0, "", true, false, &disconnected_icon, quitting}

}

// Save the current state of the application, un-register any keybindings.
func (state *TrayState) SaveCurrentState() {

	dataFile, err := xdg.DataFile("kb_ui/state.json")
	if err != nil {
		state.logger.Printf("Failed to create state file: %s\n", err.Error())
		return
	}

	endState := SaveState{state.layer_id, state.layer_name, state.is_connected}
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

// Load the previous run file, if it exists.
// This can be used to setup the default layer.
// Don't worry about errors, just ignore it since its only a save state of
// the previous state.
func (state *TrayState) LoadPreviousState() {
	dataFile, err := xdg.DataFile("kb_ui/state.json")

	if err != nil {
		state.logger.Printf("Could not find state file: %s\n", err.Error())
		return
	}

	file, err := ioutil.ReadFile(dataFile)

	if err != nil {
		state.logger.Printf("Could not read state file: %s\n", err.Error())
		return
	}

	prevState := SaveState{}
	err = json.Unmarshal(file, &prevState)

	if err != nil {
		state.logger.Printf("Could not unmarshal state file: %s\n", err.Error())
		return
	}

	state.logger.Printf("Loaded previous state: %+v\n", prevState)

	// If this is the state we left off in last time, set it.
	state.tray.layer.SetTitle(fmt.Sprintf("%s Layer", prevState.LayerName))
	i := slices.IndexFunc(*state.keybinds, func(k Keybinding) bool {
		return k.id == prevState.LayerId && k.name == prevState.LayerName
	})
	keybind := (*state.keybinds)[i]

	state.layer_id = i
	state.layer_name = keybind.name
	state.is_connected = prevState.IsConnected

	// Finally, update the tray with this loaded state.
	systray.SetIcon(*keybind.GetIcon(state))
}
