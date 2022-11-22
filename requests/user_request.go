package requests

type CreateUser struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=7,alphanum"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=7,alphanum"`
}

type Verfiy struct {
	Digits string `json:"digits" binding:"required,min=6,max=6"`
}

type AskToReset struct {
	Email string `uri:"email" binding:"required,email"`
}

type ResertPassword struct {
	OldPassword string `json:"old_password" binding:"required,min=7,alphanum"`
	NewPassword string `json:"new_password" binding:"required,min=7,alphanum"`
	Digits      string `json:"digits" binding:"required,min=6,max=6"`
}
