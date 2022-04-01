package tray

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"

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
}

// On ready, load the user configuration, setup the keybindings, then just wait
// and react to them.
func traySetup(state *TrayState) {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTooltip("Keyboard Status")

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	keys := []hotkey.Key{hotkey.KeyS, hotkey.KeyT}

	for _, key := range keys {
		hk, ok := setupKeybinding([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, key)

		if !ok {
			continue
		}

		*state.keybinds = append(*state.keybinds, hk)
	}

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func setupKeybinding(mods []hotkey.Modifier, key hotkey.Key) (bind Keybinding, ok bool) {

	hk := hotkey.New(mods, key)
	err := hk.Register()

	if err != nil {
		return Keybinding{}, false
	}

	go func() {
		hotKeyPressed := fmt.Sprintf("%s was pressed", hk.String())
		for hk != nil {
			<-hk.Keydown()
			beeep.Notify("Pressed", hotKeyPressed, "")
		}
	}()

	return Keybinding{hk, mods, key, 0, ""}, true
}
