package console

import (
	"os"
	"os/exec"
	"runtime"
)

func CMDClearWnd() *CMD {
	return &CMD{
		Name:  "cls",
		HelpS: "clear console.",
		Action: func(ps []string) {
			switch runtime.GOOS {
			case "linux":
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()

			case "windows":
				cmd := exec.Command("cmd", "/c", "cls")
				cmd.Stdout = os.Stdout
				cmd.Run()
			}
		},
	}
}
