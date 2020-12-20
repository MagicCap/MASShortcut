#include "_cgo_export.h"

void unload_all_shortcuts() {
    [[MASShortcutMonitor sharedMonitor] unregisterAllShortcuts];
}

void register_shortcut(int Keys, int Modifiers) {
    MASShortcut* shortcut = [MASShortcut shortcutWithKeyCode:Keys modifierFlags:Modifiers];
    [[MASShortcutMonitor sharedMonitor] registerShortcut:shortcut withAction:^(void) { CHotkeyCallback(Keys, Modifiers); }];
}

void unload_shortcut(MASShortcut* shortcut) {
    [[MASShortcutMonitor sharedMonitor] unregisterShortcut:shortcut];
}
