package response

// TagDeleteResponse is the response when deleting a tag
type TagDeleteResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
}

// ImageTagResponse is the response for tag operations on images
type ImageTagResponse struct {
	ImageID uint   `json:"image_id"`
	Success bool   `json:"success"`
	Tags    []uint `json:"tags,omitempty"`
}
