.PHONY: clean

WiiWill:
	go build

tar:
	fyne-cross linux -name "WiiWill" -icon "assets/Icon.svg" -release -app-id "dev.prithvi.WiiWill"

clean:
	go clean
	go clean --cache
