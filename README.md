# kb_ui - Keyboard System Tray Tool

A tiny systray utility to help me keep track of keyboard layers, showing the current layer.

<p align="center">
  <img src="https://user-images.githubusercontent.com/10038688/162264281-c1dbca2d-025b-4091-bdd8-d931679aafe8.PNG">
</p>

## What?

I use a [ZMK](https://zmk.dev) firmware powered board which can have multiple layers of
keybindings, but I don't have a screen or LEDs on the board to keep track of the current
layer.

Since I use certain bindings that conflict depending on the use case (i.e. `a` acts as
both the letter `a` and as a layer swap key, depending on if its tapped or held), I wanted
a very simple way to check the current layer. That way, when I need `a` to be held (say if
I'm playing a game), I can make sure the keyboard is in game mode and not typing mode.
Similarly, I have a different layout when I'm on Mac vs Windows.

## Examples

Starting from some basic, default state...

<img src="https://user-images.githubusercontent.com/10038688/162264281-c1dbca2d-025b-4091-bdd8-d931679aafe8.PNG">

Lets swap to a gaming mode by pressing the layer change button on my board...

<img src="https://user-images.githubusercontent.com/10038688/162264282-ffb7d41c-5a48-4ae0-b8d9-5a15b53cef58.PNG">

Or perhaps to one setup for mac...

<img src="https://user-images.githubusercontent.com/10038688/162264284-48982bcd-5d9d-4797-8be6-2792a991e452.PNG">

Or toggle the output from USB to bluetooth...

<img src="https://user-images.githubusercontent.com/10038688/162264280-ab094c1a-7524-4c77-91b4-4bd54303fb07.PNG">


## Technology

Written in golang, to produce a single, self-contained binary for every platform. That
is, when I press a layer swap key (to move into gaming mode for example), as well as
swapping layer, a unique keycode is sent to the system that `kb_ui` reacts to, swapping
the current state of the tray, updating the tray icon to a user selected icon.
Layers are defined in a json file with their name, icon and the input key-press
that activates it.

## Setup

As stated previously, this doesn't use any HID features or anything so smart, it
relies on keycodes being sent to the host. This means my layer change bindings
had to be updated to include these unique keycodes, and also tell `kb_ui` what they
are.

On the `kb_ui` side, launch it by double clicking, then right click the system tray
icon, configure. That should open up the default JSON config in your editor (or copy
the path and open it yourself if it defaults to your browser).

Most importantly, there are the icon swap bindings, defined as follows:

```json
{
    "key": "1",
    "mods": "ctrl-shift-win-alt",
    "name": "Default",
    "icon": "kb_light",
    "dark_icon": "kb_dark"
}
```

Where `key` defines the input key that will be pressed, `mods` the sequence of
modifier keys that will be held at the same time. `name` is the name you want to
give the layer, `icon` and `dark_icon` are relative paths to the icon that you
want to use, in `ico` format (`kb_light`, `kb_dark` and `disconnected` are built
in icons, so just use strings).

The connect toggle binding is set individually, but works the same.
The connect toggle binding swaps the icon between the current layer icon, and
the disconnect icon (to show output device state).

An example of my config can be found
[here](https://github.com/CrossR/dotfiles/tree/master/kb_ui/.config/kb_ui).

Similarly, my ZMK config, with the macros to actually output these layer change
keybinds can be found
[here](https://github.com/CrossR/zmk_config/blob/master/config/sofle.keymap).
These are just simple macros that send the corresponding keybinding that I've
defined in my config when I swap layer / swap output device.

## Limitations

There is the obvious limitation here, that if I swap my board to my Mac, and
then swap to a mac based layout... my PC has no idea about the layer change, so
when I swap back its out of sync. A pain sure, but one swap and it gets back in
sync and it hardly seems worth fixing when ZMK should in the future be able to
tell the current layer easily enough. So I'll swap to HID codes then, rather
than fixing the odd edge case now.
