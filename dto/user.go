package dto

type CreateUserReq struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateUserRes struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginUserReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	AccessToken string
	ID          string `json:"id"`
	Username    string `json:"username"`
}
