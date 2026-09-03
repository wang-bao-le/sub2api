package handler

import "testing"

func TestExtractResponseText(t *testing.T) {
	tests := []struct {
		name, body, want string
		stream bool
	}{
		{"openai", `{"choices":[{"message":{"content":"hello"}}]}`, "hello", false},
		{"anthropic", `{"content":[{"type":"text","text":"hello"}]}`, "hello", false},
		{"gemini", `{"candidates":[{"content":{"parts":[{"text":"hello"}]}}]}`, "hello", false},
		{"responses", `{"output_text":"hello"}`, "hello", false},
		{"sse", "data: {\"choices\":[{\"delta\":{\"content\":\"hel\"}}]}\n\ndata: {\"choices\":[{\"delta\":{\"content\":\"lo\"}}]}\n\ndata: [DONE]\n", "hello", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractResponseText([]byte(tt.body), tt.stream); got != tt.want { t.Fatalf("extractResponseText() = %q, want %q", got, tt.want) }
		})
	}
}

func TestConversationTextSafety(t *testing.T) {
	if got := redactConversationText("Bearer abcdefghijkl token=secret-value"); got != "Bearer [REDACTED] token=[REDACTED]" { t.Fatalf("redaction = %q", got) }
	got, truncated := trimConversationText("你好世界", 7)
	if !truncated || got != "你好" { t.Fatalf("trimConversationText() = %q, %v", got, truncated) }
}
