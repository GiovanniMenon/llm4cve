package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	individual bool
	verbose    bool
	rootCmd    = &cobra.Command{
		Use:   "llm4cve [CVE_ID]",
		Short: "llm4cve is a CLI tool that analyzes and summarizes CVEs using local LLMs.",
		Long:  "llm4cve is a CLI tool that analyzes and summarizes CVEs using local LLMs.",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				logrus.Fatalln("Missing argument. Please provide at least one argument between CVE-ID and CAPEC-ID.")
			}
			// var responses []string
			// for _, arg := range args {
			// 	if match, _ := regexp.MatchString(`^CVE-\d{4}-\d{4,}$`, args[0]); !match {
			// 		logrus.Fatalln("Invalid CVE-ID format. Please provide a valid CVE-ID.")
			// 	}

			// var url = "https://cve.circl.lu/api/vulnerability/" + arg

			// logrus.Debug(url)

			// resp, err := http.Get(url)
			// if err != nil {
			// 	logrus.Errorf("HTTP request failed for %s: %v", arg, err)
			// 	continue
			// }
			// defer resp.Body.Close()
			// if resp.StatusCode != http.StatusOK {
			// 	logrus.Errorf("Failed request for %s: %v", arg, resp.Status)
			// 	continue
			// }

			// bodyBytes, err := ioutil.ReadAll(resp.Body)
			// if err != nil {
			// 	logrus.Errorf("Failed reading response body for %s: %v", arg, err)
			// 	continue
			// }

			// responses = append(responses, string(bodyBytes))

			// 	}
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

	logrus.Debug("Checking CVEs Database")
	if _, err := os.Stat("cves"); os.IsNotExist(err) {
		logrus.Warn("Database not found")

		resp, err := http.Get("https://api.github.com/repos/CVEProject/cvelistV5/releases/latest")
		if err != nil {
			logrus.Fatalf("API error: %s", err)
		}
		defer resp.Body.Close()

		var release Release
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			log.Fatalf("Errore parsing JSON: %v", err)
		}

		logrus.Debugf("Downloading last CVEs database: %s", release.TagName)
		logrus.Info("Selected Asset: ", release.Assets[0].BrowserDownloadURL)

		// Missing logic for downloading -- Using wget
		wget := exec.Command("wget", "-q", "--show-progress", release.Assets[0].BrowserDownloadURL)
		wget.Stdout = os.Stdout
		wget.Stderr = os.Stderr
		if err := wget.Run(); err != nil {
			log.Fatalf("Error downloading Database %s ", err)
		}

		// Unzip

		unzip := exec.Command("unzip", release.Assets[0].Name)
		unzip.Stdout = os.Stdout
		unzip.Stderr = os.Stderr
		if err := unzip.Run(); err != nil {
			log.Fatalf("Error downloading Database %s ", err)
		}

	}

	logrus.Info("Initialized project")

	logrus.Fatal("Test")

}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}
