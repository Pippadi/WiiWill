package main

import (
	"github.com/Pippadi/WiiWill/ui"
	"github.com/Pippadi/loggo"
	actor "gitlab.com/prithvivishak/goactor"
)

func main() {
	loggo.SetLevel(loggo.DebugLevel)

	fromUI, _, err := actor.SpawnRoot(ui.New(), "UI")
	if err != nil {
		loggo.Error(err)
		return
	}

	// window.ShowAndRun must be run from the main thread,
	// so we're passing a nil actor to the first message we receive,
	// which calls window.ShowAndRun directly
	(<-fromUI)(nil)
}
