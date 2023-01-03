package migration

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"github.com/tcolgate/mp3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Uploads a file to the mongodb database
func uploadFile(path string) primitive.ObjectID {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("No se encontró el archivo:", path)
	}

	id := mongodb.UploadFile(data, fmt.Sprintf("%v.%v", primitive.NewObjectID().Hex(), filepath.Ext(path)))

	return id
}

// Calculates the total duration in seconds of an MP3 file.
func getFileDuration(path string) float64 {
	t := 0.0

	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		t = t + f.Duration().Seconds()
	}

	return t
}
