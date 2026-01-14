// Package comet provides cookie extraction for the Comet browser.
// Comet is a Chromium-based browser with its own keychain entry "Comet Safe Storage".
package comet

import (
	"context"

	"github.com/ASISolutions/kooky"
	"github.com/ASISolutions/kooky/internal/chrome"
	"github.com/ASISolutions/kooky/internal/cookies"
)

func ReadCookies(ctx context.Context, filename string, filters ...kooky.Filter) ([]*kooky.Cookie, error) {
	return cookies.SingleRead(cookieStore, filename, filters...).ReadAllCookies(ctx)
}

func TraverseCookies(filename string, filters ...kooky.Filter) kooky.CookieSeq {
	return cookies.SingleRead(cookieStore, filename, filters...)
}

// CookieStore has to be closed with CookieStore.Close() after use.
func CookieStore(filename string, filters ...kooky.Filter) (kooky.CookieStore, error) {
	return cookieStore(filename, filters...)
}

func cookieStore(filename string, filters ...kooky.Filter) (*cookies.CookieJar, error) {
	s := &chrome.CookieStore{}
	s.FileNameStr = filename
	s.BrowserStr = `comet`
	// Comet uses "Comet Safe Storage" in macOS Keychain
	s.SetSafeStorage(`Comet`, `Comet Safe Storage`)

	return cookies.NewCookieJar(s, filters...), nil
}
