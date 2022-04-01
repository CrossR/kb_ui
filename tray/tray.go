package tray

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"

	"github.com/gen2brain/beeep"
	"golang.design/x/hotkey"
)

func Start() {
	systray.Run(onReady, onExit)
}

// On exit, save the current state of the application, un-register any keybindings.
func onExit() {
}

// On ready, load the user configuration, setup the keybindings, then just wait
// and react to them.
func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTooltip("Keyboard Status")

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

        setupKeybinding([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	setupKeybinding([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyT)

	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
}

func setupKeybinding(mods []hotkey.Modifier, key hotkey.Key) {

	hk := hotkey.New(mods, key)
	err := hk.Register()
	if err != nil {
		return
	}

	go func() {
		for {
			<-hk.Keydown()
			beeep.Notify("Pressed", "Down", "")
		}
	}()
}
