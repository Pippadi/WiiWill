#!/bin/bash

cd "$(dirname $0)"/..
repodir="$PWD"
cd -
tarfile="WiiWill.tar.xz"

cd "$repodir"
make clean
make tar

cd "$repodir"/fyne-cross/dist/linux-amd64/
tar -xvf "$tarfile"
rm "$tarfile"

mkdir -p etc/modules-load.d
mkdir -p etc/udev/rules.d
cp "$repodir"/pkg/uinput.conf etc/modules-load.d
cp "$repodir"/pkg/49-wiiwill.rules etc/udev/rules.d
cp "$repodir"/pkg/Makefile .

tar -cJvf "$repodir/pkg/$tarfile" *
cd "$(dirname $0)"
