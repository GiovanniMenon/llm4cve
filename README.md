# LLM4CVE

**LLM4CVE** is a command line tool that use LLM
to analyze and summarize **C**ommon **V**ulnerabilities and **E**xposures (CVEs).
Given a CVE (or a list) the system generate a summary of their description.

### Set up

```go
// go run main.go -h
`  _ _           _  _
  | | |_ __ ___ | || |   _____   _____
  | | | '_ ' _ \| || |_ / __\ \ / / _ \
  | | | | | | | |__   _| (__ \ V /  __/
  |_|_|_| |_| |_|  |_|  \___| \_/ \___|`
llm4cve is a CLI tool that analyzes and summarizes CVEs using local LLMs.

Usage:
llm4cve [CVE_ID] [flags]

Flags:
-h, --help                help for llm4cve
-m, --model string        Chose LLM model for analysis ['llama3.2','deepseek-r1:14b'] (default "deepseek-r1:14b")
-u, --ollama-url string   Use custom URL for Ollama API (default "http://127.0.0.1:11434")
-o, --output string       Save output to file. Output will be in Markdown
-v, --verbose             Display additional information
```

Example of usage :

```bash
# Run
go run main.go CVE-2025-29927 CVE-1999-1000 -v -o -u 'http://your-ollama-url:11434'
# Output produced is in output.md

# Build
go build
./llm4cve CVE-2025-29927 CVE-1999-1000 -v
```

### Tools and Technologies

-   Langchain Framework
-   Ollama _(DeepSeek, llama 3.2)_
-   Go

### Future Work

-   [x] Custom URL for Ollama as flag
-   [x] Custom output file as flag
-   [x] Model choice as flag
-   [ ] Add Support for CWE, Capec
-   [ ] Add different Database Source

### Authors

```
@GiovanniMenon
@NicoloPellegrinelli
```
