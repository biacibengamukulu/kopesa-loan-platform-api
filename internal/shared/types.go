package shared

import "time"

type Money int64

type Meta struct {
	RequestID  string `json:"requestId"`
	Page       int    `json:"page,omitempty"`
	PageSize   int    `json:"pageSize,omitempty"`
	Total      int    `json:"total,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type FieldError struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
}

type ErrorBody struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type ResponseEnvelope struct {
	Data  any        `json:"data"`
	Meta  Meta       `json:"meta"`
	Error *ErrorBody `json:"error"`
}

type Event struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Source        string         `json:"source"`
	SpecVersion   string         `json:"specVersion"`
	OccurredAt    time.Time      `json:"occurredAt"`
	CorrelationID string         `json:"correlationId,omitempty"`
	Actor         map[string]any `json:"actor,omitempty"`
	Data          map[string]any `json:"data,omitempty"`
}
