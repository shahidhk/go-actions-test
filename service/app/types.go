package main

type requestBody struct {
	Input       map[string]interface{} `json:"input"`
	SessionVars map[string]interface{} `json:"session_vars"`
}

type responseBody struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}
