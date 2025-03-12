package params

import (
	"bytes"
	"fmt"
	"github.com/slntopp/core-chatting/cc"
	"strings"
	"text/template"
	"time"
)

type Data struct {
	Values map[string]any
}

func add(a, b int) int {
	return a + b
}

func mergeMaps(maps ...map[string]any) map[string]any {
	result := make(map[string]any)
	for _, m := range maps {
		for key, value := range m {
			result[key] = value
		}
	}
	return result
}

func ParseParameters(text string, values ...map[string]any) string {
	tmpl, err := template.New("custom").Funcs(template.FuncMap{
		"add": add,
	}).Parse(text)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return text
	}
	data := Data{
		Values: mergeMaps(values...),
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Println("error: " + err.Error())
		return text
	}
	result := buf.String()
	return strings.ReplaceAll(result, `\n`, "\n")
}

func ParametersFromChat(chat *cc.Chat) map[string]any {
	last := chat.GetMeta().GetLastMessage()
	if last == nil {
		last = &cc.Message{}
	}
	sentSecs := last.Sent / 1000
	sent := time.Unix(sentSecs, 0)

	return map[string]any{
		"CHAT_TOPIC":                chat.GetTopic(),
		"CHAT_ID":                   chat.GetUuid(),
		"CHAT_LAST_MESSAGE_MINUTES": sent.Minute(),
		"CHAT_LAST_MESSAGE_HOURS":   sent.Hour(),
		"CHAT_LAST_MESSAGE_YEAR":    sent.Year(),
		"CHAT_LAST_MESSAGE_MONTH":   int(sent.Month()),
		"CHAT_LAST_MESSAGE_DAY":     sent.Day(),
	}
}
