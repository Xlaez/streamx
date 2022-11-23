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

type VerfiyEmail struct {
	Digits string `uri:"digits" binding:"required,min=6,max=6"`
}

type AskToReset struct {
	Email string `uri:"email" binding:"required,email"`
}

type ResertPassword struct {
	OldPassword string `json:"old_password" binding:"required,min=7,alphanum"`
	NewPassword string `json:"new_password" binding:"required,min=7,alphanum"`
	Digits      string `json:"digits" binding:"required,min=6,max=6"`
}

type UploadAvatarReq struct {
	ID string `form:"id" binding:"required"`
}

type GetUserById struct {
	ID string `uri:"id" binding:"required"`
}

type GetUsers struct {
	Limit int64 `form:"limit" binding:"required"`
	Page  int64 `form:"page" binding:"required"`
}
