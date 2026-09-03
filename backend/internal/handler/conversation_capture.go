package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/Wei-Shaw/sub2api/internal/securityaudit"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type ConversationCapture struct {
	monitor *service.ConversationMonitorService
}

func NewConversationCapture(monitor *service.ConversationMonitorService) *ConversationCapture {
	return &ConversationCapture{monitor: monitor}
}

type conversationWriter struct {
	gin.ResponseWriter
	body      bytes.Buffer
	limit     int
	truncated bool
}

func (w *conversationWriter) Write(p []byte) (int, error) {
	if w.body.Len() < w.limit {
		n := w.limit - w.body.Len()
		if len(p) > n {
			w.body.Write(p[:n])
			w.truncated = true
		} else {
			w.body.Write(p)
		}
	} else {
		w.truncated = true
	}
	return w.ResponseWriter.Write(p)
}
func (w *conversationWriter) WriteString(s string) (int, error) { return w.Write([]byte(s)) }

func isConversationPath(path string) bool {
	return strings.HasSuffix(path, "/messages") || strings.HasSuffix(path, "/chat/completions") || strings.HasSuffix(path, "/responses") || strings.Contains(path, "/models/") && strings.Contains(path, ":generateContent")
}

func requestProtocol(path string) string {
	switch {
	case strings.HasSuffix(path, "/messages"):
		return "anthropic_messages"
	case strings.HasSuffix(path, "/chat/completions"):
		return "openai_chat_completions"
	case strings.HasSuffix(path, "/responses"):
		return "openai_responses"
	default:
		return "gemini"
	}
}

func jsonString(root map[string]any, key string) string {
	if v, ok := root[key].(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}

func extractModel(body []byte) string {
	var root map[string]any
	if json.Unmarshal(body, &root) != nil {
		return ""
	}
	return jsonString(root, "model")
}

func trimConversationText(value string, max int) (string, bool) {
	if max <= 0 || len(value) <= max {
		return value, false
	}
	b := []byte(value[:max])
	for len(b) > 0 && (b[len(b)-1]&0xc0) == 0x80 {
		b = b[:len(b)-1]
	}
	return string(b), true
}

var conversationSecretPattern = regexp.MustCompile(`(?i)(bearer\s+|(?:api[_-]?key|token|secret|password)[\s:=]+)[A-Za-z0-9._~+\-/]{8,}`)

func redactConversationText(value string) string {
	return conversationSecretPattern.ReplaceAllString(value, "$1[REDACTED]")
}

func textFromJSON(value any) string {
	root, ok := value.(map[string]any)
	if !ok {
		return ""
	}
	if s := jsonString(root, "output_text"); s != "" {
		return s
	}
	if delta, ok := root["delta"].(map[string]any); ok {
		if s := jsonString(delta, "text"); s != "" {
			return s
		}
		if s := jsonString(delta, "content"); s != "" {
			return s
		}
	}
	if content, ok := root["content"].([]any); ok {
		var out strings.Builder
		for _, item := range content {
			if m, ok := item.(map[string]any); ok {
				if s := jsonString(m, "text"); s != "" {
					out.WriteString(s)
				}
			}
		}
		if out.Len() > 0 {
			return out.String()
		}
	}
	if choices, ok := root["choices"].([]any); ok && len(choices) > 0 {
		if c, ok := choices[0].(map[string]any); ok {
			if m, ok := c["message"].(map[string]any); ok {
				if s := jsonString(m, "content"); s != "" {
					return s
				}
				if s := textFromJSON(m); s != "" {
					return s
				}
			}
			if d, ok := c["delta"].(map[string]any); ok {
				if s := jsonString(d, "content"); s != "" {
					return s
				}
				if s := jsonString(d, "text"); s != "" {
					return s
				}
			}
		}
	}
	if candidates, ok := root["candidates"].([]any); ok && len(candidates) > 0 {
		if c, ok := candidates[0].(map[string]any); ok {
			if content, ok := c["content"].(map[string]any); ok {
				return textFromJSON(content)
			}
		}
	}
	if parts, ok := root["parts"].([]any); ok {
		var out strings.Builder
		for _, p := range parts {
			if m, ok := p.(map[string]any); ok {
				out.WriteString(jsonString(m, "text"))
			}
		}
		return out.String()
	}
	return ""
}

func extractResponseText(body []byte, stream bool) string {
	if !stream {
		var v any
		if json.Unmarshal(body, &v) == nil {
			return textFromJSON(v)
		}
		return ""
	}
	var out strings.Builder
	for _, line := range bytes.Split(body, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, []byte("data:")) {
			continue
		}
		data := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data:")))
		if bytes.Equal(data, []byte("[DONE]")) {
			continue
		}
		var v any
		if json.Unmarshal(data, &v) == nil {
			out.WriteString(textFromJSON(v))
		}
	}
	return out.String()
}

func (m *ConversationCapture) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m == nil || m.monitor == nil || c.Request == nil || c.Request.Method != http.MethodPost || !isConversationPath(c.Request.URL.Path) {
			c.Next()
			return
		}
		cfg, err := m.monitor.GetConfig(c.Request.Context())
		if err != nil || !cfg.Enabled {
			c.Next()
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		apiKey, ok := middleware.GetAPIKeyFromContext(c)
		if !ok || apiKey == nil || apiKey.User == nil {
			c.Next()
			return
		}
		requestID, _ := c.Request.Context().Value(ctxkey.ClientRequestID).(string)
		promptReq := securityaudit.Request{Body: body, Protocol: requestProtocol(c.Request.URL.Path), RequestID: requestID, UserID: apiKey.UserID, Username: apiKey.User.Username, UserEmail: apiKey.User.Email, APIKeyID: apiKey.ID, APIKeyName: apiKey.Name, GroupID: apiKey.GroupID, Provider: "", Endpoint: c.Request.URL.Path, Model: extractModel(body)}
		if apiKey.Group != nil {
			promptReq.Provider = apiKey.Group.Platform
			promptReq.GroupName = apiKey.Group.Name
		}
		snapshot, _ := securityaudit.ExtractPromptSnapshot(promptReq)
		stream := strings.Contains(strings.ToLower(c.GetHeader("Accept")), "text/event-stream") || bytes.Contains(body, []byte(`"stream":true`))
		cw := &conversationWriter{ResponseWriter: c.Writer, limit: cfg.MaxResponseBytes}
		c.Writer = cw
		started := time.Now()
		c.Next()
		prompt, promptTruncated := trimConversationText(snapshot.FullPrompt, cfg.MaxPromptBytes)
		prompt = redactConversationText(prompt)
		responseText := extractResponseText(cw.body.Bytes(), stream)
		responseText, responseTruncated := trimConversationText(responseText, cfg.MaxResponseBytes)
		responseText = redactConversationText(responseText)
		complete := c.Writer.Status() >= 200 && c.Writer.Status() < 300
		status := http.StatusText(c.Writer.Status())
		if status == "" {
			status = "unknown"
		}
		if !complete {
			status = "failed"
		}
		input := service.ConversationRecordInput{RequestID: promptReq.RequestID, UserID: apiKey.UserID, APIKeyID: &apiKey.ID, GroupID: apiKey.GroupID, Provider: promptReq.Provider, Endpoint: promptReq.Endpoint, RequestedModel: promptReq.Model, PromptContent: prompt, ResponseContent: responseText, Stream: stream, Status: status, Complete: complete, Truncated: promptTruncated || responseTruncated || cw.truncated, DurationMS: time.Since(started).Milliseconds()}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.monitor.Record(ctx, input)
	}
}
