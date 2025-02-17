.PHONY: clean tagversion

NAME := TNFS Server Manager
VERSION := $(shell cat FyneApp.toml | sed -n '/Version *=* */{s///;s/^ *"//;s/"$$//;p;}')

all: clean macos windows linux

tagversion:
	git tag -a $(VERSION) -m "Version $(VERSION)"

clean:
	rm -f ./tnfsd/src/*.o
	rm -rf ./dist
	rm -rf ./fyne-cross

macos:
	fyne-cross darwin -arch=amd64
	mkdir -p ./dist/macos
	mv "./fyne-cross/dist/darwin-amd64/$(NAME).app" ./dist/macos/
	cd dist/macos && zip -r "$(NAME) (macOS).zip" "$(NAME).app"

windows:
	fyne-cross windows -arch=amd64
	mkdir -p ./dist/windows/
	mv "fyne-cross/dist/windows-amd64/$(NAME).exe.zip" "./dist/windows/$(NAME) (Windows).zip"

# linux:
# 	fyne-cross linux -arch=amd64
