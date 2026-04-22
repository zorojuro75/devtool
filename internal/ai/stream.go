package ai

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type streamDelta struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// PrintStream reads an SSE stream from r and writes each token to w as it arrives.
func PrintStream(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var delta streamDelta
		if err := json.Unmarshal([]byte(data), &delta); err != nil {
			continue // skip malformed chunks
		}

		if len(delta.Choices) > 0 {
			io.WriteString(w, delta.Choices[0].Delta.Content)
		}
	}
	io.WriteString(w, "\n")
	return scanner.Err()
}