package handles

import (
	"net/http"
	"strconv"
	"strings"

	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/internal/sign"
	"github.com/FXAZfung/image-board/pkg/random"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/FXAZfung/image-board/server/static"
	"github.com/gin-gonic/gin"
)

// ResetToken 重置token
// @Summary 重置token
// @Description 重置token
// @Tags 设置
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Success 200 {object} string
// @Router /api/setting/token [post]
func ResetToken(c *gin.Context) {
	token := random.SecretKey()
	item := model.SettingItem{Key: "token", Value: token, Type: conf.TypeString, Group: model.SINGLE, Flag: model.PRIVATE}
	if err := op.SaveSettingItem(&item); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}
	sign.Instance()
	common.SuccessResp(c, token)
}

// GetSetting 获取设置
// @Summary 获取设置
// @Description 获取设置
// @Tags 设置
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param key query string false "key"
// @Param keys query string false "keys"
// @Success 200 {object} string "设置"
// @Router /api/setting/setting [get]
func GetSetting(c *gin.Context) {
	key := c.Query("key")
	keys := c.Query("keys")
	if key != "" {
		item, err := op.GetSettingItemByKey(key)
		if err != nil {
			common.ErrorResp(c, http.StatusBadRequest, err)
			return
		}
		common.SuccessResp(c, item)
	} else {
		items, err := op.GetSettingItemInKeys(strings.Split(keys, ","))
		if err != nil {
			common.ErrorResp(c, http.StatusBadRequest, err)
			return
		}
		common.SuccessResp(c, items)
	}
}

// SaveSettings 保存设置
// @Summary 保存系统设置
// @Description 保存一组系统设置项
// @Tags 设置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param settings body []model.SettingItem true "设置项列表"
// @Success 200 {object} common.Resp "保存成功"
// @Failure 400 {object} common.Resp "请求格式错误"
// @Failure 401 {object} common.Resp "未授权，需要登录"
// @Failure 500 {object} common.Resp "保存失败"
// @Router /api/setting/save [post]
func SaveSettings(c *gin.Context) {
	var req []model.SettingItem
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}
	if err := op.SaveSettingItems(req); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
	} else {
		common.SuccessResp(c)
		static.UpdateIndex()
	}
}

// ListSettings 列出设置
// @Summary 获取系统设置列表
// @Description 按分组列出系统设置项（需要认证）
// @Tags 设置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param group query string false "设置分组ID"
// @Param groups query string false "多个设置分组ID，逗号分隔"
// @Success 200 {object} common.Resp{data=[]model.SettingItem} "设置项列表"
// @Failure 400 {object} common.Resp "参数格式错误"
// @Failure 401 {object} common.Resp "未授权，需要登录"
// @Failure 500 {object} common.Resp "查询失败"
// @Router /api/setting/list [get]
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
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}
	common.SuccessResp(c, settings)
}

// DeleteSetting 删除设置
// @Summary 删除设置
// @Description 删除设置
// @Tags 设置
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Param key query string true "key"
// @Router /api/setting/delete [delete]
func DeleteSetting(c *gin.Context) {
	key := c.Query("key")
	if err := op.DeleteSettingItemByKey(key); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}
	common.SuccessResp(c)
}

// PublicSettings 获取公共设置
// @Summary 获取公共设置
// @Description 获取所有公开的系统设置
// @Tags 设置
// @Accept json
// @Produce json
// @Success 200 {object} common.Resp{data=map[string]string} "公共设置键值对"
// @Failure 500 {object} common.Resp "获取设置失败"
// @Router /api/setting [get]
func PublicSettings(c *gin.Context) {
	common.SuccessResp(c, op.GetPublicSettingsMap())
}
