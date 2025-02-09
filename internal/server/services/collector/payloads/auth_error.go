package payloads

type AuthError struct{}

func NewAuthErrorPayload() *AuthError {
	return &AuthError{}
}

func (p AuthError) String() string {
	return "Ваши учетные данные не актуальны. Пожалуйста обновите учетные данные в приложении"
}
