// pk1.go - OpenCode会话数据完整转存工具 (v4.0)
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	SourcePath   string
	TargetPath   string
	MessageDir   string
	PartDir      string
	DiffDir      string
	SessionDir   string
	MetadataFile string
}

type SessionGlobal struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	ProjectID string `json:"projectID"`
	Directory string `json:"directory"`
	Title     string `json:"title"`
	Time      struct {
		Created int64 `json:"created"`
		Updated int64 `json:"updated"`
	} `json:"time"`
}

type MessageSummary struct {
	Title     string `json:"title"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type MessageModel struct {
	ModelID   string `json:"modelID"`
	Provider  string `json:"provider"`
	TokenType string `json:"tokenType"`
}

type MessageTime struct {
	Created int64 `json:"created"`
}

type Message struct {
	ID      string         `json:"id"`
	Role    string         `json:"role"`
	Summary MessageSummary `json:"summary"`
	Model   MessageModel   `json:"model"`
	Time    MessageTime    `json:"time"`
	Parts   []MessagePart  `json:"parts"`
}

type MessagePart struct {
	Type   string      `json:"type"`
	Text   string      `json:"text,omitempty"`
	Title  string      `json:"title,omitempty"`
	Tool   string      `json:"tool,omitempty"`
	CallID string      `json:"callID,omitempty"`
	Status string      `json:"status,omitempty"`
	Input  interface{} `json:"input,omitempty"`
	Output interface{} `json:"output,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type DiffFile struct {
	Path          string `json:"path"`
	BeforeContent string `json:"before_content,omitempty"`
	AfterContent  string `json:"after_content,omitempty"`
	Additions     int    `json:"additions"`
	Deletions     int    `json:"deletions"`
}

type SessionDiff struct {
	Files []DiffFile `json:"files"`
}

type SessionMetadata struct {
	SessionID       string `json:"session_id"`
	Slug            string `json:"slug"`
	ProjectID       string `json:"project_id"`
	Directory       string `json:"directory"`
	Title           string `json:"title"`
	Model           string `json:"model"`
	MessageCount    int    `json:"message_count"`
	TimeRange       string `json:"time_range"`
	TimeCreated     int64  `json:"time_created"`
	TimeUpdated     int64  `json:"time_updated"`
	SourceDirectory string `json:"source_directory"`
}

type SessionData struct {
	Metadata     SessionMetadata `json:"metadata"`
	Conversation []Message       `json:"conversation"`
	Diffs        *SessionDiff    `json:"diffs,omitempty"`
}

type Metadata struct {
	LastRun  map[string]interface{} `json:"last_run"`
	Sessions map[string]SessionRec  `json:"sessions"`
}

type SessionRec struct {
	Hash         string `json:"hash"`
	OutputFile   string `json:"output_file"`
	ProcessedAt  string `json:"processed_at"`
	MessageCount int    `json:"message_count"`
}

var reUnicode = regexp.MustCompile(`\\u[0-9a-fA-F]{4}|\\x[0-9a-fA-F]{2}`)

func decodeUnicode(text string) string {
	result := reUnicode.ReplaceAllStringFunc(text, func(m string) string {
		code, _ := strconv.ParseUint(m[2:], 16, 32)
		return string(rune(code))
	})
	result = strings.ReplaceAll(result, "\\r\\n", "\n")
	result = strings.ReplaceAll(result, "\\n", "\n")
	result = strings.ReplaceAll(result, "\\r", "\n")
	result = strings.ReplaceAll(result, "\\\\n", "\n")
	result = strings.ReplaceAll(result, "\\\\r", "\n")
	result = strings.ReplaceAll(result, "\\\\r\\\\n", "\n")
	result = strings.ReplaceAll(result, "\\u003c", "<")
	result = strings.ReplaceAll(result, "\\u003e", ">")
	result = strings.ReplaceAll(result, "\\u0026", "&")
	result = strings.ReplaceAll(result, "\\u0022", `"`)
	result = strings.ReplaceAll(result, "\\u0027", "'")
	reCRLF := regexp.MustCompile(`\r\n|\r|\n`)
	result = reCRLF.ReplaceAllString(result, "\n")
	return result
}

