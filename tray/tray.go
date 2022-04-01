package tray

import (
	"fmt"

	"github.com/CrossR/kb_ui/tray/icons"

	"github.com/getlantern/systray"

	"github.com/gen2brain/beeep"
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

	// TODO: Save the current state of the application to load on next run.
}

// On ready, load the user configuration, setup the keybindings, then just wait
// and react to them.
func traySetup(state *TrayState) {

	systray.SetTemplateIcon(icons.KB_Dark_Data, icons.KB_Dark_Data)
	systray.SetTooltip("Keyboard Status")

	config, err := LoadConfiguration()

	if len(config.layerInfo) == 0 {
		systray.Quit()
		return
	}

	defaultName := fmt.Sprintf("%s Layer", config.layerInfo[0].Name)
	mCurrentLayer := systray.AddMenuItem(defaultName, defaultName)

	if err != nil {
		systray.Quit()
		return
	}

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for i, binding := range config.layerInfo {

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

		keybind := Keybinding{nil, mods, key, i, binding.Name, icon}
		setupKeybinding(&keybind, mCurrentLayer)

		if keybind.bind == nil {
			continue
		}

		*state.keybinds = append(*state.keybinds, keybind)
	}

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

// Setup the actual keybinds to notify the user of layer changes.
func setupKeybinding(bind *Keybinding, trayItem *systray.MenuItem) {

	hk := hotkey.New(bind.mods, bind.key)
	err := hk.Register()

	if err != nil {
		return
	}

	go func() {
		layerSwapName := fmt.Sprintf("Swapped to %s layer", bind.name)
		for hk != nil {
			<-hk.Keydown()
			beeep.Notify("Layer Swapped", layerSwapName, "")
			trayItem.SetTitle(fmt.Sprintf("%s Layer", bind.name))
		}
	}()

	bind.bind = hk
}
