// +build darwin

package masshortcut

// #cgo CFLAGS: -x objective-c -framework Cocoa -framework CoreFoundation -framework Carbon -framework AppKit
// #cgo LDFLAGS: -framework Cocoa -framework CoreFoundation -framework Carbon -framework AppKit
// #import "./MASShortcutMonitor.h"
// void CHotkeyCallback(int CIndex);
// void unload_all_shortcuts();
// MASShortcut* register_shortcut(int Keys, int Modifiers, int HotkeyID);
// void unload_shortcut(MASShortcut* shortcut);
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

// Shortcut is used to define a shortcut structure.
type Shortcut struct {
	cShortcut *C.MASShortcut
	keyId     int
}

// Unload is used to unload a shortcut.
func (s Shortcut) Unload() {
	C.unload_shortcut(s.cShortcut)
	delete(cbMap, s.keyId)
}

// RegisterShortcut is used to register a shortcut.
func RegisterShortcut(Keys, Modifiers int, Callback func()) *Shortcut {
	currentShortcutLock.Lock()
	ThisID := currentShortcut
	currentShortcut++
	cbMap[ThisID] = Callback
	currentShortcutLock.Unlock()
	return &Shortcut{cShortcut: C.register_shortcut(C.int(Keys), C.int(Modifiers), C.int(ThisID)), keyId: ThisID}
}

// UnregisterShortcuts is used to unregister all of the shortcuts.
func UnregisterShortcuts() {
	currentShortcutLock.Lock()
	C.unload_all_shortcuts()
	cbMap = map[int]func(){}
	currentShortcutLock.Unlock()
}
