package db

import (
	"errors"
	"log"
	"time"

	pb "github.com/Sotatek-ThomasNguyen2/example-service.git/proto"
	"github.com/Sotatek-ThomasNguyen2/example-service.git/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	engine *gorm.DB
}

// ConnectDb open connection to db
// ConnectDb expose ...
func (d *DB) ConnectDb(sqlDSN string) error {
	db, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN:                  sqlDSN,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			// NamingStrategy: schema.NamingStrategy{
			// 	TablePrefix:   "partner.", // schema name
			// 	SingularTable: false,
			// },
			Logger: logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		return err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			<-ticker.C
			if err := sqlDb.Ping(); err != nil {
				log.Print(err)
			}
		}
	}()
	// // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	// sqlDb.SetMaxIdleConns(10)

	// // SetMaxOpenConns sets the maximum number of open connections to the database.
	// sqlDb.SetMaxOpenConns(100)
	d.engine = db
	return nil
}

func (d *DB) listUserPartnerQuery(req *pb.UserPartnerRequest) *gorm.DB {
	ss := d.engine.Table(tblUserPartner)
	if req.GetUserId() != "" {
		ss.Where(tblUserPartner+".user_id = ?", req.GetUserId())
	}
	if req.GetPhone() != "" {
		ss.Where(tblUserPartner+".phone = ?", req.GetPhone())
	}
	if req.GetUserId() != "" {
		ss.Where(tblUserPartner+".user_id = ?", req.GetUserId())
	}
	if req.GetPhone() != "" {
		ss.Where(tblUserPartner+".phone = ?", req.GetPhone())
	}
	return ss
}

func (d *DB) ListUserPartner(rq *pb.UserPartnerRequest) ([]*pb.UserPartner, error) {
	log.Println("req: ", rq)
	var userParters []*pb.UserPartner
	ss := d.listUserPartnerQuery(rq)
	if rq.GetLimit() != 0 {
		ss.Limit(int(rq.GetLimit()))
	} else {
		if rq.GetLimit() != 0 {
			ss.Limit(int(rq.GetLimit()))
		}
	}
	err := ss.Order("created desc").Find(&userParters).Error
	if err != nil {
		return nil, err
	}
	return userParters, nil
}

// FindUserVoucher expose ...
func (d *DB) FindUserPartner(rq *pb.UserPartnerRequest) (*pb.UserPartner, error) {
	raw := &pb.UserPartner{}
	if err := d.listUserPartnerQuery(rq).Take(raw).Error; err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, errors.New(utils.E_not_found)
	}
	return raw, nil
}

// InsertUserVoucher expose ...
func (d *DB) InsertUserPartner(up ...*pb.UserPartner) error {
	tx := d.engine.Table(tblUserPartner).Create(up)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// UpdateUserVoucher expose ...
func (d *DB) UpdateUserPartner(updator, selector *pb.UserPartner) error {
	tx := d.engine.Table(tblUserPartner).Model(selector).Updates(updator)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		log.Println("update no affected")
	}
	return nil
}
