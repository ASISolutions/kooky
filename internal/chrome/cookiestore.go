package chrome

import (
	"errors"
	"os"
	"time"

	"github.com/go-sqlite/sqlite3"

	"github.com/browserutils/kooky/internal/cookies"
	"github.com/browserutils/kooky/internal/utils"
)

type CookieStore struct {
	cookies.DefaultCookieStore
	Database             *sqlite3.DbFile
	KeyringPasswordBytes []byte
	PasswordBytes        []byte
	DecryptionMethod     func(data, password []byte, dbVersion int64) ([]byte, error)
	storage              safeStorage
	dbVersion            int64
	tempDBPath           string // Path to temp DB copy (if used)
}

func (s *CookieStore) Open() error {
	if s == nil {
		return errors.New(`cookie store is nil`)
	}
	if s.Database != nil {
		return nil
	}

	// Try to open directly first with retries
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt*100) * time.Millisecond)
		}

		f, err := utils.OpenFile(s.FileNameStr)
		if err != nil {
			lastErr = err
			if utils.IsDBLocked(err) {
				continue
			}
			return err
		}

		db, err := sqlite3.OpenFrom(f)
		if err != nil {
			f.Close()
			lastErr = err
			if utils.IsDBLocked(err) {
				continue
			}
			return err
		}

		s.Database = db
		return nil
	}

	// If retries failed due to locked DB, try copying to temp
	if utils.IsDBLocked(lastErr) {
		return s.openFromCopy()
	}

	return lastErr
}

// openFromCopy copies the database to a temp file and opens that instead.
// This works around locked database issues when the browser is running.
func (s *CookieStore) openFromCopy() error {
	tmpPath, err := utils.CopyDBToTemp(s.FileNameStr)
	if err != nil {
		return err
	}
	s.tempDBPath = tmpPath

	f, err := utils.OpenFile(tmpPath)
	if err != nil {
		os.Remove(tmpPath)
		s.tempDBPath = ""
		return err
	}

	db, err := sqlite3.OpenFrom(f)
	if err != nil {
		f.Close()
		os.Remove(tmpPath)
		s.tempDBPath = ""
		return err
	}

	s.Database = db
	return nil
}

func (s *CookieStore) Close() error {
	if s == nil {
		return errors.New(`cookie store is nil`)
	}
	if s.Database == nil {
		return nil
	}
	err := s.Database.Close()
	if err == nil {
		s.Database = nil
	}

	// Clean up temp DB copy if one was used
	if s.tempDBPath != "" {
		os.Remove(s.tempDBPath)
		s.tempDBPath = ""
	}

	return err
}

var _ cookies.CookieStore = (*CookieStore)(nil)

// returns the previous password for later restoration
// used in tests
func (s *CookieStore) SetKeyringPassword(password []byte) []byte {
	if s == nil {
		return nil
	}
	oldPassword := s.KeyringPasswordBytes
	s.KeyringPasswordBytes = password
	return oldPassword
}
