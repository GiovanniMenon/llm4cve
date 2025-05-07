package model

import (
	"context"
	"fmt"
	"strings"

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
		ollama.WithServerURL(""),
		ollama.WithSystemPrompt(SystemPrompt),
	)
	if err != nil {
		return fmt.Errorf("failed to contact LLM: %w", err)
	}

	ctx := context.Background()
	_, err = llm.Call(ctx, strings.Join(cves, "\n"),
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("LLM call failed: %w", err)
	}

	return nil
}
