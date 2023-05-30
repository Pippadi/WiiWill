# WiiWill [Work in progress]

A simple Wii remote gamepad mapper for Linux.

I'm writing this because other similar programs seem to be ancient, unmaintained, and difficult to use.
My hope with this is to be able to distribute a single package or binary that works out-of-the-box, with no separate driver or library installation necessary.

## How it works

> #### Disclaimer
> I'm learning about how Linux handles input devices on the fly.
> There may well be better ways to achieve what I'm doing here; I'm just going by what worked for me.
> That said, suggestions and contributions are welcome.

1. Scans for bluetooth devices which have `RVL-CNT-01` in their advertised name, and connects to the first one
2. Monitors `udev` for `uevent`s in which files under `/dev/input` (specified by `DEVNAME`) are created
3. Finds `DEVNAME` for `uevent` containing `MAJOR="13"` and `MINOR="79"`
4. Reads events from this file

Of the several `/dev/input/eventX` files created, only one registers all the buttons (including the D-pad). The `uevent` which reports this file also reports that it is created by the `input` driver (`MAJOR="13"`), and has a `MINOR="79"`.
I am not sure what the significance of the number 79 is, but I feel it is related to [this](https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h#L31).

This is all that has been implemented so far. `uinput` will be used to generate mouse/keypresses after the buttons have been remapped.

## Acknowledgements

- [nervo/wiican](https://github.com/nervo/wiican) for the inspiration
- [Oblomov/wiimote-pad](https://github.com/Oblomov/wiimote-pad) and [xwiimote](https://github.com/xwiimote/xwiimote) as references
