package songs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSongs(p graphql.ResolveParams) (interface{}, error) {
	args := p.Args

	if _, ok := args["songbook_id"]; ok {
		temp, _ := args["songbook_id"].(string)
		songbook_id, err := primitive.ObjectIDFromHex(temp)
		if err != nil {
			panic(err)
		}
		args["songbook_id"] = songbook_id
	}

	if _, ok := args["category_id"]; ok {
		temp, _ := args["category_id"].(string)
		category_id, err := primitive.ObjectIDFromHex(temp)
		if err != nil {
			panic(err)
		}
		args["category_id"] = nil
		args["categories_id"] = category_id
	}

	songs, err := new(models.Song).GetAllSongs(args)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func GetSongById(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["_id"].(string)
	if !ok {
		return nil, nil
	}

	song, err := new(models.Song).GetSongByID(id)
	if err != nil {
		return nil, err
	}

	if song.ID.Hex() == "000000000000000000000000" {
		return nil, nil
	}

	return song, nil
}

func validateNumber(song models.Song, songbook_id string) bool {
	songbook, err := new(models.Songbook).GetSongbookByID(songbook_id, "EN")
	if err != nil {
		panic(err)
	}

	if !songbook.Numeration {
		return true
	}

	db := mongodb.GetMongoDBConnection()
	object_id, err := primitive.ObjectIDFromHex(songbook_id)
	if err != nil {
		panic(err)
	}

	cursor, err := db.Collection("Songs").Aggregate(context.TODO(), []bson.M{
		{"$match": bson.M{
			"songbook_id": object_id,
			"number":      song.Number,
			"_id": bson.M{
				"$ne": song.ID,
			},
		}},
	})

	if err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		return false
	}

	return true
}

func UpdateSong(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["_id"].(string)
	lang := p.Context.Value("language").(string)
	new_song := p.Args["song"].(map[string]interface{})

	song, err := new(models.Song).GetSongByID(id)
	if err != nil {
		return nil, err
	}

	if song.ID.Hex() == "000000000000000000000000" {
		return nil, fmt.Errorf("song not found")
	}

	if err := helpers.BindJSON(new_song, &song); err != nil {
		return nil, err
	}

	if !validateNumber(song, song.SongbookID.Hex()) {
		models.CleanUpSongTemp(new_song)
		return nil, fmt.Errorf(locale.GetLocalizedMessage(lang, "song.error.invalid_number"), strconv.Itoa(song.Number))
	}

	files := []string{"music_sheet_path",
		"music_audio_path",
		"music_audio_only_path",
		"soprano_voice_audio_path",
		"contralto_voice_audio_path",
		"tenor_voice_audio_path",
		"bass_voice_audio_path",
		"all_voice_audio_path",
	}

	for _, file := range files {
		value := ""
		if new_song[file] != nil {
			value = new_song[file].(string)
		}

		var err error = nil
		switch file {
		case "music_sheet_path":
			err = song.UpdateMusicSheet(value)
		case "music_audio_path":
			err = song.UpdateMusicAudio(value)
		case "music_audio_only_path":
			err = song.UpdateMusicAudioOnly(value)
		default:
			err = song.UpdateVoices(value, strings.Split(file, "_")[0])
		}
		if err != nil {
			return nil, err
		}
	}

	if err := song.UpdateSong(); err != nil {
		return nil, err
	}

	return song, nil
}

func CreateSong(p graphql.ResolveParams) (interface{}, error) {
	lang := p.Context.Value("language").(string)
	new_song := p.Args["song"].(map[string]interface{})

	var song models.Song

	if err := helpers.BindJSON(new_song, &song); err != nil {
		return nil, err
	}

	if !validateNumber(song, song.SongbookID.Hex()) {
		models.CleanUpSongTemp(new_song)
		return nil, fmt.Errorf(locale.GetLocalizedMessage(lang, "song.error.invalid_number"), strconv.Itoa(song.Number))
	}

	if err := song.CreateSong(); err != nil {
		return nil, err
	}

	files := []string{"music_sheet_path",
		"music_audio_path",
		"music_audio_only_path",
		"soprano_voice_audio_path",
		"contralto_voice_audio_path",
		"tenor_voice_audio_path",
		"bass_voice_audio_path",
		"all_voice_audio_path",
	}

	for _, file := range files {
		value := ""
		if new_song[file] != nil {
			value = new_song[file].(string)
		}

		var err error = nil
		switch file {
		case "music_sheet_path":
			err = song.UpdateMusicSheet(value)
		case "music_audio_path":
			err = song.UpdateMusicAudio(value)
		case "music_audio_only_path":
			err = song.UpdateMusicAudioOnly(value)
		default:
			err = song.UpdateVoices(value, strings.Split(file, "_")[0])
		}
		if err != nil {
			return nil, err
		}
	}

	return song, nil
}

func DeleteSong(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["_id"].(string)

	song, err := new(models.Song).GetSongByID(id)
	if err != nil {
		return nil, err
	}

	if song.ID.Hex() == "000000000000000000000000" {
		return nil, fmt.Errorf("song not found")
	}

	if err := song.DeleteSong(); err != nil {
		return nil, err
	}

	return song, nil
}
