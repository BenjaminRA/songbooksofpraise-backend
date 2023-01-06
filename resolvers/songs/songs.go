package songs

import (
	"fmt"
	"strings"

	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
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

	songs := new(models.Song).GetAllSongs(args)

	return songs, nil
}

func GetSongById(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["_id"].(string)
	if !ok {
		return nil, nil
	}

	song := new(models.Song).GetSongByID(id)
	if song.ID.Hex() == "000000000000000000000000" {
		return nil, nil
	}

	return song, nil
}

func UpdateSong(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["_id"].(string)
	new_song := p.Args["song"].(map[string]interface{})

	song := new(models.Song).GetSongByID(id)

	files := []string{"music_sheet_path",
		"music_audio_path",
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
		default:
			err = song.UpdateVoices(value, strings.Split(file, "_")[0])
		}
		if err != nil {
			return nil, err
		}
	}

	if song.ID.Hex() == "000000000000000000000000" {
		return nil, fmt.Errorf("song not found")
	}

	if err := helpers.BindJSON(new_song, &song); err != nil {
		return nil, err
	}

	if err := song.UpdateSong(); err != nil {
		return nil, err
	}

	return song, nil
}
