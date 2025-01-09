package ui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

int hideFromDock(void) {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    [[NSRunningApplication currentApplication] activateWithOptions:NSApplicationActivateAllWindows];
    return 0;
}

int showInDock(void) {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    return 0;
}
*/
import "C"

func HideFromDock() {
	C.hideFromDock()
}

func ShowInDock() {
	C.showInDock()
}
