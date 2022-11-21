/* Copyright© 2022 kzz KManager@gmail.com*/

package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"w2w.io/cmn"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "w2w.io",
	Short: "qNear serve framework",
	Long: `qNear serve framework provide below service
		configured by .config_OSTYPE.json
		db connect
		zapLogger
		user management
		user login
		authorize/authenticate
		SMS(ALIYUN)`,
	Run: serve,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	//Settle configure parameters, db connector, logger
	cmn.Configure()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
