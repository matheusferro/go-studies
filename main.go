package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Album struct {
    ID      string `json:"id" binding:"required"`
    Title   string `json:"title" binding:"required"`
    Artist  string `json:"artist" binding:"required"`
    Price   float64 `json:"price" binding:"required"`
}


var albums = []Album{
    {ID: "1", Title: "The Dark Side Of the Moon", Artist: "Pink Floyd", Price: 99.99},
    {ID: "2", Title: "Black Sabbath", Artist: "Black Sabbath", Price: 99.99},
}


func getAlbums(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
    id := c.Param("id")

    for _, album := range albums {
        if album.ID == id {
            c.IndentedJSON(http.StatusOK, album)
            return
        }
    }

    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func postAlbums(c *gin.Context) {
    var newAlbum Album
    if err := c.BindJSON(&newAlbum); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid fields received"})
        return
    }
    for _, album := range albums {
        if album.ID == newAlbum.ID {
            c.JSON(http.StatusConflict, gin.H{"message": "Resource already created"})
            return
        }
    }

    albums = append(albums, newAlbum)
    c.IndentedJSON(http.StatusCreated, newAlbum)
}

func setupRouter() *gin.Engine {
    router := gin.Default()
    router.GET("/albums", getAlbums)
    router.GET("/albums/:id", getAlbumByID)
    router.POST("/albums", postAlbums)

    return router
}

func main() {
    router := setupRouter()
    router.Run("localhost:8080")
}
