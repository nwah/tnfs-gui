APP_ID := org.fujinet.tnfsd-gui
NAME := TNFS Server Manager

all: clean build

clean:
	rm -f ./tnfsd/src/*.o
	rm -rf ./dist
	rm -rf ./fyne-cross

build_all: darwin windows linux

macos:
	$(MAKE) -C tnfsd/src OS=BSD
	fyne-cross darwin -arch=amd64 -app-id $(APP_ID) -name "$(NAME)"
	mkdir -p ./dist/macos
	mv "./fyne-cross/dist/darwin-amd64/$(NAME).app" ./dist/macos/
	cp tnfsd/bin/tnfsd "dist/macos/$(NAME).app/Contents/MacOS/"
	cd dist/macos && zip -r "$(NAME) (macOS).zip" "$(NAME).app"

windows:
	$(MAKE) -C tnfsd/src OS=Windows_NT
	fyne-cross windows -arch=amd64 -app-id $(APP_ID) -name "$(NAME)"
	mkdir -p "./dist/windows/$(NAME)"
	unzip "./fyne-cross/dist/windows-amd64/$(NAME).zip" -d "./dist/windows/$(NAME)/"
	cp tnfsd/bin/tnfsd.exe "./dist/windows/$(NAME)"
	cd dist/windows && zip -r "$(NAME) (Windows).zip" "$(NAME)"

# linux:
# 	$(MAKE) -C tnfsd/src OS=LINUX
# 	fyne-cross linux -dir src -arch=amd64 -app-id $(APP_ID) -name "$(NAME)"
