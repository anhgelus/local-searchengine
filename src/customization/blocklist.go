package customization

import "strings"

var blocklist = []string{
	"pinterest.com",
	"allocine.com",
	"jeuxvideo.com",
	"lemonde.fr",
	"w3schools.com",
	"pinterest.fr",
}

var DefaultBlocklist = []string{
	"pinterest.com",
	"allocine.com",
	"jeuxvideo.com",
	"lemonde.fr",
	"w3schools.com",
	"pinterest.fr",
}

func IsBlockedSite(s string) bool {
	for _, domain := range blocklist {
		if strings.HasSuffix(s, domain) {
			return true
		}
	}
	return false
}

func UpdateBlockList(l []string) {
	for _, site := range l {
		blocklist = append(blocklist, site)
	}
}
