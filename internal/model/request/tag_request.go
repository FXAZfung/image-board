package request

// CreateTagReq is the request model for creating a tag
type CreateTagReq struct {
	Name string `json:"name" binding:"required"`
}

// AddTagReq represents a request to add a tag to an image
type AddTagReq struct {
	ID  uint   `json:"id" binding:"required"`
	Tag string `json:"tag" binding:"required"`
}

type RemoveTagReq struct {
	ImageID uint `json:"image_id" binding:"required"`
	TagID   uint `json:"tag_id" binding:"required"`
}

// AddTagsReq is the request model for adding tags to an image
type AddTagsReq struct {
	Tags []string `json:"tags" binding:"required"`
}
