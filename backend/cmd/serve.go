/* Copyright©2022 kzz KManager@gmail.com */

package cmd

import (
	"fmt"
	"w2w.io/cmn"
	"w2w.io/sckserve"
	"w2w.io/service"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "service hub",
	Long:  `services starter.`,
	Run:   serve,
}

func serve(cmd *cobra.Command, args []string) {
	go service.WebServe(cmd, args)
	go sckserve.SocketServe(cmd, args)

	z.Info("serve gate started")
	_ = <-cmn.GetTerminateSignal()
	fmt.Println("service stopped")
}

func init() {
	rootCmd.AddCommand(serveCmd)

}
