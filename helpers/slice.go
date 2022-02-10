package slice

import (
	"reflect"

	"github.com/BenjaminRA/himnario-backend/models"
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

func SongsMap(t *[]models.Himno, function func(himno models.Himno) models.Song) []models.Song {
	arr := []models.Song{}

	for _, value := range *t {
		arr = append(arr, function(value))
	}
	return arr
}

func VerseMap(t *[]models.Parrafo, function func(himno models.Parrafo) models.Verse) []models.Verse {
	arr := []models.Verse{}

	for _, value := range *t {
		arr = append(arr, function(value))
	}
	return arr
}
