#!/bin/sh
fyne-cross --targets=linux/amd64 .
mv build/main-linux-amd64 cistatus
fyne package -os linux -icon icon.png
tar -zcvf linux.tar.gz usr
rm -rf usr
mv linux.tar.gz release/.

fyne-cross --targets=windows/amd64 .
mv build/main-windows-amd64.exe cistatus.exe
fyne package -os windows -icon icon.png
fyne-cross --targets=windows/amd64 .
mv build/main-windows-amd64.exe cistatus.exe
fyne package -os windows -icon icon.png
mv cistatus.exe release/.

fyne-cross --targets=darwin/amd64 .
mv build/main-darwin-amd64 cistatus
fyne package -os darwin -icon icon.png
rm -rf release/cistatus.app
mv cistatus.app release/.
tar -zcvf release/osx.tar.gz release/cistatus.app
cp -r release/cistatus.app /Applications/.