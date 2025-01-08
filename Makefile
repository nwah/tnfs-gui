.PHONY: clean

APP_ID := org.fujinet.tnfsd-gui
NAME := TNFS Server Manager

all: clean macos windows linux

clean:
	rm -f ./tnfsd/src/*.o
	rm -rf ./dist
	rm -rf ./fyne-cross

macos:
ifeq (,$(wildcard bin/tnfsd-bsd))
	$(MAKE) -C tnfsd/src OS=BSD
	mv tnfsd/bin/tnfsd bin/tnfsd-bsd
endif
	fyne-cross darwin -arch=amd64 -app-id $(APP_ID) -name "$(NAME)"
	mkdir -p ./dist/macos
	mv "./fyne-cross/dist/darwin-amd64/$(NAME).app" ./dist/macos/
	cp tnfsd/bin/tnfsd-bsd "dist/macos/$(NAME).app/Contents/MacOS/tnfsd"
	cd dist/macos && zip -r "$(NAME) (macOS).zip" "$(NAME).app"

windows:
ifeq (,$(wildcard bin/tnfsd.exe))
	$(MAKE) -C tnfsd/src OS=Windows_NT
	mv tnfsd/bin/tnfsd.exe bin/tnfsd.exe
endif
	fyne-cross windows -arch=amd64 -app-id $(APP_ID) -name "$(NAME)"
	mkdir -p "./dist/windows/$(NAME)"
	unzip "./fyne-cross/dist/windows-amd64/$(NAME).zip" -d "./dist/windows/$(NAME)/"
	cp tnfsd/bin/tnfsd.exe "./dist/windows/$(NAME)"
	cd dist/windows && zip -r "$(NAME) (Windows).zip" "$(NAME)"

# linux:
# ifeq (,$(wildcard bin/tnfsd-linux))
# 	$(MAKE) -C tnfsd/src OS=LINUX
# 	mv tnfsd/bin/tnfsd bin/tnfsd-linux
# endif
# 	fyne-cross linux -arch=amd64 -app-id $(APP_ID) -name "$(NAME)"
