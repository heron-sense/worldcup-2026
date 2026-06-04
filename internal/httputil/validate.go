package httputil

import (
	"encoding/json"
	"net/http"
	"regexp"
)

var namePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// ValidName 校验 scene/slot 名称：非空，仅字母数字及 . - _
func ValidName(name string) bool {
	return name != "" && namePattern.MatchString(name)
}

// WriteJSON 写出 JSON 响应
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError 写出错误 JSON
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}
