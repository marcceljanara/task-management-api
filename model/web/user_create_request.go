package web

type UserCreateRequest struct {
	Name string `validate:"required,min=1,max=200" json:"name"`
	Email string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=8" json:"password"`
}