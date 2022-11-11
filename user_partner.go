package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/Sotatek-ThomasNguyen2/example-service.git/common"
	pb "github.com/Sotatek-ThomasNguyen2/example-service.git/proto"
	"github.com/Sotatek-ThomasNguyen2/example-service.git/utils"
	"google.golang.org/grpc/codes"
)

const DEFAULT_LIMIT = 30

func (u *User) ListUserPartners(ctx context.Context, req *pb.UserPartnerRequest) (*pb.UserPartners, error) {
	if req.GetLimit() == 0 {
		req.Limit = DEFAULT_LIMIT
	}
	up, err := u.db.ListUserPartner(req)
	if err != nil {
		return nil, err
	}
	uPs := &pb.UserPartners{
		UserPartners: up,
	}
	return uPs, nil
}

func (u *User) CreateUserPartner(ctx context.Context, req *pb.UserPartner) (*pb.UserPartner, error) {
	req.Created = time.Now().Unix()
	if err := u.db.InsertUserPartner(req); err != nil {
		log.Print(err)
		return nil, common.Err(codes.Internal, utils.E_internal_error)
	}
	up, err := u.db.FindUserPartner(&pb.UserPartnerRequest{UserId: req.GetUserId(), Phone: req.GetPartnerId()})
	if err != nil {
		log.Print(err)
		return nil, common.Err(codes.Internal, utils.E_internal_error)
	}
	return up, nil
}

func (u *User) UpdateUserPartner(ctx context.Context, req *pb.UserPartner) (*pb.Empty, error) {
	req.UpdatedAt = time.Now().Unix()
	err := u.db.UpdateUserPartner(&pb.UserPartner{UserId: req.GetUserId(), AliasUserId: req.GetAliasUserId(), PartnerId: req.GetPartnerId()}, req)
	if err != nil {
		log.Print(err)
		return nil, common.Err(codes.Internal, err.Error())
	}
	return &pb.Empty{}, nil
}

func (u *User) GetUserPartner(ctx context.Context, req *pb.UserPartnerRequest) (*pb.UserPartner, error) {
	up, err := u.db.FindUserPartner(req)
	if err != nil {
		log.Print(err)
		return nil, common.Err(codes.Internal, err.Error())
	}
	return up, nil
}

func (u *User) ScanToCache() error {
	data := make(chan string, 100)
	wg := &sync.WaitGroup{}
	cCtx, cancel := context.WithCancel(context.Background())
	file, err := os.Open("name.txt")
	if err != nil {
		log.Panicln(err)
	}
	go func() {
		for {
			select {
			case <-cCtx.Done():
				return
			case name := <-data:
				match, _ := regexp.MatchString(name, "/s+")
				if name != "" && !match {
					err := u.red.LPush(cCtx, "names", name).Err()
					u.red.Expire(cCtx, "names", 10*time.Second)
					log.Println("adding %s to redis cache", name)
					if err != nil {
						log.Println("name: %s push fail to cache cause of %s", name, err)
						wg.Done()
						continue
					}
					log.Println("done!")
				}
				wg.Done()
				continue
			}
		}
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data <- scanner.Text()
		wg.Add(1)
	}
	wg.Wait()
	file.Close()
	cancel()
	close(data)
	log.Println("done scanning and add to cache!")
	return nil
}
