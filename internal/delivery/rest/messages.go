package rest

type ErrorMessage struct {
	Message string `json:"message"`
}

func FromErr(err error) ErrorMessage {
	return ErrorMessage{err.Error()}
}
func FromString(msg string) ErrorMessage {
	return ErrorMessage{msg}
}

type SuccessMessage[T any] struct {
	Data T `json:"data"`
}

func NewSuccessMessage[T any](data T) SuccessMessage[T] {
	return SuccessMessage[T]{Data: data}
}
