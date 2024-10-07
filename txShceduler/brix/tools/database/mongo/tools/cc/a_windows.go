//go:build windows
// +build windows

package cc

import (
	"os"
	"runtime"

	"golang.org/x/sys/windows"
)

func init() {
	//https://docs.microsoft.com/en-us/windows/console/setconsolemode
	if runtime.GOOS == "windows" {
		stdout := windows.Handle(os.Stdout.Fd())
		var originalMode uint32
		windows.GetConsoleMode(stdout, &originalMode)
		windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}

}
