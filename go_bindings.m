#include "_cgo_export.h"

void unload_all_shortcuts() {
    [[MASShortcutMonitor sharedMonitor] unregisterAllShortcuts];
}

void register_shortcut(int Keys, int Modifiers, int HotkeyID) {
    MASShortcut* shortcut = [MASShortcut shortcutWithKeyCode:Keys modifierFlags:Modifiers];
    [[MASShortcutMonitor sharedMonitor] registerShortcut:shortcut withAction:^(void) { CHotkeyCallback(HotkeyID); }];
}
