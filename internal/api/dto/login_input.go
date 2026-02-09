package dto

type LoginInput struct {
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required"`
}
