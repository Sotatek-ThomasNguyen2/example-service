package db

import (
	"log"
	"testing"

	pb "github.com/Sotatek-ThomasNguyen2/example-service.git/proto"
)

func Test_connection(t *testing.T) {
	d := &DB{}
	err := d.ConnectDb("host=localhost user=postgres password=123 dbname=example_service port=5432 sslmode=disable")
	if err != nil {
		log.Print(err)
	}
}

func Test_listPartners(t *testing.T) {
	d := &DB{}
	err := d.ConnectDb("host=localhost user=postgres password=123 dbname=example_service port=5432 sslmode=disable")
	if err != nil {
		log.Print(err)
	}
	list, err := d.ListUserPartner(&pb.UserPartnerRequest{Limit: 5})
	log.Print("list: ", list)
}
