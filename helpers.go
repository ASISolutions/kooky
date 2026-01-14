package kooky

import (
	"strings"
)

// ToCookieHeader converts a slice of cookies to an HTTP Cookie header value.
// The format is "name1=value1; name2=value2; ..."
func ToCookieHeader(cookies []*Cookie) string {
	if len(cookies) == 0 {
		return ""
	}

	parts := make([]string, 0, len(cookies))
	for _, c := range cookies {
		if c == nil {
			continue
		}
		parts = append(parts, c.Name+"="+c.Value)
	}
	return strings.Join(parts, "; ")
}

// ToCookieHeaderFromSeq converts a CookieSeq to an HTTP Cookie header value.
// Errors in the sequence are ignored.
func ToCookieHeaderFromSeq(seq CookieSeq) string {
	if seq == nil {
		return ""
	}

	var parts []string
	for cookie, err := range seq {
		if err != nil || cookie == nil {
			continue
		}
		parts = append(parts, cookie.Name+"="+cookie.Value)
	}
	return strings.Join(parts, "; ")
}

// Cookies helper method
func (c Cookies) ToCookieHeader() string {
	return ToCookieHeader(c)
}
