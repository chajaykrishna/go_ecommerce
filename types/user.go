package types

type UserSigninRquest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=72"`
}

type UserSignupRequest struct {
	Username string `json:"username" validate:"required,min=3,alphanum"`
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Phone    string `json:"phone" validate:"omitempty,e164"`
	Address  string `json:"address" validate:"omitempty,min=3"`
}
