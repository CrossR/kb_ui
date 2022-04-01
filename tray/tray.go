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

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// TODO: Load the user's configuration, such that this is dynamic.
	// TODO: Config should include the layer name, and the keybindings.
	keys := []hotkey.Key{hotkey.KeyS, hotkey.KeyT}

	for _, key := range keys {

		bind := Keybinding{nil, []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, key, 0, ""}
		setupKeybinding(&bind)

		if bind.bind == nil {
			continue
		}

		*state.keybinds = append(*state.keybinds, bind)
	}

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

// Setup the actual keybinds to notify the user of layer changes.
func setupKeybinding(bind *Keybinding) {

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
		}
	}()

	bind.bind = hk
}
