package ai

import (
	"strings"
	"testing"
)

func TestPrintStream(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "single token",
			input: `data: {"choices":[{"delta":{"content":"Hello"}}]}
data: [DONE]
`,
			want: "Hello\n",
		},
		{
			name: "multiple tokens",
			input: `data: {"choices":[{"delta":{"content":"Hello"}}]}
data: {"choices":[{"delta":{"content":" world"}}]}
data: [DONE]
`,
			want: "Hello world\n",
		},
		{
			name: "skips non-data lines",
			input: `: keep-alive
data: {"choices":[{"delta":{"content":"Hi"}}]}
data: [DONE]
`,
			want: "Hi\n",
		},
		{
			name: "skips malformed json",
			input: `data: not-json
data: {"choices":[{"delta":{"content":"OK"}}]}
data: [DONE]
`,
			want: "OK\n",
		},
		{
			name:  "empty stream",
			input: "data: [DONE]\n",
			want:  "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out strings.Builder
			err := PrintStream(strings.NewReader(tt.input), &out)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if out.String() != tt.want {
				t.Errorf("output = %q, want %q", out.String(), tt.want)
			}
		})
	}
}