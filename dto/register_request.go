package dto

type RegisterRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8"`
	Name      string   `json:"name" validate:"required"`
	Age       int      `json:"age" validate:"required,gte=18"`
	Bio       string   `json:"bio"`
	Gender    string   `json:"gender" validate:"required,oneof=male female"`
	Location  string   `json:"location"`
	Interests string   `json:"interests"`
	Photos    []string `json:"photos"`
	IsPremium bool     `json:"is_premium"`
}
