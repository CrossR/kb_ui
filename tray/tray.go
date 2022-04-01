package tray

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"

    "golang.design/x/hotkey"
)

func Start() {
	onExit := func() {
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}

	systray.Run(onReady, onExit)
	fmt.Println("We are passed the run...")
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Lantern")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

        hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	err := hk.Register()
	if err != nil {
		return
	}
	fmt.Printf("hotkey: %v is registered\n", hk)

        go func() {
        for {
	<-hk.Keydown()
	fmt.Printf("hotkey: %v is down\n", hk)
        }
        }()

	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("Awesome App")
		systray.SetTooltip("Pretty awesome")

		mChange := systray.AddMenuItem("Click Me", "Click Me")

		systray.AddSeparator()

		for {
			select {
			case <-mChange.ClickedCh:
				mChange.SetTitle("Clicked")
			}
		}
	}()
}
