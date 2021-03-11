package service

// error codes
const (
	ErrOK = iota
	ErrValidationRequest
	ErrValidationFieldIsTooLarge
	ErrValidationFieldIsTooSmall
	ErrValidationFieldIsNotRectangular
	ErrValidationFieldHasInvalidData
	ErrStorageFailed
)
