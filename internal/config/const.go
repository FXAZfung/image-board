package config

const (
	TypeString = "string"
	TypeSelect = "select"
	TypeBool   = "bool"
	TypeText   = "text"
	TypeNumber = "number"
)

const (

	// image
	ImageMaxSize = "image_max_size"

	// Site
	VERSION      = "version"
	SiteTitle    = "site_title"
	Announcement = "announcement"
	RobotsTxt    = "robots_txt"

	Logo    = "logo"
	Favicon = "favicon"

	// Preview
	ImageTypes = "image_types"

	// Gloabl
	PrivacyRegs         = "privacy_regs"
	FilenameCharMapping = "filename_char_mapping"
	LinkExpiration      = "link_expiration"

	// single
	Token = "token"
)

const (
	UNKNOWN = iota
	FOLDER
	// OFFICE
	VIDEO
	AUDIO
	TEXT
	IMAGE
)
