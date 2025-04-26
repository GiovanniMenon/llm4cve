package cli

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	individual bool
	verbose    bool
	rootCmd    = &cobra.Command{
		Use:   "llm4cve [CVE_ID] [CAPEC_ID]",
		Short: "",
		Long:  "",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				logrus.Fatalln("Missing argument. Please provide at least one argument between CVE-ID and CAPEC-ID.")
			}
			var responses []string
			for _, arg := range args {

				var url string

				switch {
				case regexp.MustCompile(`^CVE-\d{4}-\d{4,}$`).MatchString(arg):
					url = "https://cve.circl.lu/api/vulnerability/" + arg

				case regexp.MustCompile(`^CAPEC-\d+$`).MatchString(arg):

					url = "https://cve.circl.lu/api/capec/" + arg
				default:
					logrus.Warnf("Skipping invalid ID: %s", arg)
					continue
				}

				logrus.Debug(url)

				resp, err := http.Get(url)
				if err != nil {
					logrus.Errorf("HTTP request failed for %s: %v", arg, err)
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					logrus.Errorf("Failed request for %s: %v", arg, resp.Status)
					continue
				}

				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					logrus.Errorf("Failed reading response body for %s: %v", arg, err)
					continue
				}

				responses = append(responses, string(bodyBytes))

			}
			llm, err := ollama.New(
				ollama.WithModel("deepseek-R1:14B"),
				ollama.WithServerURL(""),
				ollama.WithSystemPrompt(`You are an assistant that receives CVEs and CAPECs and summarizes them,
				providing important information about this in a complete summary.`),
			)
			if err != nil {
				log.Fatal(err)
			}
			ctx := context.Background()
			_, err = llm.Call(ctx, string(strings.Join(responses, "\n")),
				llms.WithTemperature(0.8),
				llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
					fmt.Print(string(chunk))
					return nil
				}),
			)

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
