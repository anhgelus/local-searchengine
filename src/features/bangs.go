package features

import (
	"fmt"
	"net/url"
	"strings"
)

var filterBangs = map[string]string{
	"!r":  "site:reddit.com",
	"!gr": "site:grafikart.fr",
	"!i":  "site:stackoverflow.com OR site:github.com",
}

var RedirectBangs = map[string]string{
	"!rt":   "https://www.rottentomatoes.com/search?search=%s",
	"!npm":  "https://www.npmjs.com/search?q=%s",
	"!tr":   "https://www.wordreference.com/enfr/%s",
	"!imdb": "https://www.imdb.com/find?s=all&q=%s",
	"!gh":   "https://github.com/search?q=%s",
	"!yt":   "https://www.youtube.com/results?search_query=%s",
	"!d":    "https://www.larousse.fr/dictionnaires/francais/%s",
	"!g":    "https://www.google.com/search?q=%s",
	"!w":    "/weather?q=%s",
	"!gmap": "https://www.google.fr/maps?hl=fr&q=%s",
	"!img":  "https://duckduckgo.com/?q=%s&iax=images&ia=images",
}

func ParseFilterBangs(q string) string {
	for bang, replace := range filterBangs {
		if strings.HasPrefix(q, bang+" ") ||
			strings.HasSuffix(q, " "+bang) {
			q = strings.ReplaceAll(q, bang, replace)
		}
	}
	return q
}

func ParseRedirectBangs(q string) string {
	if q == "" {
		return ""
	}
	for bang, redirect := range RedirectBangs {
		if strings.HasPrefix(q, bang+" ") ||
			strings.HasSuffix(q, " "+bang) {
			return fmt.Sprintf(
				redirect,
				url.QueryEscape(
					strings.TrimSpace(
						strings.ReplaceAll(q, bang, ""),
					),
				),
			)
		}
	}
	return ""
}
