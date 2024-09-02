package utils

import (
	"encoding/json"
	"html"
	"strings"
)

func IsEmptyString(value string) bool {
	return len(CleanString(value)) == 0
}

func CleanString(value string) string {
	return html.EscapeString(strings.TrimSpace(value))
}

func DataString(data interface{}) string {
	out, err := json.Marshal(data)

	if err != nil {
		return ""
	}

	return string(out)
}

func MergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))

	for k, v := range a {
		out[k] = v
	}

	for k, v := range b {
		// If you use map[string]interface{}, ok is always false here.
		// Because yaml.Unmarshal will give you map[interface{}]interface{}.
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = MergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
