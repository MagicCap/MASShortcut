// +build darwin

package masshortcut

// #cgo CFLAGS: -x objective-c -framework Cocoa -framework CoreFoundation -framework Carbon -framework AppKit
// #cgo LDFLAGS: -framework Cocoa -framework CoreFoundation -framework Carbon -framework AppKit
// #import "./MASShortcutMonitor.h"
// void CHotkeyCallback(int CIndex);
// MASShortcutMonitor* get_default_shortcut_monitor();
// void unload_all_shortcuts(MASShortcutMonitor* monitor);
// void register_shortcut(MASShortcutMonitor* monitor, int Keys, int Modifiers, int HotkeyID);
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

// ShortcutMonitor is used to define a shortcut monitor.
type ShortcutMonitor struct {
	monitor *C.MASShortcutMonitor
	cbs     []int
	cbLock  sync.Mutex
}

// RegisterShortcut is used to register a shortcut.
func (s *ShortcutMonitor) RegisterShortcut(Keys, Modifiers int, Callback func()) {
	currentShortcutLock.Lock()
	ThisID := currentShortcut
	currentShortcut++
	cbMap[ThisID] = Callback
	currentShortcutLock.Unlock()
	s.cbLock.Lock()
	s.cbs = append(s.cbs, ThisID)
	s.cbLock.Unlock()
	C.register_shortcut(s.monitor, C.int(Keys), C.int(Modifiers), C.int(ThisID))
}

// UnregisterShortcuts is used to unregister all of the shortcuts.
func (s *ShortcutMonitor) UnregisterShortcuts() {
	s.cbLock.Lock()
	C.unload_all_shortcuts(s.monitor)
	for _, v := range s.cbs {
		delete(cbMap, v)
	}
	s.cbLock.Unlock()
}

func getGlobalShortcutMonitor() *ShortcutMonitor {
	return &ShortcutMonitor{monitor: C.get_default_shortcut_monitor(), cbs: []int{}, cbLock: sync.Mutex{}}
}

// GlobalShortcutMonitor defines the global shortcut monitor.
var GlobalShortcutMonitor = getGlobalShortcutMonitor()
