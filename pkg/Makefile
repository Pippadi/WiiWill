# If PREFIX isn't provided, we check for $(DESTDIR)/usr/local and use that if it exists.
# Otherwice we fall back to using /usr.

LOCAL != test -d $(DESTDIR)/usr/local && echo -n "/local" || echo -n ""
LOCAL ?= $(shell test -d $(DESTDIR)/usr/local && echo "/local" || echo "")
PREFIX ?= /usr$(LOCAL)

Name := "WiiWill"
Exec := "WiiWill"
Icon := "WiiWill.png"

default:
	# User install
	# Run "make user-install" to install in ~/.local/
	# Run "make user-uninstall" to uninstall from ~/.local/
	#
	# System install
	# Run "sudo make install" to install the application.
	# Run "sudo make uninstall" to uninstall the application.

install:
	install -Dm00644 usr/local/share/applications/$(Name).desktop $(DESTDIR)$(PREFIX)/share/applications/$(Name).desktop
	install -Dm00755 usr/local/bin/$(Exec) $(DESTDIR)$(PREFIX)/bin/$(Exec)
	install -Dm00644 usr/local/share/pixmaps/$(Icon) $(DESTDIR)$(PREFIX)/share/pixmaps/$(Icon)
	install -Dm00644 etc/modules-load.d/uinput.conf /etc/modules-load.d/uinput.conf
	install -Dm00644 etc/udev/rules.d/49-wiiwill.rules /etc/udev/rules.d/49-wiiwill.rules
	udevadm control --reload-rules
	modprobe uinput
uninstall:
	-rm $(DESTDIR)$(PREFIX)/share/applications/$(Name).desktop
	-rm $(DESTDIR)$(PREFIX)/bin/$(Exec)
	-rm $(DESTDIR)$(PREFIX)/share/pixmaps/$(Icon)
	-rm /etc/modules-load.d/uinput.conf
	-rm /etc/udev/rules.d/49-wiiwill.rules
