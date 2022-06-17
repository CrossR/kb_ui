package tray

import (
	"fmt"

	"github.com/CrossR/kb_ui/tray/icons"

	"github.com/getlantern/systray"
)

type TrayItems struct {
	layer  *systray.MenuItem
	quiet  *systray.MenuItem
	config *systray.MenuItem
	quit   *systray.MenuItem
}

var Version string

// On ready, load the user configuration, setup the keybindings, then just wait
// and react to them.
func SetupInitialTrayState(state *TrayState) {

	// Setup the default parts of the system tray.
	systray.SetTemplateIcon(icons.KB_Dark_Data, icons.KB_Dark_Data)
	systray.SetTooltip(fmt.Sprintf("Keyboard Status (%s)", GetVersion()))

	// Set a default layer, so there is something set before checking for a
	// previous run file.
	mCurrentLayer := systray.AddMenuItem("Default Layer", "The current keyboard layer")

	// Add the final entries to configure or quit the application.
	systray.AddSeparator()
	mConfigure := systray.AddMenuItem("Configure", "Open the app config file")
	mQuiet := systray.AddMenuItem("Quiet notifications", "Don't show notifications")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

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
			state.setQuiet()
		}
	}()

	state.tray = &TrayItems{mCurrentLayer, mQuiet, mConfigure, mQuit}
}

// Get version string, this will be set dynamically for releases to git hash.
func GetVersion() string {
	if Version != "" {
		return Version
	}

	return "Dev"
}
