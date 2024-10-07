package cc

import (
	"fmt"
	"testing"
)

func Test_COLOR(t *testing.T) {
	const (
		//nnn  = "\033[(n/b)"
		black  = "\033[1;30m%s\033[0m"
		red    = "\033[1;31m%s\033[0m"
		green  = "\033[1;32m%s\033[0m"
		yellow = "\033[1;33m%s\033[0m"
		blue   = "\033[1;34m%s\033[0m"
		purple = "\033[0;35m%s\033[0m"
		cyan   = "\033[0;36m%s\033[0m"
		white  = "\033[0;37m%s\033[0m"
	)
	// handle := syscall.Handle(os.Stdout.Fd())
	// kernel32DLL := syscall.NewLazyDLL("kernel32.dll")
	// setConsoleModeProc := kernel32DLL.NewProc("SetConsoleMode")
	// setConsoleModeProc.Call(uintptr(handle), 0x0001|0x0002|0x0004)
	bar := func() {
		fmt.Println("========================================")
	}
	ct := func() {
		defer bar()
		bar()
		fmt.Printf(black, "black")
		fmt.Println("")
		fmt.Printf(red, "red")
		fmt.Println("")
		fmt.Printf(green, "green")
		fmt.Println("")
		fmt.Printf(yellow, "yellow")
		fmt.Println("")
		fmt.Printf(blue, "blue")
		fmt.Println("")
		fmt.Printf(purple, "purple")
		fmt.Println("")
		fmt.Printf(cyan, "cyan")
		fmt.Println("")
		fmt.Printf(white, "white")
		fmt.Println("")
	}
	_ = ct
	ct()

	fmt.Println("\033[0;33m test_color \033[0m")        //yellow
	fmt.Println("\033[4m\033[1;36m test_color \033[0m") //cyan
	fmt.Println("\033[0;64m test_color \033[0m")        //yellow(b)
}

func Test_Loop(t *testing.T) {
	for i := 0; i < 20; i++ {
		switch i {
		case 2, 5, 6, 7, 8, 9:
			continue
		case 16, 17, 18, 19:
			continue
		}
		for j := 30; j < 50; j++ {
			s := j % 10
			if s == 8 || s == 9 {
				continue
			}
			fmt.Println("[", i, "][", j, "] ", "\033["+fmt.Sprint(i)+";"+fmt.Sprint(j)+"m abcdefg \033[0m")
		}
	}
}

func TestItalic(t *testing.T) {
	fmt.Println("\033[0;31m  TESTABCD abcdefg \033[0m")
	fmt.Println("\033[3;31m  TESTABCD abcdefg \033[0m")
}
