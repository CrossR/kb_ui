package tray

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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
	keybinds   *[]Keybinding
	layer_id   int
	layer_name string
}

type SaveState struct {
	PreviousId   int    `json:"id"`
	PreviousName string `json:"name"`
}

func Start() {
	var keybinds []Keybinding
	trayState := TrayState{&keybinds, 0, ""}

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
		hk.bind.Unregister()
		hk.bind = nil
	}

	dataFile, err := xdg.DataFile("kb_ui/state.json")
	if err != nil {
		return
	}

	endState := SaveState{state.layer_id, state.layer_name}
	json, err := json.MarshalIndent(endState, "", "    ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(dataFile, json, 0644)

	if err != nil {
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
		systray.Quit()
		return
	}

	// Set a default layer, so there is something set before checking for a
	// previous run file.
	defaultName := fmt.Sprintf("%s Layer", config.LayerInfo[0].Name)
	mCurrentLayer := systray.AddMenuItem(defaultName, defaultName)

	if err != nil {
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
		return
	}
	file, _ := ioutil.ReadFile(dataFile)

	prevState := SaveState{}
	json.Unmarshal(file, &prevState)

	for i, binding := range config.LayerInfo {

		// Parse the configurations strings into its mods / keys and icon.
		mods := ParseModifiers(binding.Mods)
		if len(mods) == 0 {
			systray.Quit()
			break
		}

		key, err := ParseKey(binding.Key)
		if err != nil {
			systray.Quit()
			break
		}

		icon, err := ParseIcon(binding.Icon)
		if err != nil {
			systray.Quit()
			break
		}

		keybind := Keybinding{nil, mods, key, i, binding.Name, &icon}
		setupKeybinding(state, &keybind, mCurrentLayer)

		if keybind.bind == nil {
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
	infoKeybind(state, &config)
}

// Setup the actual keybinds to notify the user of layer changes.
func setupKeybinding(state *TrayState, keybind *Keybinding, trayItem *systray.MenuItem) {

	hk := hotkey.New(keybind.mods, keybind.key)
	err := hk.Register()

	if err != nil {
		return
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
}

// A small helper function that just alerts the user on the current state.
func infoKeybind(state *TrayState, config *Config) *Keybinding {

	mods := ParseModifiers(config.InfoMods)
	if len(mods) == 0 {
		return nil
	}

	key, err := ParseKey(config.InfoKey)
	if err != nil {
		systray.Quit()
	}

	keybind := Keybinding{nil, mods, key, -1, "Info", nil}
	hk := hotkey.New(keybind.mods, keybind.key)
	err = hk.Register()

	if err != nil {
		return nil
	}

	go func() {
		for hk != nil {
			<-hk.Keydown()
			beeep.Notify(state.layer_name, "The current keybinding layer is "+state.layer_name, "")
		}
	}()

	keybind.bind = hk

	return &keybind
}
