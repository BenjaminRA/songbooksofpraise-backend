package songs

import (
	"net/http"

	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
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

// func GetSongs(c *gin.Context) {
// 	songs := new(models.Song).GetAllSongs()

// 	c.IndentedJSON(http.StatusOK, songs)
// }

// func GetSongsById(c *gin.Context) {
// 	id := c.Param("id")

// 	song := new(models.Song).GetSongByID(id)

// 	if song.ID.Hex() == "000000000000000000000000" {
// 		c.IndentedJSON(http.StatusNotFound, song)
// 	} else {
// 		c.IndentedJSON(http.StatusOK, song)
// 	}

// }

func GetMusicSheet(c *gin.Context) {
	id := c.Param("id")
	data, filename := new(models.Song).GetMusicSheet(id)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func GetVoicesByVoice(c *gin.Context) {
	id := c.Param("id")
	voice := c.Param("voice")
	data, filename, err := new(models.Song).GetVoice(id, voice)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}
