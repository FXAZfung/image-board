package handles

import (
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/model/request"
	"github.com/FXAZfung/image-board/internal/model/response"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ListTags lists all tags with pagination
// @Summary List tags
// @Description Get a paginated list of all tags
// @Tags tag
// @Accept json
// @Produce json
// @Param page body model.PageReq true "Pagination parameters"
// @Success 200 {object} common.PageResp{content=[]model.Tag} "Tags list and count"
// @Router /api/public/tags [post]
func ListTags(c *gin.Context) {
	var req model.PageReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	req.Validate()
	tags, total, err := op.ListTags(req.Page, req.PerPage)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, common.PageResp{
		Content: tags,
		Total:   total,
	})
}

// GetTagByID gets a tag by its ID
// @Summary Get tag by ID
// @Description Get tag details by its ID
// @Tags tag
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} model.Tag "Tag details"
// @Router /api/public/tags/{id} [get]
func GetTagByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	tag, err := op.GetTagByID(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, err)
		return
	}

	common.SuccessResp(c, tag)
}

// GetTagByName gets a tag by its name
// @Summary Get tag by name
// @Description Get tag details by its name
// @Tags tag
// @Accept json
// @Produce json
// @Param name query string true "Tag name"
// @Success 200 {object} model.Tag "Tag details"
// @Router /api/public/tags/name [get]
func GetTagByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		common.ErrorStrResp(c, http.StatusBadRequest, "Tag name is required")
		return
	}

	tag, err := op.GetTagByName(name)
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, err)
		return
	}

	common.SuccessResp(c, tag)
}

// MostPopularTags gets the most used tags
// @Summary Get popular tags
// @Description Get the most popular tags by usage count
// @Tags tag
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of tags to return" default(10)
// @Success 200 {array} model.Tag "List of popular tags"
// @Router /api/public/tags/popular [get]
func MostPopularTags(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default to 10 if invalid
	}

	tags, err := op.GetMostPopularTags(limit)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, tags)
}

// SearchTags searches tags by prefix
// @Summary Search tags
// @Description Search for tags that start with the given prefix
// @Tags tag
// @Accept json
// @Produce json
// @Param prefix query string true "Tag prefix to search for"
// @Param limit query int false "Maximum number of results" default(20)
// @Success 200 {array} model.Tag "List of matching tags"
// @Router /api/public/tags/search [get]
func SearchTags(c *gin.Context) {
	prefix := c.Query("prefix")
	if prefix == "" {
		common.ErrorStrResp(c, http.StatusBadRequest, "Search prefix is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default to 20 if invalid
	}

	tags, err := op.SearchTagsByPrefix(prefix, limit)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, tags)
}

// CreateTag creates a new tag
// @Summary Create a new tag
// @Description Create a new tag in the system
// @Tags tag
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer user token"
// @Param tag body request.CreateTagReq true "Tag to create"
// @Success 200 {object} model.Tag "Created tag"
// @Router /api/auth/tags [post]
func CreateTag(c *gin.Context) {
	var req request.CreateTagReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	if req.Name == "" {
		common.ErrorStrResp(c, http.StatusBadRequest, "Tag name is required")
		return
	}

	// Check if tag already exists
	_, err := op.GetTagByName(req.Name)
	if err == nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Tag already exists")
		return
	}

	// Create new tag
	tag := &model.Tag{
		Name:  req.Name,
		Count: 0, // Initial count is 0
	}

	if err := op.CreateTag(tag); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, tag)
}

// UpdateTag updates an existing tag
// @Summary Update a tag
// @Description Update an existing tag's information
// @Tags tag
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer user token"
// @Param id path int true "Tag ID"
// @Param tag body request.CreateTagReq true "Updated tag information"
// @Success 200 {object} model.Tag "Updated tag"
// @Router /api/auth/tags/{id} [put]
func UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req request.CreateTagReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	if req.Name == "" {
		common.ErrorStrResp(c, http.StatusBadRequest, "Tag name is required")
		return
	}

	// Get existing tag
	tag, err := op.GetTagByID(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, errs.ErrTagNotFound)
		return
	}

	// Update tag name
	tag.Name = req.Name
	if err := op.UpdateTag(tag); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, tag)
}

// DeleteTag deletes a tag
// @Summary Delete a tag
// @Description Delete a tag and remove it from all associated images
// @Tags tag
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer user token"
// @Param id path int true "Tag ID"
// @Success 200 {object} response.TagDeleteResponse "Tag deletion response"
// @Router /api/auth/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Get tag for response
	tag, err := op.GetTagByID(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, errs.ErrTagNotFound)
		return
	}

	if err := op.DeleteTag(uint(id)); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, response.TagDeleteResponse{
		ID:      uint(id),
		Name:    tag.Name,
		Success: true,
	})
}

// GetTagsByImage gets all tags for an image
// @Summary Get tags for image
// @Description Get all tags associated with a specific image
// @Tags tag
// @Accept json
// @Produce json
// @Param image_id path int true "Image ID"
// @Success 200 {array} model.Tag "List of tags"
// @Router /api/public/tags/image/{image_id} [get]
func GetTagsByImage(c *gin.Context) {
	imageIDStr := c.Param("image_id")
	imageID, err := strconv.ParseUint(imageIDStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid image ID format")
		return
	}

	// First check if image exists
	_, err = op.GetImageByID(uint(imageID))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, errs.ImageNotFound)
		return
	}

	tags, err := op.GetTagsForImage(uint(imageID))
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, tags)
}
