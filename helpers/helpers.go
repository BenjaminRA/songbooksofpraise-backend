package helpers

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func Map(t interface{}, function func(interface{}) interface{}) interface{} {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		arr := []interface{}{}
		s := reflect.Indirect(reflect.ValueOf(t))
		for i := 0; i < s.Len(); i++ {
			arr = append(arr, function(s.Index(i).Interface()))
		}
		return arr
	}
	return nil
}

// func SongsMap(t *[]models.Himno, function func(himno models.Himno) models.Song) []models.Song {
// 	arr := []models.Song{}

// 	for _, value := range *t {
// 		arr = append(arr, function(value))
// 	}
// 	return arr
// }

// func VerseMap(t *[]models.Parrafo, function func(himno models.Parrafo) models.Verse) []models.Verse {
// 	arr := []models.Verse{}

// 	for _, value := range *t {
// 		arr = append(arr, function(value))
// 	}
// 	return arr
// }

func BindJSON(jsonObject interface{}, object interface{}) error {
	jsonBody, err := json.Marshal(jsonObject)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBody, &object)
	if err != nil {
		return err
	}

	return nil
}

func LoadLocalEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func HashValue(value string) string {
	hash := sha256.New()
	key := hash.Sum([]byte(value))

	return fmt.Sprintf("%x", key)
}

func GetUniqueMD5ID() string {
	LoadLocalEnv()
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s____%s", os.Getenv("SECRET"), time.Now())))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetFilenameFromPath(path string) string {
	var filename string
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			break
		}
		filename = string(path[i]) + filename
	}

	// Remove MD5ID
	var finalFilename string
	isMD5ID := false

	for i := 0; i < len(filename); i++ {
		if filename[i] == '-' {
			isMD5ID = true
			continue
		}

		if filename[i] == '.' {
			isMD5ID = false
		}

		if !isMD5ID {
			finalFilename += string(filename[i])
		}
	}

	return finalFilename
}

// Helper functions for SQL escaping
func SqlEscape(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func SqlEscapeNullString(s *string) string {
	if s == nil {
		return "NULL"
	}
	return SqlEscape(*s)
}

func SqlEscapeNullInt(i *int) string {
	if i == nil {
		return "NULL"
	}
	return fmt.Sprintf("%d", *i)
}
