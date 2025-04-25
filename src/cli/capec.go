package cli

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)


var capec = &cobra.Command{
	Use: "capec [CAPEC_ID]",
	Short: "Get information about a specific CAPEC (Common Attack Pattern Enumeration and Classification) ID.",
	Long: `Get information about a specific CAPEC (Common Attack Pattern Enumeration and Classification) ID.
CAPEC is a comprehensive dictionary of attack patterns that can be used to identify and classify potential vulnerabilities in software systems.
CAPEC is maintained by the MITRE Corporation and is used by security professionals to understand and mitigate potential threats.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Fatalln("Please provide a CAPEC ID.")
		}
		
		logrus.Infoln("CVE Provided: ", args)

		if match, _ := regexp.MatchString(`^CAPEC-\d+$`, args[0]); !match {
			logrus.Fatalln("Invalid CAPEC ID format. It should be in the format CAPEC-XXXX.")
		}

		response, err := http.Get("https://cve.circl.lu/api/capec/" + args[0])

		if err != nil {
			logrus.Fatalln("Error fetching CAPEC data: ", err)
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(response.Body)

		if response.StatusCode != http.StatusOK {
			logrus.Fatalln("Error: ", response.Status)
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			logrus.Fatalln(err)
		}

		trimmed := strings.TrimSpace(string(body))
		if trimmed == "null" {
			logrus.Fatalln("No CAPEC data found for the provided ID.")
		}

		llm, err := ollama.New(
			ollama.WithModel("llama3.2"),
			ollama.WithSystemPrompt("You are an assistant that analyzes and summarizes CAPECs. Given a list (or one), generate a summary of their description"))
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
	rootCmd.AddCommand(capec)
}