func decodeMapUnicode(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		switch val := v.(type) {
		case string:
			result[k] = decodeUnicode(val)
		case map[string]interface{}:
			result[k] = decodeMapUnicode(val)
		case []interface{}:
			result[k] = decodeSliceUnicode(val)
		default:
			result[k] = val
		}
	}
	return result
}

func decodeSliceUnicode(data []interface{}) []interface{} {
	result := make([]interface{}, len(data))
	for i, v := range data {
		switch val := v.(type) {
		case string:
			result[i] = decodeUnicode(val)
		case map[string]interface{}:
			result[i] = decodeMapUnicode(val)
		case []interface{}:
			result[i] = decodeSliceUnicode(val)
		default:
			result[i] = val
		}
	}
	return result
}

func escapeMarkdown(text string) string {
	special := []string{"\\", "`", "*", "_", "#", "+", ".", ">", "|", "[", "]", "(", ")", "{", "}", "\""}
	for _, ch := range special {
		text = strings.ReplaceAll(text, ch, "\\"+ch)
	}
	return text
}

func highlightErrors(text string) string {
	errorPatterns := []string{
		"错误", "ERROR", "error", "Bug", "bug",
		"问题", "修复", "修改", "fix", "FIX",
		"失败", "fail", "FAIL", "异常",
		"解决", "issue", "Issue",
	}

	for _, pattern := range errorPatterns {
		if strings.Contains(text, pattern) {
			highlighted := fmt.Sprintf("<span style=\"color:red\">【%s】</span>", pattern)
			text = strings.Replace(text, pattern, highlighted, 1)
		}
	}
	return text
}

