package cmd

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"w2w.io/cmn"
	_ "w2w.io/service"
)

var z *zap.Logger

func init() {
	//Setup package scope variables, just like logger, db connector, configure parameters, etc.
	cmn.PackageStarters = append(cmn.PackageStarters, func() {
		z = cmn.GetLogger()
		z.Info("cmd zLogger settled")
	})
}

type cleanProc func()

//cleanProcesses put cleanup job to this slice array
var cleanProcesses []cleanProc

//Cleanup executed when terminated by OS/User
func Cleanup(s chan os.Signal) {
	if s == nil {
		return
	}

	cmn.SetQuitChannel(s)

	select {
	case signal := <-s:
		fmt.Printf("\n\n")
		z.Warn(fmt.Sprintf("\n***********************************\nservice terminated by %v", signal))
		cmn.UtilCleanup()

		for _, v := range cleanProcesses {
			v()
		}

		z.Sync()
		os.Exit(0)
	}
}
