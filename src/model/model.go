package model

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"slices"

	"github.com/charmbracelet/glamour"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const SystemPrompt = `
Carefully analyze the content of the following text related to single or multiple CVE (Common Vulnerabilities and Exposures). Ignore the structure or formatting of the JSON itself and focus entirely on the technical and semantic details of the vulnerability.

Provide a clear, concise, and technically-oriented summary that includes:

    A natural language description of the vulnerability

    The affected products and versions

    The vulnerability mechanism (e.g., race condition, buffer overflow, etc.)

    The security impact (e.g., privilege escalation, remote code execution, auth bypass, etc.)

    The severity level (e.g., CVSS score, textual severity if available)

    Any known mitigations or patches

    Any additional details relevant to a security analyst

Do not describe the JSON structure or include phrases like “this JSON represents…”. Focus strictly on the CVE content.
`

var ollamaURL string

var analysisModel string

var availableModels = []string{}

func SetURL(u string) error {
	parsedURL, err := url.ParseRequestURI(u)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("Invalid URL format. Please provide a valid URL.")
	}

	ollamaURL = parsedURL.String()

	err = GetAvailableModels()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func GetAvailableModels() error {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(ollamaURL + "/api/tags")
	if err != nil {
		return fmt.Errorf("Failed to connect to LLM API ")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to get available models: %s ", resp.Status)
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	availableModels = make([]string, 0, len(result.Models))
	for _, model := range result.Models {
		availableModels = append(availableModels, model.Name)
	}

	return nil
}

func IsModelAvailable(model string) bool {
	if !strings.Contains(model, ":") {
		model = model + ":latest"
	}
	return slices.Contains(availableModels, model)
}

func SetAnalysisModel(model string) bool {
	if IsModelAvailable(model) {
		analysisModel = model
		return true
	}
	return false
}

func CreateLLM(model string) (*ollama.LLM, error) {
	llm, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(ollamaURL),
	)
	return llm, err
}

func Analysis(cves []string, outputFile string) error {
	llm, err := CreateLLM(analysisModel)
	if err != nil {
		return fmt.Errorf("failed to contact LLM: %w ", err)
	}

	ctx := context.Background()
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, SystemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, strings.Join(cves, "\n------------------------------\n")),
	}

	var buffer strings.Builder
	text := ""
	_, err = llm.GenerateContent(ctx, content,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			buffer.Write(chunk)
			text += string(chunk)
			for {
				content := buffer.String()
				if idx := strings.Index(content, "\n\n"); idx != -1 {
					paragraph := content[:idx+2]
					buffer.Reset()
					buffer.WriteString(content[idx+2:])

					rendered, err := glamour.Render(paragraph, "dracula")
					if err != nil {
						fmt.Print(paragraph)
					} else {
						cleaned := strings.ReplaceAll(rendered, "\n\n", "\n")
						fmt.Print(cleaned)
						// fmt.Print(rendered)
					}
				} else {
					break
				}
			}

			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("LLM call failed: %w", err)
	}

	// Handle any remaining text after the stream ends
	if remaining := buffer.String(); len(remaining) > 0 {
		rendered, err := glamour.Render(remaining, "dracula")
		if err != nil {
			fmt.Print(remaining)
		} else {
			fmt.Print(rendered)
		}
	}

	if outputFile != "" {
		return WriteToFile(outputFile, text)
	}
	return nil
}

func WriteToFile(outputFile string, text string) error {
	logrus.Debug("Writing output file")		
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("Error creating file: %w ", err)
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return fmt.Errorf("Error writing to file: %w ", err)
	}

	logrus.Infoln("Output file:", outputFile)
	return nil
}

func Summary(cve string) (string, error) {
	llm, err := CreateLLM("llama3.2")
	if err != nil {
		return "", fmt.Errorf("failed to contact LLM: %w ", err)
	}

	ctx := context.Background()

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Your job is to summarize the given JSON vulnerability in a single continuous text string. Do not exclude or alter any data. Do not separate into paragraphs or bullet points. The result must be a flat, plain-text sentence-style summary containing every information."),
		llms.TextParts(llms.ChatMessageTypeHuman, cve),
	}
	text := ""
	completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		text += string(chunk)
		return nil
	}))
	_ = completion

	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w ", err)
	}

	return text, nil
}
