// +build windows

package process

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"
)

func (p process) Signal(sig syscall.Signal) {
	err:= p.Process.SendSignal(windows.Signal(sig))

	if err != nil {
		fmt.Printf("err [%s] signal [%d]", err.Error(), sig)
	}
}
