package payloads

type InternalError struct{}

func NewInternalErrorPayload() *InternalError {
	return &InternalError{}
}

func (p *InternalError) String() string {
	return "Не удалось получить актуальный наряд. Пожалуйста, проверьте его вручную"
}
