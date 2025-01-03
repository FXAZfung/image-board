package op

import (
	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

// Setting
type SettingItemHook func(item *model.SettingItem) error

var settingItemHooks = map[string]SettingItemHook{
	//conf.VideoTypes: func(item *model.SettingItem) error {
	//	conf.SlicesMap[conf.VideoTypes] = strings.Split(item.Value, ",")
	//	return nil
	//},
	//conf.AudioTypes: func(item *model.SettingItem) error {
	//	conf.SlicesMap[conf.AudioTypes] = strings.Split(item.Value, ",")
	//	return nil
	//},
	//conf.ImageTypes: func(item *model.SettingItem) error {
	//	conf.SlicesMap[conf.ImageTypes] = strings.Split(item.Value, ",")
	//	return nil
	//},
	//conf.TextTypes: func(item *model.SettingItem) error {
	//	conf.SlicesMap[conf.TextTypes] = strings.Split(item.Value, ",")
	//	return nil
	//},
	//conf.ProxyTypes: func(item *model.SettingItem) error {
	//	conf.SlicesMap[conf.ProxyTypes] = strings.Split(item.Value, ",")
	//	return nil
	//},
	//conf.ProxyIgnoreHeaders: func(item *model.SettingItem) error {
	//	conf.SlicesMap[conf.ProxyIgnoreHeaders] = strings.Split(item.Value, ",")
	//	return nil
	//},
	conf.PrivacyRegs: func(item *model.SettingItem) error {
		regStrs := strings.Split(item.Value, "\n")
		regs := make([]*regexp.Regexp, 0, len(regStrs))
		for _, regStr := range regStrs {
			reg, err := regexp.Compile(regStr)
			if err != nil {
				return errors.WithStack(err)
			}
			regs = append(regs, reg)
		}
		conf.PrivacyReg = regs
		return nil
	},
	conf.FilenameCharMapping: func(item *model.SettingItem) error {
		err := utils.Json.UnmarshalFromString(item.Value, &conf.FilenameCharMap)
		if err != nil {
			return err
		}
		log.Debugf("filename char mapping: %+v", conf.FilenameCharMap)
		return nil
	},
	//conf.IgnoreDirectLinkParams: func(item *model.SettingItem) error {
	//	conf.SlicesMap[conf.IgnoreDirectLinkParams] = strings.Split(item.Value, ",")
	//	return nil
	//},
}

func RegisterSettingItemHook(key string, hook SettingItemHook) {
	settingItemHooks[key] = hook
}

func HandleSettingItemHook(item *model.SettingItem) (hasHook bool, err error) {
	if hook, ok := settingItemHooks[item.Key]; ok {
		return true, hook(item)
	}
	return false, nil
}

// Category
type CategoryHook func(category *model.Category) error

var categoryHooks = map[string]CategoryHook{}

func RegisterCategoryHook(key string, hook CategoryHook) {
	categoryHooks[key] = hook
}

func HandleCategoryHook(category *model.Category) (hasHook bool, err error) {
	if hook, ok := categoryHooks[category.Name]; ok {
		return true, hook(category)
	}
	return false, nil
}
