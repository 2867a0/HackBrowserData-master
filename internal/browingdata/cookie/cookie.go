package cookie

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"hack-browser-data/internal/decrypter"
	"hack-browser-data/internal/item"
	"hack-browser-data/internal/log"
	"hack-browser-data/internal/utils/typeutil"

	// import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

type ChromiumCookie []cookie

type cookie struct {
	Host         string
	Path         string
	KeyName      string
	encryptValue []byte
	Value        string
	IsSecure     bool
	IsHTTPOnly   bool
	HasExpire    bool
	IsPersistent bool
	CreateDate   time.Time
	ExpireDate   time.Time
}

const (
	queryChromiumCookie = `SELECT name, encrypted_value, host_key, path, creation_utc, expires_utc, is_secure, is_httponly, has_expires, is_persistent FROM cookies`
)

func (c *ChromiumCookie) Parse(masterKey []byte) error {
	cookieDB, err := sql.Open("sqlite3", item.TempChromiumCookie)
	if err != nil {
		return err
	}
	defer os.Remove(item.TempChromiumCookie)
	defer cookieDB.Close()
	rows, err := cookieDB.Query(queryChromiumCookie)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			key, host, path                               string
			isSecure, isHTTPOnly, hasExpire, isPersistent int
			createDate, expireDate                        int64
			value, encryptValue                           []byte
		)
		if err = rows.Scan(&key, &encryptValue, &host, &path, &createDate, &expireDate, &isSecure, &isHTTPOnly, &hasExpire, &isPersistent); err != nil {
			log.Warn(err)
		}

		cookie := cookie{
			KeyName:      key,
			Host:         host,
			Path:         path,
			encryptValue: encryptValue,
			IsSecure:     typeutil.IntToBool(isSecure),
			IsHTTPOnly:   typeutil.IntToBool(isHTTPOnly),
			HasExpire:    typeutil.IntToBool(hasExpire),
			IsPersistent: typeutil.IntToBool(isPersistent),
			CreateDate:   typeutil.TimeEpoch(createDate),
			ExpireDate:   typeutil.TimeEpoch(expireDate),
		}
		if len(encryptValue) > 0 {
			var err error
			if masterKey == nil {
				value, err = decrypter.DPAPI(encryptValue)
			} else {
				value, err = decrypter.Chromium(masterKey, encryptValue)
			}
			if err != nil {
				log.Error(err)
			}
		}
		cookie.Value = string(value)
		*c = append(*c, cookie)
	}
	sort.Slice(*c, func(i, j int) bool {
		return (*c)[i].CreateDate.After((*c)[j].CreateDate)
	})
	return nil
}

func (c *ChromiumCookie) Name() string {
	return "cookie"
}

func (c *ChromiumCookie) Length() int {
	return len(*c)
}

func (c *ChromiumCookie) SaveCookie(outputPath string, outFormat string) {
	// 存放 处理完之后的cookie
	cookie_data := make(map[string]string)

	// 遍历原始cookie
	for _, data := range *c {
		// 拿到原始cookie里面的host作为处理后的key
		cookie_key := data.Host
		if data.Host[:1] == "." {
			cookie_key = data.Host[1:]
		}

		// 原始cookie里面的key
		key := data.KeyName
		value := data.Value
		expire := strconv.FormatInt(data.ExpireDate.Unix(), 10)
		httpOnly := boolToI(data.IsHTTPOnly)
		secure := boolToI(data.IsSecure)

		map_value := fmt.Sprintf("^%s=%s`%s%s%s", key, value, expire, httpOnly, secure)

		if oldValue, exists := cookie_data[cookie_key]; exists {
			// 存在, 追加
			cookie_data[cookie_key] = oldValue + map_value
		} else {
			// 不存在, 添加
			cookie_data[cookie_key] = map_value
		}

	}

	save(&cookie_data, outputPath, outFormat)
}

type FirefoxCookie []cookie

const (
	queryFirefoxCookie = `SELECT name, value, host, path, creationTime, expiry, isSecure, isHttpOnly FROM moz_cookies`
)

func (f *FirefoxCookie) Parse(masterKey []byte) error {
	cookieDB, err := sql.Open("sqlite3", item.TempFirefoxCookie)
	if err != nil {
		return err
	}
	defer os.Remove(item.TempFirefoxCookie)
	defer cookieDB.Close()
	rows, err := cookieDB.Query(queryFirefoxCookie)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name, value, host, path string
			isSecure, isHTTPOnly    int
			creationTime, expiry    int64
		)
		if err = rows.Scan(&name, &value, &host, &path, &creationTime, &expiry, &isSecure, &isHTTPOnly); err != nil {
			log.Warn(err)
		}
		*f = append(*f, cookie{
			KeyName:    name,
			Host:       host,
			Path:       path,
			IsSecure:   typeutil.IntToBool(isSecure),
			IsHTTPOnly: typeutil.IntToBool(isHTTPOnly),
			CreateDate: typeutil.TimeStamp(creationTime / 1000000),
			ExpireDate: typeutil.TimeStamp(expiry),
			Value:      value,
		})
	}
	return nil
}

func (f *FirefoxCookie) Name() string {
	return "cookie"
}

func (f *FirefoxCookie) Length() int {
	return len(*f)
}

func (f *FirefoxCookie) SaveCookie(outputPath string, outFormat string) {
	// 存放 处理完之后的cookie
	cookie_data := make(map[string]string)

	// 遍历原始cookie
	for _, data := range *f {
		// 拿到原始cookie里面的host作为处理后的key
		cookie_key := data.Host
		if data.Host[:1] == "." {
			cookie_key = data.Host[1:]
		}

		// 原始cookie里面的key
		key := data.KeyName
		value := data.Value
		expire := strconv.FormatInt(data.ExpireDate.Unix(), 10)
		httpOnly := boolToI(data.IsHTTPOnly)
		secure := boolToI(data.IsSecure)

		map_value := fmt.Sprintf("^%s=%s`%s%s%s", key, value, expire, httpOnly, secure)

		if oldValue, exists := cookie_data[cookie_key]; exists {
			// 存在, 追加
			cookie_data[cookie_key] = oldValue + map_value
		} else {
			// 不存在, 添加
			cookie_data[cookie_key] = map_value
		}
	}

	save(&cookie_data, outputPath, outFormat)
}

func boolToI(b bool) string {
	if b == true {
		return "1"
	} else {
		return "0"
	}
}

func save(cookie_data *map[string]string, outputPath string, outputFormat string) {
	// 保存
	fName := outputPath + "\\" + strings.ReplaceAll(outputFormat, ".", "_") + "_cookie.txt"
	file, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Errorf("save file %s error: ", fName, err.Error())
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("close file %s error: ", fName, err.Error())
		}
	}(file)
	write := bufio.NewWriter(file)
	for k, v := range *cookie_data {
		_, _ = write.WriteString(k + "/" + v + "\n")
	}
	_ = write.Flush()
	log.Noticef("output to file %s success", fName)
}
