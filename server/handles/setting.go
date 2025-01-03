package handles

import (
	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/internal/sign"
	"github.com/FXAZfung/image-board/pkg/random"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/FXAZfung/image-board/server/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// ResetToken 重置token
// @Summary 重置token
// @Description 重置token
// @Tags setting
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Success 200 {object} string
// @Router /api/private/setting/token [post]
func ResetToken(c *gin.Context) {
	token := random.SecretKey()
	item := model.SettingItem{Key: "token", Value: token, Type: conf.TypeString, Group: model.SINGLE, Flag: model.PRIVATE}
	if err := op.SaveSettingItem(&item); err != nil {
		common.ErrorStrResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	sign.Instance()
	common.SuccessResp(c, token)
}

// GetSetting 获取设置
// @Summary 获取设置
// @Description 获取设置
// @Tags setting
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param key query string false "key"
// @Param keys query string false "keys"
// @Success 200 {object} string "设置"
// @Router /api/private/setting [get]
func GetSetting(c *gin.Context) {
	key := c.Query("key")
	keys := c.Query("keys")
	if key != "" {
		item, err := op.GetSettingItemByKey(key)
		if err != nil {
			common.ErrorStrResp(c, http.StatusBadRequest, err.Error())
			return
		}
		common.SuccessResp(c, item)
	} else {
		items, err := op.GetSettingItemInKeys(strings.Split(keys, ","))
		if err != nil {
			common.ErrorStrResp(c, http.StatusBadRequest, err.Error())
			return
		}
		common.SuccessResp(c, items)
	}
}

// SaveSettings 保存设置
// @Summary 保存设置
// @Description 保存设置
// @Tags setting
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param settings body []model.SettingItem true "设置"
// @Router /api/private/setting [post]
func SaveSettings(c *gin.Context) {
	var req []model.SettingItem
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := op.SaveSettingItems(req); err != nil {
		common.ErrorStrResp(c, http.StatusInternalServerError, err.Error())
	} else {
		common.SuccessResp(c)
		static.UpdateIndex()
	}
}

// ListSettings 列出设置
// @Summary 列出设置
// @Description 列出设置
// @Tags setting
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param group query string false "group"
// @Param groups query string false "groups"
// @Success 200 {object} []model.SettingItem "设置列表"
// @Router /api/private/settings [get]
func ListSettings(c *gin.Context) {
	groupStr := c.Query("group")
	groupsStr := c.Query("groups")
	var settings []model.SettingItem
	var err error
	if groupsStr == "" && groupStr == "" {
		settings, err = op.GetSettingItems()
	} else {
		var groupStrings []string
		if groupsStr != "" {
			groupStrings = strings.Split(groupsStr, ",")
		} else {
			groupStrings = append(groupStrings, groupStr)
		}
		var groups []int
		for _, str := range groupStrings {
			group, err := strconv.Atoi(str)
			if err != nil {
				common.ErrorStrResp(c, http.StatusBadRequest, err.Error())
				return
			}
			groups = append(groups, group)
		}
		settings, err = op.GetSettingItemsInGroups(groups)
	}
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	common.SuccessResp(c, settings)
}

// DeleteSetting 删除设置
// @Summary 删除设置
// @Description 删除设置
// @Tags setting
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param key query string true "key"
// @Router /api/private/setting [delete]
func DeleteSetting(c *gin.Context) {
	key := c.Query("key")
	if err := op.DeleteSettingItemByKey(key); err != nil {
		common.ErrorStrResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.SuccessResp(c)
}

// PublicSettings 获取公共设置
// @Summary 获取公共设置
// @Description 获取公共设置
// @Tags setting
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "公共设置"
// @Router /api/public/settings [get]
func PublicSettings(c *gin.Context) {
	common.SuccessResp(c, op.GetPublicSettingsMap())
}
