package model

import (
	"testing"
)

func TestAnalysis(t *testing.T) {
	responses := []string{
		"CVE-2023-XXXX: Buffer overflow in XYZ software can lead to remote code execution.",
		"CVE-2023-YYYY: Improper input validation in ABC library allows privilege escalation.",
	}

	err := Analysis(responses)
	if err != nil {
		t.Fatalf("Summarize failed: %v", err)
	}
}
