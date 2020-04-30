// +build darwin

package masshortcut

// #cgo CFLAGS: -x objective-c -framework Cocoa -framework CoreFoundation -framework Carbon -framework AppKit
// #cgo LDFLAGS: -framework Cocoa -framework CoreFoundation -framework Carbon -framework AppKit
// #import "./MASShortcutMonitor.h"
// void CHotkeyCallback(int CIndex);
// void unload_all_shortcuts();
// void register_shortcut(int Keys, int Modifiers, int HotkeyID);
import "C"
import "sync"

var (
	cbMap               = map[int]func(){}
	currentShortcut     = 0
	currentShortcutLock = sync.Mutex{}
)

// CHotkeyCallback is the hotkey callback for the shortcuts.
//export CHotkeyCallback
func CHotkeyCallback(CIndex C.int) {
	Index := int(CIndex)
	s, ok := cbMap[Index]
	if !ok {
		return
	}
	go s()
}


// RegisterShortcut is used to register a shortcut.
func RegisterShortcut(Keys, Modifiers int, Callback func()) {
	currentShortcutLock.Lock()
	ThisID := currentShortcut
	currentShortcut++
	cbMap[ThisID] = Callback
	currentShortcutLock.Unlock()
	C.register_shortcut(C.int(Keys), C.int(Modifiers), C.int(ThisID))
}

// UnregisterShortcuts is used to unregister all of the shortcuts.
func UnregisterShortcuts() {
	currentShortcutLock.Lock()
	C.unload_all_shortcuts()
	cbMap = map[int]func(){}
	currentShortcutLock.Unlock()
}
