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
		// handle := syscall.Handle(os.Stdout.Fd())
		// kernel32DLL := syscall.NewLazyDLL("kernel32.dll")
		// setConsoleModeProc := kernel32DLL.NewProc("SetConsoleMode")
		// _, _, err := setConsoleModeProc.Call(uintptr(handle), 0x0001|0x0002|0x0004)
		// if err != nil {
		// 	fmt.Println("[CC.COLOR] --> ", err)
		// 	fmt.Println(err.Error())

		// 	modeText := strings.ToLower(err.Error())
		// 	if strings.Contains(modeText, "incorrect") { //The parameter is incorrect.
		// 		isLegacyMode = true
		// 	}
		// }

		stdout := windows.Handle(os.Stdout.Fd())
		var originalMode uint32
		windows.GetConsoleMode(stdout, &originalMode)
		windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}

	initOS("windows")
}
