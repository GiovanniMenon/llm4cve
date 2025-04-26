package cli

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	individual bool
	verbose    bool
	rootCmd    = &cobra.Command{
		Use:   "llm4cve",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initProject)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Display additional information")
	rootCmd.PersistentFlags().BoolVarP(&individual, "individual", "i", false, "Compute each vulnerability individually and provide response in a file.")
}

func initProject() {
	fmt.Println("\033[34m" + `
 _ _           _  _                  
| | |_ __ ___ | || |   _____   _____ 
| | | '_ ' _ \| || |_ / __\ \ / / _ \
| | | | | | | |__   _| (__ \ V /  __/
|_|_|_| |_| |_|  |_|  \___| \_/ \___|
` + "\033[0m")

	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}
