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
)

var (
	cve = &cobra.Command{
		Use:   "cve [CVE-ID]",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				logrus.Fatalln("Missing CVE-ID argument. Please provide at least one CVE-ID.")
			}
			logrus.Infoln("CVEs Provided: ", args)

			llm, err := ollama.New(
				ollama.WithModel("deepseek-r1:14b"),
				ollama.WithServerURL(""),
				ollama.WithSystemPrompt(`You are an assistant that receives CVEs and summarizes them,
providing important information about the vulnerability.`),
			)
			if err != nil {
				log.Fatal(err)
			}

			for _, cveID := range args {
				if match, _ := regexp.MatchString("^CVE-\\d{4}-\\d{4,}$", cveID); !match {
					logrus.Warnf("Skipping invalid CVE-ID: %s", cveID)
					continue
				}

				logrus.Debug("Fetching details for CVE-ID: %s", cveID)

				resp, err := http.Get("https://cve.circl.lu/api/vulnerability/" + cveID)
				if err != nil {
					logrus.Errorf("Failed to fetch CVE %s: %v", cveID, err)
					continue
				}
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						fmt.Println(err)
					}
				}(resp.Body)

				if resp.StatusCode != http.StatusOK {
					logrus.Errorf("Failed to fetch CVE %s: %v", cveID, resp.Status)
					continue
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					logrus.Errorf("Failed to read response for CVE %s: %v", cveID, err)
					continue
				}

				// Move this in another file after merge
				_, err = llm.Call(nil, string(body),
					llms.WithTemperature(0.8),
					llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
						fmt.Print(string(chunk))
						return nil
					}),
				)

				if err != nil {
					logrus.Errorf("Failed to process CVE %s: %v", cveID, err)
					continue
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(cve)
}
