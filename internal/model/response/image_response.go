package response

// ImageUploadResponse defines the upload image response format
type ImageUploadResponse struct {
	ID   uint   `json:"id" example:"1"`
	Path string `json:"path" example:"abc123.jpg"`
}

// ImageDeleteResponse defines the image deletion response format
type ImageDeleteResponse struct {
	Message string `json:"message" example:"Image deleted successfully"`
	ID      uint   `json:"id" example:"1"`
}

// ImageCountResponse defines the response for image count
type ImageCountResponse struct {
	Count int64 `json:"count" example:"42"`
}
