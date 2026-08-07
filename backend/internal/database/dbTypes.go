package database

type Language string

const (
	LanguageEN Language = "en"
	LanguageES Language = "es"
)

func (lang Language) IsValid() bool {
	if lang == LanguageEN || lang == LanguageES {
		return true
	}
	return false
}

////

type Role string

const (
	RoleVisitor Role = "visitor"
	RoleAgency Role = "agency"
	RoleCreator Role = "creator"
	RoleBrand Role = "brand"
)

func (role Role) IsValid() bool {
	if role == RoleVisitor || 
		role == RoleAgency || 
		role == RoleCreator || 
		role == RoleBrand {
			return true
		}
	return false
}

///

type Service string

const (
	ServiceYouTube Service = "youtube"
	ServiceInstagram Service = "instagram"
	ServiceTiktok Service = "tiktok"
	ServiceTwitter Service = "twitter"
	ServiceTwitch Service = "twitch"
)

func (service Service) IsValid() bool {
	if service == ServiceYouTube ||
	service == ServiceInstagram ||
	service == ServiceTiktok ||
	service == ServiceTwitter ||
	service == ServiceTwitch {
		return true
	}
	return false
}