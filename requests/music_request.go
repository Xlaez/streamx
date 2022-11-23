package requests

type UploadMusic struct {
	// Id        primitive.ObjectID `json:"_id" required:"true"`
	Title  string `form:"title" required:"true"`
	Cover  string `form:"cover" required:"true"`
	Artist string `form:"artist" required:"true"`
}

type GetOneMusic struct {
	Id string `uri:"id" required:"true"`
}
