// +build darwin

package masshortcut

// #cgo CFLAGS: -x objective-c -framework Cocoa -framework CoreFoundation -framework Carbon -framework AppKit
// #cgo LDFLAGS: -framework Cocoa -framework CoreFoundation -framework Carbon -framework AppKit
// #import "./MASShortcutMonitor.h"
// void CHotkeyCallback(int Keys, int Modifiers);
// void register_shortcut(int Keys, int Modifiers);
import "C"
import (
	"sync"
	"unsafe"
)

var (
	cbMap               = map[string]func(){}
	currentShortcuts    = []string{}
	currentShortcutLock = sync.Mutex{}
)

// Creates a hash.
func hashgen(Keys, Modifiers int) string {
	a := make([]byte, 16)
	for i, v := range *(*[8]byte)(unsafe.Pointer(&Keys)) {
		a[i] = v
	}
	for i, v := range *(*[8]byte)(unsafe.Pointer(&Modifiers)) {
		a[i+8] = v
	}
	return string(a)
}

// CHotkeyCallback is the hotkey callback for the shortcuts.
//export CHotkeyCallback
func CHotkeyCallback(CKeys, CModifiers C.int) {
	Keys := int(CKeys)
	Modifiers := int(CModifiers)
	s, ok := cbMap[hashgen(Keys, Modifiers)]
	if !ok {
		return
	}
	go s()
}

// RegisterShortcut is used to register a shortcut.
func RegisterShortcut(Keys, Modifiers int, Callback func()) {
	currentShortcutLock.Lock()
	includes := false
	hash := hashgen(Keys, Modifiers)
	for _, v := range currentShortcuts {
		if v == hash {
			includes = true
			break
		}
	}
	cbMap[hash] = Callback
	currentShortcutLock.Unlock()
	if !includes {
		C.register_shortcut(C.int(Keys), C.int(Modifiers))
	}
}

// UnregisterShortcuts is used to unregister all of the shortcuts.
func UnregisterShortcuts() {
	currentShortcutLock.Lock()
	cbMap = map[string]func(){}
	currentShortcutLock.Unlock()
}
