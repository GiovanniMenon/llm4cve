package cli

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

var cve = &cobra.Command{
	Use:   "cve [CVE-ID]",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			logrus.Fatalln("Missing CVE-ID argument. Please provide at least one CVE-ID.")
		}
		logrus.Infoln("CVE Provided: ", args)

		if match, _ := regexp.MatchString("^CVE-\\d{4}-\\d{4,}$", args[0]); !match {
			logrus.Fatalln("CVE-ID is invalid")
		}
		resp, err := http.Get("https://cve.circl.lu/api/vulnerability/" + args[0])
		if err != nil {
			logrus.Fatalln(err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			logrus.Fatalln("Failed to fetch CVEs: %s ", resp.Status)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Fatalln(err)
		}

		llm, err := ollama.New(
			ollama.WithModel("deepseek-r1:14b"),
			ollama.WithServerURL(""),
			ollama.WithSystemPrompt("You are an assistant that receives CVEs and summarize it providing important information about the vulnerability"))
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		completion, err := llm.Call(ctx, string(body),
			llms.WithTemperature(0.8),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				fmt.Print(string(chunk))
				return nil
			}),
		)
		if err != nil {
			log.Fatal(err)
		}

		_ = completion
	},
}

func init() {
	rootCmd.AddCommand(cve)
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339Nano, timeStr)
	return t
}
