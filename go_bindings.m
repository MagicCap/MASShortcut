#include "_cgo_export.h"

MASShortcutMonitor* get_default_shortcut_monitor() {
    return [MASShortcutMonitor sharedMonitor];
}

void unload_all_shortcuts(MASShortcutMonitor* monitor) {
    [monitor unregisterAllShortcuts];
}

void register_shortcut(MASShortcutMonitor* monitor, int Keys, int Modifiers, int HotkeyID) {
    MASShortcut* shortcut = [MASShortcut shortcutWithKeyCode:Keys modifierFlags:Modifiers];
    [monitor registerShortcut:shortcut withAction:^(void) { CHotkeyCallback(HotkeyID); }];
}
