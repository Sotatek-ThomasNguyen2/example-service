package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/Sotatek-ThomasNguyen2/example-service.git/db"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

var Global = 1

func main() {
	go freeMemory()
	if err := CliTool(); err != nil {
		log.Panicln(err)
	}
}

func freeMemory() {
	for {
		fmt.Println("run gc")
		start := time.Now()
		runtime.GC()
		debug.FreeOSMemory()
		elapsed := time.Since(start)
		fmt.Printf("gc took %s\n", elapsed)
		time.Sleep(15 * time.Minute)
	}
}

func createTableDb(ctx *cli.Context) error {
	d := &db.DB{}
	if err := d.ConnectDb("host=localhost user=postgres password=123 dbname=example_service port=5432 sslmode=disable"); err != nil {
		debug.PrintStack()
		log.Panicln(err)
	}
	if err := d.CreateDb(); err != nil {
		return err
	}
	log.Println("create table successfully!")
	return nil
}

func start(ctx *cli.Context) error {
	u := initServe()
	r := &Router{
		route: gin.Default(),
	}
	r.u = u
	err := r.HttpRouter(u)
	if err != nil {
		log.Println("err: ", err)
		return err
	}
	if err := StartGRPCServe(3000, u); err != nil {
		log.Println("err: ", err)
		return err
	}
	return nil
}

func scan(ctx *cli.Context) error {
	u := initServe()
	if err := u.ScanToCache(); err != nil {
		return err
	}
	return nil
}

func CliTool() error {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		return errors.New("what?")
	}
	app.Commands = []*cli.Command{
		{Name: "start", Action: start},
		{Name: "scan", Action: scan},
		{Name: "createDb", Usage: "Creating database table", Action: createTableDb},
	}
	return app.Run(os.Args)
}
