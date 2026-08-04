package model

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TodoAddPayload struct {
	Task string `json:"task"`
}

type TodoUpdatePayload struct {
	Task *string `json:"task"`
	Done *bool   `json:"done"`
}

type ErrorResponseAPI struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
