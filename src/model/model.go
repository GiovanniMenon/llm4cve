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
You are a cybersecurity assistant. Your task is to analyze and summarize one or more CVEs (Common Vulnerabilities and Exposures).

For each CVE provided:
- Summarize the vulnerability in plain language.
- Highlight the affected software or system.
- Describe the impact (e.g., data loss, privilege escalation).
- Mention the severity if available (e.g., CVSS score or rating).
- Note any known exploits or mitigation steps.

If multiple CVEs are provided:
- Identify any common components, attack surfaces, or affected systems.
- Assess whether the vulnerabilities could be chained or combined.
- Summarize the overall risk or attack potential if exploited together.

Be concise, accurate, and structured.
`

func Analysis(cves []string) error {
	llm, err := ollama.New(
		ollama.WithModel("llama3.2"),
		ollama.WithServerURL("http://127.0.0.1:11434"),
		ollama.WithSystemPrompt(SystemPrompt),
	)
	if err != nil {
		return fmt.Errorf("failed to contact LLM: %w", err)
	}

	ctx := context.Background()
	var buffer strings.Builder

	_, err = llm.Call(ctx, strings.Join(cves, "\n"),
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
