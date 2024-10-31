package main

import (
	"context"
	"fmt"
	"log"
	pb "matheusferro/go-webservice/proto"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Album struct {
	ID     string  `json:"id" binding:"required"`
	Title  string  `json:"title" binding:"required"`
	Artist string  `json:"artist" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}

var (
	dbClient pb.DatabaseClient
	mutex    sync.Mutex
)

var albums = []Album{
	{ID: "1", Title: "The Dark Side Of the Moon", Artist: "Pink Floyd", Price: 99.99},
	{ID: "2", Title: "Black Sabbath", Artist: "Black Sabbath", Price: 99.99},
}

func databaseConnection() pb.DatabaseClient {
	mutex.Lock()
	defer mutex.Unlock()
	if dbClient == nil {
		addr := fmt.Sprintf("%s:%s", os.Getenv("DB_ACL_HOST"), os.Getenv("DB_ACL_PORT"))
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
			return nil
		}
		dbClient = pb.NewDatabaseClient(conn)
	}
	return dbClient
}

func getAlbums(c *gin.Context) {
	databaseConn := databaseConnection()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := databaseConn.GetAllAlbums(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get albums: %v", err)
	}
	allAlbums := make([]*Album, 0, len(r.Albums))
	for _, ab := range r.Albums {
		//log.Printf("Album title: %s", ab.Title)
		allAlbums = append(allAlbums, &Album{ID: ab.Id, Title: ab.Title, Artist: ab.Artist, Price: ab.Price})
	}
	c.IndentedJSON(http.StatusOK, allAlbums)

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
