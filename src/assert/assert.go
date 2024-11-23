package assert

import (
	"fmt"
	"os"
	"runtime/debug"
	"syscall"
)

const (
	Panic PanicMode = iota
	Exit
	SIGTERM
)

type PanicMode int

func log(arg ...interface{}) {

	fmt.Fprintf(os.Stderr, "ASSERT: %+v\n", arg)
	fmt.Fprintln(os.Stderr, string(debug.Stack()))
}

func handlePanic(mode PanicMode) {
	switch mode {
	case Panic:
		panic("ASSERTION FAILED")
	case Exit:
		os.Exit(1)
	case SIGTERM:
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		if err != nil {
			fmt.Fprintf(os.Stderr, "PANIC: %s\n", err)
			os.Exit(1)
		}

		panic("SYSCALL PANIC")
	}
}

func ErrNotNil(err error, mode PanicMode, arg ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "PANIC: %s\n", err)
		log(arg...)
		handlePanic(mode)
	}
}

func NotNil(intf interface{}, mode PanicMode, arg ...interface{}) {
	if intf == nil {
		log(arg...)
		handlePanic(mode)
	}
}
