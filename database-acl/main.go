package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
    "strconv"

	"google.golang.org/grpc"
	pb "matheusferro/database-acl/proto"

	spannergorm "github.com/googleapis/go-gorm-spanner"
	"gorm.io/gorm"
)

var (
	port   = flag.Int("port", 50051, "The server port")
	dbConn *gorm.DB
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedDatabaseServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) GetAllAlbums(ctx context.Context, in *pb.Empty) (*pb.AllAlbums, error) {
	allAlbums, _ := GetAllAlbumsDB(ctx)
	log.Printf("found %d items", len(allAlbums))
	return &pb.AllAlbums{Albums: allAlbums}, nil
}

func OpenConnection() {
	db, err := gorm.Open(spannergorm.New(spannergorm.Config{
		DriverName: "spanner",
		DSN:        "projects/test-project/instances/test-instance/databases/test-database",
	}), &gorm.Config{PrepareStmt: true})
	if err != nil {
		fmt.Errorf("failed to open database connection: %w", err)
		panic(err)
	}
    
	dbConn = db
}

type Album struct {
	gorm.Model
	Id     int64
	Title  string
	Artist string
	Price  float64
}

func (Album) TableName() string {
    return "album"
}

func GetAllAlbumsDB(ctx context.Context) ([]*pb.Album, error) {
	var albums []*Album
	dbConn.Unscoped().Find(&albums)
    albumsProto := make([]*pb.Album, len(albums))

    for i, a := range albums {
        albumsProto[i] = &pb.Album {Id: strconv.FormatInt(a.Id, 10), Title: a.Title, Artist: a.Artist, Price: a.Price}
    }

	return albumsProto, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDatabaseServer(s, &server{})
	os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")
	OpenConnection()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
