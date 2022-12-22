//go:build linux

package chromium

import (
	"crypto/sha1"
	"errors"
	"os"

	"github.com/godbus/dbus/v5"
	keyring "github.com/ppacher/go-dbus-keyring"
	"golang.org/x/crypto/pbkdf2"

	"hack-browser-data/internal/item"
	"hack-browser-data/internal/log"
)

func (c *chromium) GetMasterKey() ([]byte, error) {
	// what is d-bus @https://dbus.freedesktop.org/
	var chromiumSecret []byte
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	defer os.Remove(item.TempChromiumKey)
	svc, err := keyring.GetSecretService(conn)
	if err != nil {
		return nil, err
	}
	session, err := svc.OpenSession()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := session.Close(); err != nil {
			log.Errorf("close session failed: %v", err)
		}
	}()
	collections, err := svc.GetAllCollections()
	if err != nil {
		return nil, err
	}
	for _, col := range collections {
		items, err := col.GetAllItems()
		if err != nil {
			return nil, err
		}
		for _, i := range items {
			label, err := i.GetLabel()
			if err != nil {
				log.Error(err)
				continue
			}
			if label == c.storage {
				se, err := i.GetSecret(session.Path())
				if err != nil {
					return nil, errors.New("get storage from dbus error:" + err.Error())
				}
				chromiumSecret = se.Value
			}
		}
	}
	if chromiumSecret == nil {
		// @https://source.chromium.org/chromium/chromium/src/+/main:components/os_crypt/os_crypt_linux.cc;l=100
		chromiumSecret = []byte("peanuts")
	}
	chromiumSalt := []byte("saltysalt")
	// @https://source.chromium.org/chromium/chromium/src/+/master:components/os_crypt/os_crypt_linux.cc
	key := pbkdf2.Key(chromiumSecret, chromiumSalt, 1, 16, sha1.New)
	c.masterKey = key
	log.Infof("%s initialized master key success", c.name)
	return key, nil
}
