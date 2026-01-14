package comet

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/browserutils/kooky"
	"github.com/browserutils/kooky/internal/chrome"
	"github.com/browserutils/kooky/internal/cookies"
)

type cometFinder struct{}

var _ kooky.CookieStoreFinder = (*cometFinder)(nil)

func init() {
	kooky.RegisterFinder(`comet`, &cometFinder{})
}

func (f *cometFinder) FindCookieStores() kooky.CookieStoreSeq {
	return func(yield func(kooky.CookieStore, error) bool) {
		// Comet only runs on macOS
		if runtime.GOOS != "darwin" {
			return
		}

		cfgDir, err := os.UserConfigDir()
		if err != nil {
			yield(nil, err)
			return
		}

		// Comet stores cookies in ~/Library/Application Support/Comet/Default/Cookies
		root := filepath.Join(cfgDir, "Comet")

		// Check for Default profile
		cookiesPaths := []string{
			filepath.Join(root, "Default", "Network", "Cookies"), // Chrome 96+
			filepath.Join(root, "Default", "Cookies"),            // Older versions
		}

		for _, cookiesPath := range cookiesPaths {
			if _, err := os.Stat(cookiesPath); err == nil {
				s := &chrome.CookieStore{}
				s.FileNameStr = cookiesPath
				s.BrowserStr = `comet`
				s.ProfileStr = `Default`
				s.IsDefaultProfileBool = true
				s.OSStr = runtime.GOOS
				// Set Comet's keychain credentials
				s.SetSafeStorage(`Comet`, `Comet Safe Storage`)

				st := &cookies.CookieJar{CookieStore: s}
				if !yield(st, nil) {
					return
				}
			}
		}
	}
}
