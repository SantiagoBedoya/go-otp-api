package model

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	Secret2FA string `json:"secret_2fa,omitempty"`
	Valid2FA  bool   `json:"valid_2fa,omitempty"`
}