func main() {
	fmt.Println("============================================================")
	fmt.Println("         pk1 - OpenCode会话数据完整转存 (v4.0)")
	fmt.Println("============================================================")
	fmt.Println("")
	fmt.Println("提取数据类型:")
	fmt.Println("  [x] text      - 普通文本消息")
	fmt.Println("  [x] tool      - 工具调用 (完整input/output/error)")
	fmt.Println("  [x] reasoning - AI推理过程")
	fmt.Println("  [x] step      - 步骤标记")
	fmt.Println("  [x] session   - 会话上下文 (slug/projectID/directory)")
	fmt.Println("  [x] diff      - 完整代码变更 (before/after)")
	fmt.Println("")

	cfg := Config{
		SourcePath: filepath.Join(DetectHomeDir(), ".local", "share", "opencode", "storage"),
		TargetPath: func() string {
			// 自动检测PK_PATH
			if pkPath := os.Getenv("PK_PATH"); pkPath != "" {
				return pkPath
			}
			// 平台默认
			switch runtime.GOOS {
			case "windows":
				return "D:/pkskill"
			case "linux":
				return "/opt/pkskill"
			case "darwin":
				return "/Users/" + DetectUsername() + "/pkskill"
			default:
				return "D:/pkskill"
			}
		}(),
		MessageDir: "message",
		PartDir:    "part",
		DiffDir:    "session_diff",
		SessionDir: "session\\global",
	}

	// 解析++分隔符参数
	// 格式: pk1.exe / pk1.exe++<source> / pk1.exe++<source>++<target>
	if len(os.Args) > 1 {
		args := os.Args[1]
		parts := strings.Split(args, "++")
		if len(parts) >= 1 && parts[0] != "" {
			cfg.SourcePath = parts[0]
		}
		if len(parts) >= 2 && parts[1] != "" {
			cfg.TargetPath = parts[1]
		}
	}

	cfg.MessageDir = filepath.Join(cfg.SourcePath, "message")
	cfg.PartDir = filepath.Join(cfg.SourcePath, "part")
	cfg.DiffDir = filepath.Join(cfg.SourcePath, "session_diff")
	cfg.SessionDir = filepath.Join(cfg.SourcePath, "session", "global")
	cfg.MetadataFile = filepath.Join(cfg.TargetPath, "pk1_ingest_metadata.json")

	fmt.Printf("数据来源: %s\n", cfg.SourcePath)
	fmt.Printf("目标路径: %s\n", cfg.TargetPath)
	fmt.Println("")

	if _, err := os.Stat(cfg.MessageDir); os.IsNotExist(err) {
		fmt.Println("ERROR: message/ directory not found!")
		os.Exit(1)
	}
	fmt.Println("[OK] message/ directory verified")

	var metadata Metadata
	metadata.Sessions = make(map[string]SessionRec)
	if data, err := os.ReadFile(cfg.MetadataFile); err == nil {
		var loaded Metadata
		if json.Unmarshal(data, &loaded) == nil {
			metadata = loaded
			if metadata.Sessions == nil {
				metadata.Sessions = make(map[string]SessionRec)
			}
			fmt.Printf("[OK] Loaded metadata: %d sessions tracked\n", len(metadata.Sessions))
		}
	} else {
		fmt.Println("[!] Starting fresh")
	}
	fmt.Println("")

	entries, err := os.ReadDir(cfg.MessageDir)
	if err != nil {
		fmt.Printf("ERROR: Cannot read message dir: %v\n", err)
		os.Exit(1)
	}

	var sessionDirs []os.DirEntry
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "ses_") {
			sessionDirs = append(sessionDirs, e)
		}
	}

	fmt.Printf("[INFO] Found %d sessions\n", len(sessionDirs))
	fmt.Println("")

	if len(sessionDirs) == 0 {
		fmt.Println("[!] No sessions found. Exiting.")
		os.Exit(0)
	}

	var (
		processed int
		skipped   int
		updated   int
		errors    int
		startTime = time.Now()
	)

	outputDir := filepath.Join(cfg.TargetPath, "pk1_output")
	os.MkdirAll(outputDir, 0755)

	for _, sessionEntry := range sessionDirs {
		sessionID := sessionEntry.Name()
		msgDir := filepath.Join(cfg.MessageDir, sessionID)

		msgFiles, _ := filepath.Glob(filepath.Join(msgDir, "msg_*.json"))
		if len(msgFiles) == 0 {
			continue
		}

		var msgDataList []struct {
			File      string
			Content   map[string]interface{}
			Timestamp int64
		}

		for _, msgFile := range msgFiles {
			data, err := os.ReadFile(msgFile)
			if err != nil {
				continue
			}
			var content map[string]interface{}
			if json.Unmarshal(data, &content) != nil {
				continue
			}
			timeMap, ok := content["time"].(map[string]interface{})
			if !ok {
				continue
			}
			created, ok := timeMap["created"].(float64)
			if !ok {
				continue
			}

			msgDataList = append(msgDataList, struct {
				File      string
				Content   map[string]interface{}
				Timestamp int64
			}{
				File:      msgFile,
				Content:   content,
				Timestamp: int64(created),
			})
		}

		if len(msgDataList) == 0 {
			continue
		}

		sort.Slice(msgDataList, func(i, j int) bool {
			return msgDataList[i].Timestamp < msgDataList[j].Timestamp
		})

		var messages []Message
		var sessionHash string
		var firstTime, lastTime int64 = math.MaxInt64, 0
		var modelUsed, slug, projectID, directory, title string
		var timeCreated, timeUpdated int64

		for _, item := range msgDataList {
			fileHash, _ := hashFile(item.File)
			sessionHash += fileHash

			ts := item.Timestamp
			if ts < firstTime {
				firstTime = ts
			}
			if ts > lastTime {
				lastTime = ts
			}

			if modelUsed == "" {
				if modelMap, ok := item.Content["model"].(map[string]interface{}); ok {
					if m, ok := modelMap["modelID"].(string); ok {
						modelUsed = m
					}
				}
			}

			msgID := ""
			if id, ok := item.Content["id"].(string); ok {
				msgID = id
			}

			msgSummary := MessageSummary{}
			if summaryMap, ok := item.Content["summary"].(map[string]interface{}); ok {
				if s, ok := summaryMap["title"].(string); ok {
					msgSummary.Title = s
				}
				if a, ok := summaryMap["additions"].(float64); ok {
					msgSummary.Additions = int(a)
				}
				if d, ok := summaryMap["deletions"].(float64); ok {
					msgSummary.Deletions = int(d)
				}
			}

			msgModel := MessageModel{}
			if modelMap, ok := item.Content["model"].(map[string]interface{}); ok {
				if m, ok := modelMap["modelID"].(string); ok {
					msgModel.ModelID = m
				}
				if p, ok := modelMap["provider"].(string); ok {
					msgModel.Provider = p
				}
				if t, ok := modelMap["tokenType"].(string); ok {
					msgModel.TokenType = t
				}
			}

			msgTime := MessageTime{}
			if timeMap, ok := item.Content["time"].(map[string]interface{}); ok {
				if c, ok := timeMap["created"].(float64); ok {
					msgTime.Created = int64(c)
				}
			}

			role := ""
			if r, ok := item.Content["role"].(string); ok {
				role = r
			}

			msgEntry := Message{
				ID:      msgID,
				Role:    role,
				Summary: msgSummary,
				Model:   msgModel,
				Time:    msgTime,
			}

			msgEntry.Parts = extractAllParts(msgID, cfg.PartDir)

			if msgSummary.Title != "" && title == "" {
				title = msgSummary.Title
			}

			if msgEntry.Model.ModelID != "" && modelUsed == "" {
				modelUsed = msgEntry.Model.ModelID
			}

			hasContent := len(msgEntry.Parts) > 0
			isUser := role == "user"
			if hasContent || isUser {
				messages = append(messages, msgEntry)
			}
		}

		if len(messages) == 0 {
			fmt.Printf("[SKIP] %s - all messages filtered (empty)\n", sessionID)
			skipped++
			continue
		}

		sessionGlobal := readSessionGlobal(filepath.Join(cfg.SessionDir, sessionID+".json"))
		if sessionGlobal.Slug != "" {
			slug = sessionGlobal.Slug
		}
		if sessionGlobal.ProjectID != "" {
			projectID = sessionGlobal.ProjectID
		}
		if sessionGlobal.Directory != "" {
			directory = sessionGlobal.Directory
		}
		if sessionGlobal.Title != "" {
			title = sessionGlobal.Title
		}
		timeCreated = sessionGlobal.Time.Created
		timeUpdated = sessionGlobal.Time.Updated

		hashBytes := sha256.Sum256([]byte(sessionHash))
		hashStr := fmt.Sprintf("%x", hashBytes)

		if existing, ok := metadata.Sessions[sessionID]; ok && existing.Hash == hashStr {
			skipped++
			continue
		}

		if _, ok := metadata.Sessions[sessionID]; ok {
			updated++
		}

		timeRange := ""
		if len(messages) > 0 {
			start := time.Unix(0, firstTime*1000000).In(time.FixedZone("UTC+8", 8*3600))
			end := time.Unix(0, lastTime*1000000).In(time.FixedZone("UTC+8", 8*3600))
			timeRange = fmt.Sprintf("%s - %s", start.Format("2006-01-02 15:04:05"), end.Format("15:04:05"))
		}

		var diffs *SessionDiff
		diffFile := filepath.Join(cfg.DiffDir, sessionID+".json")
		if data, err := os.ReadFile(diffFile); err == nil {
			json.Unmarshal(data, &diffs)
		}

		sessionData := SessionData{
			Metadata: SessionMetadata{
				SessionID:       sessionID,
				Slug:            slug,
				ProjectID:       projectID,
				Directory:       directory,
				Title:           title,
				Model:           modelUsed,
				MessageCount:    len(messages),
				TimeRange:       timeRange,
				TimeCreated:     timeCreated,
				TimeUpdated:     timeUpdated,
				SourceDirectory: cfg.SourcePath,
			},
			Conversation: messages,
			Diffs:        diffs,
		}

		timestamp := time.Unix(0, firstTime*1000000).In(time.FixedZone("UTC+8", 8*3600)).Format("20060102_150405")
		outputName := fmt.Sprintf("conv_%s_%s", sessionID, timestamp)
		jsonPath := filepath.Join(outputDir, outputName+".json")
		mdPath := filepath.Join(outputDir, outputName+".md")

		jsonData, _ := json.MarshalIndent(sessionData, "", "  ")
		if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
			fmt.Printf("[ERROR] Failed to write %s: %v\n", sessionID, err)
			errors++
			continue
		}
		setFileTime(jsonPath)

		if err := writeMarkdown(sessionData, mdPath); err != nil {
			fmt.Printf("[ERROR] Failed to write markdown %s: %v\n", sessionID, err)
			errors++
			continue
		}

		fmt.Printf("[OK] %s\n", sessionID)

		metadata.Sessions[sessionID] = SessionRec{
			Hash:         hashStr,
			OutputFile:   filepath.Base(jsonPath),
			ProcessedAt:  time.Now().Format(time.RFC3339),
			MessageCount: len(messages),
		}

		processed++
	}

	metadata.LastRun = map[string]interface{}{
		"timestamp":   time.Now().Format(time.RFC3339),
		"source_path": cfg.SourcePath,
		"target_path": cfg.TargetPath,
		"processed":   processed,
		"updated":     updated,
		"skipped":     skipped,
		"errors":      errors,
		"version":     "v4.0",
	}

	if data, err := json.MarshalIndent(metadata, "", "  "); err == nil {
		os.WriteFile(cfg.MetadataFile, data, 0644)
		fmt.Println("[OK] Metadata saved")
	}

	// ============================================================
	// pk_core: Generate Path Config
	// ============================================================
	fmt.Println("[PK_CORE] Generating path configuration...")
	if err := GeneratePathConfig(cfg.TargetPath); err != nil {
		fmt.Printf("[WARNING] Failed to generate path config: %v\n", err)
	} else {
		fmt.Println("[PK_CORE] Path configuration saved to output/pk_path_config.json")
	}

	fmt.Println("")
	fmt.Println("============================================================")
	fmt.Println("                    INGESTION COMPLETE")
	fmt.Println("============================================================")
	fmt.Printf("  Processed: %d sessions\n", processed)
	fmt.Printf("  Updated:   %d sessions\n", updated)
	fmt.Printf("  Skipped:   %d sessions\n", skipped)
	fmt.Printf("  Errors:    %d sessions\n", errors)
	fmt.Printf("  Time:      %v\n", time.Since(startTime))
	fmt.Println("")

	if runtime.GOOS == "windows" && os.Getenv("TERM") == "" {
		fmt.Println("按 Enter 键退出...")
		fmt.Scanln()
	}
}

func hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

func readSessionGlobal(path string) SessionGlobal {
	var global SessionGlobal
	data, err := os.ReadFile(path)
	if err != nil {
		return global
	}
	json.Unmarshal(data, &global)
	return global
}

func extractAllParts(msgID, partDir string) []MessagePart {
	partSubDir := filepath.Join(partDir, msgID)
	entries, err := os.ReadDir(partSubDir)
	if err != nil {
		return nil
	}

	var parts []MessagePart

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "prt_") || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(partSubDir, entry.Name()))
		if err != nil {
			continue
		}

		var rawPart map[string]interface{}
		if json.Unmarshal(data, &rawPart) != nil {
			continue
		}

		partType := ""
		if t, ok := rawPart["type"].(string); ok {
			partType = t
		}

		part := MessagePart{
			Type: partType,
		}

		switch partType {
		case "text":
			if t, ok := rawPart["text"].(string); ok {
				part.Text = decodeUnicode(t)
			}

		case "reasoning":
			if t, ok := rawPart["text"].(string); ok {
				part.Text = decodeUnicode(t)
			}

		case "step-start":
			if t, ok := rawPart["title"].(string); ok {
				part.Title = t
			} else {
				part.Title = "未知步骤"
			}

		case "tool":
			if t, ok := rawPart["tool"].(string); ok {
				part.Tool = t
			}
			if t, ok := rawPart["callId"].(string); ok {
				part.CallID = t
			}
			if t, ok := rawPart["call_id"].(string); ok {
				part.CallID = t
			}
			if state, ok := rawPart["state"].(map[string]interface{}); ok {
				decodedState := decodeMapUnicode(state)
				if s, ok := decodedState["status"].(string); ok {
					part.Status = s
				}
				if s, ok := decodedState["title"].(string); ok {
					part.Title = s
				}
				part.Input = decodedState["input"]
				if outputRaw, ok := decodedState["output"].(string); ok {
					part.Output = outputRaw
				}
				if s, ok := decodedState["error"].(string); ok {
					part.Error = s
				}
			}
		}

		parts = append(parts, part)
	}

	return parts
}

