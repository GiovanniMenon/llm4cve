package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
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

func Analysis(cves []string) error {
	llm, err := ollama.New(
		ollama.WithModel("deepseek-r1:14b"),
		ollama.WithServerURL("http://172.24.1.8:11434"),
	)
	if err != nil {
		return fmt.Errorf("failed to contact LLM: %w", err)
	}

	ctx := context.Background()
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, SystemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, strings.Join(cves, "\n------------------------------\n")),
	}

	var buffer strings.Builder

	_, err = llm.GenerateContent(ctx, content,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			buffer.Write(chunk)

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

	return nil
}

func Summarizes(cve string) (string, error) {
	llm, err := ollama.New(
		ollama.WithModel("llama3.2"),
		ollama.WithServerURL("http://172.24.1.8:11434"),
	)
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
		//fmt.Print(string(chunk))
		text += string(chunk)
		return nil
	}))
	_ = completion

	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w ", err)
	}

	return text, nil
}
