# LLM4CVE

**LLM4CVE** is a command line tool that use LLM 
to analyze and summarize **C**ommon **V**ulnerabilities and **E**xposures (CVEs).
Given a CVE (or a list) the system generate a summary of their description.
### Set up 

```go
// go run main.go -h
 _ _           _  _                  
| | |_ __ ___ | || |   _____   _____ 
| | | '_ ' _ \| || |_ / __\ \ / / _ \
| | | | | | | |__   _| (__ \ V /  __/
|_|_|_| |_| |_|  |_|  \___| \_/ \___|
llm4cve is a CLI tool that analyzes and summarizes CVEs using local LLMs.

Usage:
  llm4cve [CVE_ID] [flags]

Flags:
  -h, --help         help for llm4cve
  -o, --output       file output.md is created with the output
  -v, --verbose      Display additional information
```

Example of usage :
```go
// Run
go run main.go CVE-2025-29927 CVE-1999-1000 -v -o
// Output produced is in output.md

// Build
go build
./llm4cve CVE-2025-29927 CVE-1999-1000 -v 
```


### Tool Used

- Langchain Framework
- Ollama *(DeepSeek, llama 3.1)*
- Go 

### Future Work 

- [ ] Url, Output File Path and Model as CLI Flags;
- [ ] Add Support for CWE, Capec;
- [ ] Add different Database Source;
```
# Authors
@GiovanniMenon
@NicoloPellegrinelli
```