package response

// TagDeleteResponse is the response when deleting a tag
type TagDeleteResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
}

// ImageTagResponse represents a response for image tag operations
type ImageTagResponse struct {
	ImageID uint   `json:"image_id"`
	TagName string `json:"tag_name,omitempty"`
	Success bool   `json:"success"`
}
