package panik

import (
	"fmt"
	"os"
	"runtime"

	"github.com/whytehack/goftp/pkg/logs"
)

// Catch catches panics if any and terminate program with more information in addition to stack traces.
// Call Catch with defer at the start of main and goroutines.
func Catch() {
	if err := recover(); err != nil {
		logs.ERROR.Printf("panik.Catch: caught a panic: %v \n", err)
		// following loop is inspired by debug.Stack() in addition to that this gathers other goroutines stack traces.
		stack := make([]byte, 1024)
		for {
			n := runtime.Stack(stack, true)
			if n < len(stack) {
				stack = stack[:n]
				break
			}
			stack = make([]byte, 2*len(stack))
		}
		panicLogger(err, stack)
	}
}

func panicLogger(err interface{}, stack []byte) {
	panicInfo := fmt.Sprintf(`
    RecoveredPanic: %v
================== Stack Start ==================
%s
==================  Stack End  ==================
`, err, stack)

	// Make sure we output panik info to console

	logs.ERROR.Println(panicInfo + "\r\n")

	_ = os.Stderr.Sync()
	_ = os.Stdout.Sync()
	osExit(ExitPanic)
}

const ExitPanic = 2

var osExit = os.Exit
