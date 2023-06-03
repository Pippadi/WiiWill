# WiiWill [Work in progress]

A simple Wii remote gamepad mapper for Linux.

I'm writing this because other similar programs seem to be ancient, unmaintained, and difficult to use.
My hope with this is to be able to distribute a single package or binary that works out-of-the-box, with no separate driver or library installation necessary.

## How it works

> #### Disclaimer
> I'm learning about how Linux handles input devices on the fly.
> There may well be better ways to achieve what I'm doing here; I'm just going by what worked for me.
> That said, suggestions and contributions are welcome.

1. Monitors `udev` for `uevent`s in which files under `/dev/input` (specified by `DEVNAME`) are created
2. Finds `DEVNAME` for `uevent` satisfying `MAJOR="13"`, `MINOR!="0", and `ID_INPUT_KEY=1`
3. Reads events from this file
4. Writes mapped key event to `/dev/uinput`

Of the several `/dev/input/eventX` files created, only one registers all the buttons (including the D-pad).
The `uevent` which reports this file also reports that it is created by the `input` driver (`MAJOR="13"`), has a `MINOR` not equal to `"0"`, and has `ID_INPUT_KEY` set to `"1"`.
I am not sure what exactly these mean. These are just patterns I've observed.
Time for some `grep`ping in the Linux source code.

## Acknowledgements

- [nervo/wiican](https://github.com/nervo/wiican) for the inspiration
- [Oblomov/wiimote-pad](https://github.com/Oblomov/wiimote-pad) and [xwiimote](https://github.com/xwiimote/xwiimote) as references