func writeMarkdown(data SessionData, path string) error {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString("session_id: " + data.Metadata.SessionID + "\n")
	sb.WriteString("slug: " + data.Metadata.Slug + "\n")
	sb.WriteString("project_id: " + data.Metadata.ProjectID + "\n")
	sb.WriteString("directory: " + data.Metadata.Directory + "\n")
	sb.WriteString("title: " + escapeMarkdown(data.Metadata.Title) + "\n")
	sb.WriteString("model: " + data.Metadata.Model + "\n")
	sb.WriteString("message_count: " + fmt.Sprintf("%d", data.Metadata.MessageCount) + "\n")
	sb.WriteString("time_range: " + data.Metadata.TimeRange + "\n")
	sb.WriteString("---\n\n")

	sb.WriteString("# " + data.Metadata.Title + "\n\n")

	sb.WriteString("**Slug**: " + data.Metadata.Slug + " | **项目ID**: " + data.Metadata.ProjectID + "\n")
	sb.WriteString("**目录**: " + data.Metadata.Directory + "\n")
	sb.WriteString("**模型**: " + data.Metadata.Model + " | **消息数**: " + fmt.Sprintf("%d", data.Metadata.MessageCount) + "\n")
	sb.WriteString("**时间**: " + data.Metadata.TimeRange + "\n\n")

	sb.WriteString("---\n\n")

	for i, msg := range data.Conversation {
		roleChinese := "用户"
		if msg.Role == "assistant" {
			roleChinese = "助手"
		}

		sb.WriteString(fmt.Sprintf("## 消息 #%d | %s | %s\n\n", i+1, roleChinese, time.Unix(0, msg.Time.Created*1000000).In(time.FixedZone("UTC+8", 8*3600)).Format("2006-01-02 15:04:05")))

		if msg.Summary.Title != "" {
			sb.WriteString("**摘要**: " + escapeMarkdown(msg.Summary.Title) + "\n")
		}
		if msg.Model.ModelID != "" {
			sb.WriteString("**模型**: " + msg.Model.ModelID + "\n")
		}
		sb.WriteString("\n")

		for _, part := range msg.Parts {
			switch part.Type {
			case "text":
				if part.Text != "" {
					sb.WriteString(part.Text + "\n\n")
				}

			case "reasoning":
				if part.Text != "" {
					sb.WriteString("### [思考]\n\n")
					lines := strings.Split(part.Text, "\n")
					for _, line := range lines {
						line = strings.TrimSpace(line)
						if line != "" {
							sb.WriteString("> " + highlightErrors(decodeUnicode(line)) + "\n")
						}
					}
					sb.WriteString("\n")
				}

			case "tool":
				statusText := "已完成"
				if part.Status == "error" {
					statusText = "错误"
				}
				sb.WriteString(fmt.Sprintf("### [工具] %s\n\n", part.Tool))
				sb.WriteString(fmt.Sprintf("- 编号: `%s`\n", part.CallID))
				sb.WriteString(fmt.Sprintf("- 状态: %s\n", statusText))
				if part.Title != "" {
					sb.WriteString(fmt.Sprintf("- 标题: %s\n", decodeUnicode(part.Title)))
				}

				if part.Input != nil {
					inputJSON, _ := json.MarshalIndent(part.Input, "", "  ")
					decodedJSON := decodeUnicode(string(inputJSON))
					sb.WriteString(fmt.Sprintf("**输入**:\n```json\n%s\n```\n\n", decodedJSON))
				}

				if part.Error != "" {
					sb.WriteString(fmt.Sprintf("**错误**:\n```\n%s\n```\n\n", decodeUnicode(part.Error)))
				}

				if part.Output != nil {
					switch outputStr := part.Output.(type) {
					case string:
						sb.WriteString(fmt.Sprintf("**输出**:\n```\n%s\n```\n\n", decodeUnicode(outputStr)))
					default:
						outputJSON, _ := json.MarshalIndent(part.Output, "", "  ")
						decodedJSON := decodeUnicode(string(outputJSON))
						sb.WriteString(fmt.Sprintf("**输出**:\n```json\n%s\n```\n\n", decodedJSON))
					}
				}
			}
		}

		sb.WriteString("\n")
	}

	if data.Diffs != nil && len(data.Diffs.Files) > 0 {
		sb.WriteString("---\n\n## 代码变更\n\n")

		totalAdd, totalDel := 0, 0
		for _, f := range data.Diffs.Files {
			sb.WriteString(fmt.Sprintf("### %s\n", f.Path))
			sb.WriteString(fmt.Sprintf("**变更**: +%d / -%d\n\n", f.Additions, f.Deletions))

			if f.BeforeContent != "" || f.AfterContent != "" {
				sb.WriteString("```diff\n")
				if f.BeforeContent != "" {
					sb.WriteString(f.BeforeContent)
				}
				if f.AfterContent != "" {
					sb.WriteString(f.AfterContent)
				}
				sb.WriteString("\n```\n\n")
			}

			totalAdd += f.Additions
			totalDel += f.Deletions
		}

		sb.WriteString(fmt.Sprintf("**总计**: +%d / -%d\n", totalAdd, totalDel))
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

// ============================================================
// pk_core: 跨平台路径抽象层
// ============================================================

// PathConfig 路径配置结构
type PathConfig struct {
	Version     string            `json:"version"`
	GeneratedAt string            `json:"generated_at"`
	Platform    string            `json:"platform"`
	Paths       map[string]string `json:"paths"`
}

// GeneratePathConfig 生成路径配置文件
func GeneratePathConfig(pkPath string) error {
	username := DetectUsername()
	homeDir := DetectHomeDir()
	platform := DetectPlatform()

	// 标准化路径格式（统一使用正斜杠）
	opencodeStorage := filepath.ToSlash(filepath.Join(homeDir, ".local", "share", "opencode", "storage"))
	pkPathNormalized := filepath.ToSlash(pkPath)
	pk1Output := filepath.ToSlash(filepath.Join(pkPath, "pk1_output"))
	pk2Enriched := filepath.ToSlash(filepath.Join(pkPath, "pk2_enriched"))
	pk3Pool := filepath.ToSlash(filepath.Join(pkPath, "pk3_lab", "evolution_pool"))

	config := PathConfig{
		Version:     "v1",
		GeneratedAt: time.Now().Format(time.RFC3339),
		Platform:    platform,
		Paths: map[string]string{
			"pk_path":            pkPathNormalized,
			"home":               filepath.ToSlash(homeDir),
			"username":           username,
			"opencode_storage":   opencodeStorage,
			"pk1_output":         pk1Output,
			"pk2_enriched":       pk2Enriched,
			"pk3_evolution_pool": pk3Pool,
		},
	}

	// 输出到 {PK_PATH}/pk1_output/pk_path_config.json
	pk1OutputDir := pk1Output
	configPath := filepath.Join(pk1OutputDir, "pk_path_config.json")

	// 确保目录存在
	if err := os.MkdirAll(pk1OutputDir, 0755); err != nil {
		return err
	}

	// 写入JSON文件
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// DetectUsername 检测用户名
func DetectUsername() string {
	// 优先级：环境变量
	if username := os.Getenv("USERNAME"); username != "" {
		return username
	}
	if username := os.Getenv("USER"); username != "" {
		return username
	}
	if username := os.Getenv("LOGNAME"); username != "" {
		return username
	}

	// 系统API fallback
	return "unknown"
}

// DetectHomeDir 检测主目录
func DetectHomeDir() string {
	// 优先级：环境变量
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home
	}

	// Windows: HOMEDRIVE + HOMEPATH
	if homeDrive := os.Getenv("HOMEDRIVE"); homeDrive != "" {
		if homePath := os.Getenv("HOMEPATH"); homePath != "" {
			return homeDrive + homePath
		}
	}

	return "/home/" + DetectUsername()
}

// DetectPlatform 检测操作系统
func DetectPlatform() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	case "darwin":
		return "macos"
	default:
		return runtime.GOOS
	}
}

// setFileTime 设置文件修改时间为当前时间
func setFileTime(path string) error {
	now := time.Now()
	return os.Chtimes(path, now, now)
}
