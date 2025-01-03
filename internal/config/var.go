package config

import (
	"net/url"
	"regexp"
)

var (
	Version    string = "dev"
	WebVersion string
)

var (
	Conf *Config
	URL  *url.URL
)

var (
	RawIndexHtml string
	ManageHtml   string
	IndexHtml    string
)

var FilenameCharMap = make(map[string]string)
var PrivacyReg []*regexp.Regexp
