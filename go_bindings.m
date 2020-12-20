#include "_cgo_export.h"

void unload_all_shortcuts() {
    [[MASShortcutMonitor sharedMonitor] unregisterAllShortcuts];
}

MASShortcut* register_shortcut(int Keys, int Modifiers, int HotkeyID) {
    MASShortcut* shortcut = [MASShortcut shortcutWithKeyCode:Keys modifierFlags:Modifiers];
    [[MASShortcutMonitor sharedMonitor] registerShortcut:shortcut withAction:^(void) { CHotkeyCallback(HotkeyID); }];
    return shortcut;
}

void unload_shortcut(MASShortcut* shortcut) {
    [[MASShortcutMonitor sharedMonitor] unregisterShortcut:shortcut];
}
