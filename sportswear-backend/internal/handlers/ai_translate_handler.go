package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/utils"
)

// AITranslate 批量翻译文本（OpenAI-compatible chat completions）。
// 环境变量：AI_TRANSLATE_ENDPOINT（如 https://api.openai.com/v1/chat/completions）/
// AI_TRANSLATE_KEY / AI_TRANSLATE_MODEL（默认 gpt-4o-mini）。
// 未配置时返回明确错误，前端可降级为「复制英文源」。
func AITranslate(c *gin.Context) {
	var req struct {
		Texts      []string `json:"texts" binding:"required"`
		TargetLang string   `json:"target_lang" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	endpoint := os.Getenv("AI_TRANSLATE_ENDPOINT")
	apiKey := os.Getenv("AI_TRANSLATE_KEY")
	if endpoint == "" || apiKey == "" {
		utils.BadRequest(c, "未配置 AI 翻译服务（AI_TRANSLATE_ENDPOINT / AI_TRANSLATE_KEY）")
		return
	}
	model := os.Getenv("AI_TRANSLATE_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	joined := strings.Join(req.Texts, "\n---\n")
	prompt := fmt.Sprintf("Translate the following texts into %s. Keep one translation per line, matching the order. Do not include the separator line.\n\n%s", req.TargetLang, joined)

	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a professional translator."},
			{"role": "user", "content": prompt},
		},
	}
	payload, _ := json.Marshal(body)
	httpReq, err := http.NewRequest("POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		utils.InternalError(c, "构造请求失败")
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := (&http.Client{}).Do(httpReq)
	if err != nil {
		utils.InternalError(c, "调用 AI 服务失败: "+err.Error())
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		utils.BadRequest(c, "AI 服务返回错误: "+string(respBody))
		return
	}

	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &out); err != nil || len(out.Choices) == 0 {
		utils.InternalError(c, "解析 AI 响应失败")
		return
	}

	translations := []string{}
	for _, l := range strings.Split(strings.TrimSpace(out.Choices[0].Message.Content), "\n") {
		translations = append(translations, strings.TrimSpace(l))
	}
	utils.Success(c, gin.H{"translations": translations})
}
