package game

import "fmt"

// Error describes API error response
type Error struct {
	Code    int           `json:"code"`    // error code
	Message string        `json:"message"` // default message with params substituted
	Params  []interface{} `json:"params"`  // params for i18n localized messages
}

// Error makes Error to implement error
func (e Error) Error() string { return fmt.Sprintf("%d: %s", e.Code, e.Message) }
