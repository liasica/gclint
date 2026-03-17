// Package jsonmapkey contains analyzer fixtures for Chinese JSON and map key checks.
package jsonmapkey

const userNameKey = "用户名"

type payload struct {
	UserName string `json:"用户名"` // want "JSON tag key must not contain Chinese"
	Status   string `json:"status"`
}

func buildPayload() map[string]any {
	return map[string]any{
		"状态":        "active",  // want "map key must not contain Chinese"
		userNameKey: "liasica", // want "map key must not contain Chinese"
		"status":    "ok",
	}
}
