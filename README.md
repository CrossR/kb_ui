# kb_ui

A tiny systray utility to help me keep track of keyboard layers. Shows notification and
system tray icon for the current layer.

I use a [ZMK](https://zmk.dev) firmware powered board which can have multiple layers of
keybindings, but I don't have a screen or LEDs on the board to keep track of the current
layer.

Since I use certain bindings that conflict depending on the use case (i.e. `a` acts as
both the letter `a` and as a layer swap key, depending on if its tapped or held), I wanted
a very simple way to check the current layer. That way, when I need `a` to be held (say if
I'm playing a game), I can make sure the keyboard is in game mode and not typing mode.
Similarly, I have a different layout when I'm on Mac vs Windows.

Written in golang, and uses a very simple binding flag to notify of layer swaps. That
is, when I press a layer swap key (to move into gaming mode for example), as well as
swapping layer, a unique keycode is sent to the system that `kb_ui` reacts to, swapping
the current state of the tray, sending a notification using the OS level notification
system, and then updating the tray icon to a user selected icon. Layers are defined in a
json file with their name, icon and the input key-press that activates it.
