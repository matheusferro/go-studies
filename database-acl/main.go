package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	pb "matheusferro/database-acl/proto"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	port   = flag.Int("port", 50051, "The server port")
	dbConn *pgxpool.Pool
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedDatabaseServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) GetAllAlbums(ctx context.Context, in *pb.Empty) (*pb.AllAlbums, error) {
	//log.Printf("Received: ")
	//albums := []*pb.Album{
	//    {Id: "1", Title: "The Dark Side Of the Moon", Artist: "Pink Floyd", Price: 99.99},
	//    {Id: "2", Title: "Black Sabbath", Artist: "Black Sabbath", Price: 99.99},
	//    {Id: "3", Title: "Black Sabbath2", Artist: "Black Sabbath2", Price: 99.98},
	//}

	allAlbums, _ := GetAllAlbumsDB(ctx)
	log.Printf("found %d items", len(allAlbums))
	return &pb.AllAlbums{Albums: allAlbums}, nil
}

func OpenConnection() {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PWD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), "library-db", "disable")

	//dataSourceName := "postgres://postgres:postgres@localhost:5432/library-db?sslmode=disable"
	db, err := pgxpool.New(context.Background(), dataSourceName)
	if err != nil {
		fmt.Errorf("failed to open database connection: %w", err)
		panic(err)
	}

	dbConn = db
}

func GetAllAlbumsDB(ctx context.Context) ([]*pb.Album, error) {
	rows, err := dbConn.Query(ctx, "SELECT * FROM album")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var albums []*pb.Album

	for rows.Next() {
		var album pb.Album

		if err := rows.Scan(&album.Id, &album.Title, &album.Artist, &album.Price); err != nil {
			panic(err)
		}

		albums = append(albums, &album)
	}

	return albums, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDatabaseServer(s, &server{})
	OpenConnection()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
