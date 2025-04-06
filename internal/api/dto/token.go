package dto

// структура запроса на обновление refresh токена
type RefreshTokenUpdate struct {
	RefreshToken string `json:"refresh_token"`
}
