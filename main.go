package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	pb "brawl_grpc/proto"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

var db *sql.DB

type server struct {
	pb.UnimplementedBrawlerServiceServer
}

func (s server) GetBrawlersInfo(ctx context.Context, req *pb.BrawlerRequest) (*pb.BrawlerResponse, error) {
	var name, ptype, category string

	query := "SELECT * FROM brawl.brawler where Name like @Name"
	row := db.QueryRowContext(ctx, query, sql.Named("Name", "%"+req.Name+"%"))
	err := row.Scan(&name, &ptype, &category)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.BrawlerResponse{
				Name:     "Not found",
				Type:     "Not found",
				Category: "Not found",
			}, nil
		}
		return nil, err
	}
	return &pb.BrawlerResponse{
		Name:     name,
		Type:     ptype,
		Category: category,
	}, nil
}

func (s *server) GetBrawlerList(req *pb.Empty, stream pb.BrawlerService_GetBrawlerListServer) error {
	query := "SELECT * FROM brawl.brawler"
	rows, err := db.Query(query)
	if err != nil {
		log.Panic(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name, ptype, category string

		if err := rows.Scan(&name, &ptype, &category); err != nil {
			log.Panic(err)
			return err
		}
		if err := stream.Send(&pb.BrawlerResponse{
			Name:     name,
			Type:     ptype,
			Category: category,
		}); err != nil {
			log.Panic(err)
			return err
		}
	}
	return nil

}

func (s *server) AddBrawler(stream pb.BrawlerService_AddBrawlerServer) error {
	var count int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.AddBrawlerResponse{
				Count: count,
			})
		}
		if err != nil {
			log.Panic(err)
			return err
		}
		query := "insert into brawl.brawler (Name, Type, Category) values (@Name, @Type, @Category)"
		_, err = db.Exec(query,
			sql.Named("Name", req.Name),
			sql.Named("Type", req.Type),
			sql.Named("Category", req.Category))

		if err != nil {
			log.Panic(err)
			return err
		}
		count++
		log.Printf("Added %s", req.Name)
	}
}

func (s *server) GetBrawlerByTyppe(stream pb.BrawlerService_GetBrawlerByTyppeServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("End of stream")
			return nil
		}
		if err != nil {
			log.Panic(err)
			return err
		}

		query := "select * from brawl.brawler where lower(Type) = lower(@Type)"
		rows, err := db.Query(query, sql.Named("Type", req.Type))
		if err != nil {
			log.Panic(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var name, ptype, category string

			if err := rows.Scan(&name, &ptype, &category); err != nil {
				log.Panic(err)
				return err
			}

			if err := stream.Send(&pb.BrawlerResponse{
				Name:     name,
				Type:     ptype,
				Category: category,
			}); err != nil {
				log.Panic(err)
				return err
			}

		}
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	s := os.Getenv("DB_SERVER")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		user, pass, s, port, database)

	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to database")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBrawlerServiceServer(grpcServer, &server{})

	//iniciar el servidor HTTP para health check
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		log.Println("Starting health check server on port :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	log.Println("Starting server on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
