package migration

import (
	"io"
	"os"

	"github.com/tcolgate/mp3"
)

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
