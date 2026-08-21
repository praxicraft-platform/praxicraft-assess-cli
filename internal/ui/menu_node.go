package ui

// MenuNode is one level in a nested interactive menu.
type MenuNode struct {
	Label    string
	Hint     string
	Key      string     // shortcut: "1" on home, "l" on actions
	Command  string     // leaf: CLI line to run
	Children []MenuNode // branch: open submenu
}
