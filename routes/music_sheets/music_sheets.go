package music_sheets_route

// getAlbums responds with the list of all albums as JSON.

// func GetAlbumsById(c *gin.Context) {
// 	id := c.Param("id")

// 	for _, value := range albums {
// 		if value.ID == id {
// 			c.IndentedJSON(http.StatusOK, value)
// 			return
// 		}
// 	}

// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})

// }

// // postAlbums adds an album from JSON received in the request body.
// func PostAlbums(c *gin.Context) {
// 	var newAlbum album

// 	// Call BindJSON to bind the received JSON to
// 	// newAlbum.
// 	if err := c.BindJSON(&newAlbum); err != nil {
// 		return
// 	}

// 	// Add the new album to the slice.
// 	albums = append(albums, newAlbum)
// 	c.IndentedJSON(http.StatusCreated, newAlbum)
// }
