package service

type JWTService interface {
	GenerateToken(userId string) (string, error) // JWT
	ValidateToken(tokenString string) (string, error) // userId
}