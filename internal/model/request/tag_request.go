package request

// CreateTagReq is the request model for creating a tag
type CreateTagReq struct {
	Name string `json:"name" binding:"required"`
}

// AddTagsReq is the request model for adding tags to an image
type AddTagsReq struct {
	Tags []string `json:"tags" binding:"required"`
}
