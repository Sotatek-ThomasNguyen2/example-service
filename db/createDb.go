package db

import (
	pb "github.com/Sotatek-ThomasNguyen2/example-service.git/proto"
	"gorm.io/gorm"
)

const (
	tblUserPartner = "user_partner"
)

func createTable(model interface{}, tblName string, engine *gorm.DB) error {
	err := engine.Table(tblName).Migrator().AutoMigrate(model)
	return err
}

func (d *DB) CreateDb() error {
	var err error
	err = createTable(&pb.UserPartner{}, tblUserPartner, d.engine)
	return err
}
