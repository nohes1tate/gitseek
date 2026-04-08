package llm

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type Client interface {
	Generate(ctx context.Context, messages []*schema.Message) (string, error)
}

type client struct {
	chatModel model.ToolCallingChatModel
}

func NewLLMClient() (Client, error) {
	cm, err := NewChatModel()
	if err != nil {
		return nil, err
	}
	return &client{chatModel: cm}, nil
}

func (c *client) Generate(ctx context.Context, messages []*schema.Message) (string, error) {
	start := time.Now()

	sr, err := c.chatModel.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("chatModel.Stream failed: %w", err)
	}
	defer sr.Close()

	var builder strings.Builder
	chunkCount := 0

	for {
		chunk, err := sr.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("stream recv failed after %d chunks, elapsed=%s: %w", chunkCount, time.Since(start), err)
		}
		if chunk == nil {
			continue
		}

		chunkCount++
		builder.WriteString(chunk.Content)
	}

	text := cleanModelOutput(builder.String())
	if text == "" {
		return "", fmt.Errorf("empty llm response, chunks=%d, elapsed=%s", chunkCount, time.Since(start))
	}

	return text, nil
}

func cleanModelOutput(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)

	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "query:") {
		s = strings.TrimSpace(s[len("query:"):])
	}
	lower = strings.ToLower(s)
	if strings.HasPrefix(lower, "search query:") {
		s = strings.TrimSpace(s[len("search query:"):])
	}

	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		s = strings.TrimSpace(lines[0])
	}
	return s
}
