/* Copyright © 2022 kmanager@gmail.com */
package main

import (
	"os"
	"os/signal"
	"syscall"
	"w2w.io/cmd"
	"w2w.io/cmn"
)

var terminateSignal chan os.Signal
var buildVer = "nativeBuild"

func main() {
	terminateSignal = make(chan os.Signal)
	signal.Notify(terminateSignal,
		// syscall.SIGTSTP,
		syscall.SIGILL,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGKILL,
		syscall.SIGQUIT,
		// syscall.SIGSTOP,
		syscall.SIGHUP,
	)

	go cmd.Cleanup(terminateSignal)
	cmn.SetBuildVer(buildVer)

	cmd.Execute()
}
