package jsonmapkey

const userNameKey = "用户名"
const rawPayloadJSON = `{"用户名":"liasica","profile":{"状态":"active"}}` // want "JSON string key must not contain Chinese" "JSON string key must not contain Chinese"

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

func updatePayload(payloadMap map[string]any) {
	payloadMap["状态"] = "active" // want "map key must not contain Chinese"
	payloadMap["status"] = "ok"
}

func loadPayloadJSON() string {
	return rawPayloadJSON
}

func validPayloadJSON() string {
	return `{"user_name":"liasica","status":"active"}`
}
