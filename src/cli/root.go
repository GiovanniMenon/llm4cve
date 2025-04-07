package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initProject)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initProject() {
	fmt.Println("\033[34m" + `
 _ _           _  _                  
| | |_ __ ___ | || |   _____   _____ 
| | | '_ ' _ \| || |_ / __\ \ / / _ \
| | | | | | | |__   _| (__ \ V /  __/
|_|_|_| |_| |_|  |_|  \___| \_/ \___|
` + "\033[0m")

}
