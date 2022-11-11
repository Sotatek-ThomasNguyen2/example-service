package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"time"

	"github.com/Sotatek-ThomasNguyen2/example-service.git/db"
	pb "github.com/Sotatek-ThomasNguyen2/example-service.git/proto"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type User struct {
	db  IDatabase
	red *redis.Client
}

type IDatabase interface {
	ListUserPartner(rq *pb.UserPartnerRequest) ([]*pb.UserPartner, error)
	FindUserPartner(rq *pb.UserPartnerRequest) (*pb.UserPartner, error)
	InsertUserPartner(up ...*pb.UserPartner) error
	UpdateUserPartner(updator, selector *pb.UserPartner) error
}

func initServe() *User {
	log.SetFlags(log.Lshortfile)
	d := &db.DB{}
	if err := d.ConnectDb("host=localhost user=postgres password=123 dbname=example_service port=5432 sslmode=disable"); err != nil {
		debug.PrintStack()
		log.Panicln(err)
	}
	rd := NewRedisCache("postgres:6379", "")
	log.Println("Connect db success!")
	return &User{
		db:  d,
		red: rd,
	}
}

func (r *Router) HttpRouter(u *User) error {
	ro := r.route
	r.router()
	go ro.Run(":3001")
	return nil
}

func StartGRPCServe(port int, u *User) error {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUserServiceServer(grpcServer, u)
	grpcServer.Serve(lis)
	return nil
}

func NewRedisCache(addr, pw string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw, // no password set
		DB:       0,  // use default DB
	})
	log.Print(addr, pw)
	tick := time.NewTicker(10 * time.Minute)
	ctx := context.Background()
	go func(client *redis.Client) {
		for {
			select {
			case <-tick.C:
				if err := client.Ping(ctx).Err(); err != nil {
					panic(err)
				}
			}
		}
	}(client)
	return client
}
