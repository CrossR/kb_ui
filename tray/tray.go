package tray

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/CrossR/kb_ui/tray/icons"

	"github.com/adrg/xdg"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"golang.design/x/hotkey"
)

type Keybinding struct {
	bind *hotkey.Hotkey
	mods []hotkey.Modifier
	key  hotkey.Key
	id   int
	name string
	icon *[]byte
}

type TrayState struct {
	logger     *log.Logger
	keybinds   *[]Keybinding
	layer_id   int
	layer_name string
}

type SaveState struct {
	PreviousId   int    `json:"id"`
	PreviousName string `json:"name"`
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

	dataFile, err := xdg.DataFile("kb_ui/state.json")
	if err != nil {
		state.logger.Printf("Failed to create state file: %s\n", err.Error())
		return
	}

	endState := SaveState{state.layer_id, state.layer_name}
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
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Load the previous run file, if it exists.
	// This can be used to setup the default layer.
	// Don't worry about errors, just ignore it since its only a save state of
	// the previous state.
	dataFile, err := xdg.DataFile("kb_ui/state.json")
	if err != nil {
		state.logger.Printf("Could not find state file: %s\n", err.Error())
		return
	}
	file, _ := ioutil.ReadFile(dataFile)

	prevState := SaveState{}
	json.Unmarshal(file, &prevState)

	for i, binding := range config.LayerInfo {

		// Parse the configurations strings into its mods / keys and icon.
		mods := ParseModifiers(binding.Mods)
		if len(mods) == 0 {
			state.logger.Printf("Error parsing mods\n")
			continue
		}

		key, err := ParseKey(binding.Key)
		if err != nil {
			state.logger.Printf("Error parsing key: %s\n", err.Error())
			continue
		}

		icon, err := ParseIcon(binding.Icon)
		if err != nil {
			state.logger.Printf("Error parsing icon: %s\n", err.Error())
		}

		keybind := Keybinding{nil, mods, key, i, binding.Name, &icon}
		err = setupKeybinding(state, &keybind, mCurrentLayer)

		if err != nil {
			state.logger.Printf("Error setting up keybind %d: %s\n", i, err.Error())
			continue
		}

		// Store the binding, so we can unregister it later.
		*state.keybinds = append(*state.keybinds, keybind)

		// If this is the state we left off in last time, set it.
		if binding.Name == prevState.PreviousName {
			mCurrentLayer.SetTitle(fmt.Sprintf("%s Layer", keybind.name))
			systray.SetIcon(*keybind.icon)
			state.layer_id = i
			state.layer_name = keybind.name
		}
	}

	// On tray quit event.
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	// On configure event.
	go func() {
		<-mConfigure.ClickedCh
		OpenConfig()
	}()

	// Finally, hook up the info binding.
	infoBind, err := infoKeybind(state, &config)
	if err == nil {
		*state.keybinds = append(*state.keybinds, infoBind)
	} else {
		state.logger.Printf("Failed to create info keybind: %s\n", err.Error())
	}
}

// Setup the actual keybinds to notify the user of layer changes.
func setupKeybinding(state *TrayState, keybind *Keybinding, trayItem *systray.MenuItem) error {

	hk := hotkey.New(keybind.mods, keybind.key)
	err := hk.Register()

	if err != nil {
		return err
	}

	go func() {
		layerSwapName := fmt.Sprintf("Swapped to %s layer", keybind.name)
		for hk != nil {
			<-hk.Keydown()
			// Notify the user of the layer change.
			beeep.Notify("Layer Swapped", layerSwapName, "")

			// Update the tray icon and title.
			trayItem.SetTitle(fmt.Sprintf("%s Layer", keybind.name))
			systray.SetIcon(*keybind.icon)

			// Make sure the app state is saved.
			state.layer_id = keybind.id
			state.layer_name = keybind.name
		}
	}()

	keybind.bind = hk

	return nil
}

// A small helper function that just alerts the user on the current state.
func infoKeybind(state *TrayState, config *Config) (Keybinding, error) {

	mods := ParseModifiers(config.InfoMods)
	if len(mods) == 0 {
		return Keybinding{}, errors.New("info keybind declared with no modifiers")
	}

	key, err := ParseKey(config.InfoKey)
	if err != nil {
		return Keybinding{}, errors.New("failed to parse info key")
	}

	keybind := Keybinding{nil, mods, key, -1, "Info", nil}
	hk := hotkey.New(keybind.mods, keybind.key)
	err = hk.Register()

	if err != nil {
		return Keybinding{}, errors.New("info keybind failed to register")
	}

	go func() {
		for hk != nil {
			<-hk.Keydown()
			beeep.Notify(state.layer_name, "The current keybinding layer is "+state.layer_name, "")
		}
	}()

	keybind.bind = hk

	return keybind, nil
}

func getInitialState() TrayState {

	var keybinds []Keybinding

	log_file_path, _ := xdg.DataFile("kb_ui/kb_ui.log")
	f, _ := os.OpenFile(log_file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	logger := log.New(f, "", log.LstdFlags)

	return TrayState{logger, &keybinds, 0, ""}

}
