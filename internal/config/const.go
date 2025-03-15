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
	ImageTypes   = "image_types"

	// Site
	VERSION          = "version"
	SiteTitle        = "site_title"
	Announcement     = "announcement"
	RobotsTxt        = "robots_txt"
	IndexBackground  = "index_background"
	IndexTitle       = "index_title"
	IndexDescription = "index_description"
	PageSize         = "page_size"

	Logo    = "logo"
	Favicon = "favicon"

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
