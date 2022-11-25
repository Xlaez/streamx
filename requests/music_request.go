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

type GetMusicByArtist struct {
	Limit  int64  `form:"limit" binding:"required"`
	Page   int64  `form:"page" binding:"required"`
	Artist string `form:"artist" binding:"required"`
}
