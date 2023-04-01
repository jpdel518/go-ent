package mysql

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jpdel518/go-ent/ent"
	"github.com/jpdel518/go-ent/infrastructure/rdb/seed"
	"log"
	"os"
	"time"
)

type config struct {
	SQLDriver string
	DbName    string
	DbUser    string
	DbPass    string
}

var DB *sql.DB

func InitDatabase() {
	c := config{
		SQLDriver: os.Getenv("RDB_DRIVER"),
		DbName:    os.Getenv("RDB_NAME"),
		DbUser:    os.Getenv("RDB_USER"),
		DbPass:    os.Getenv("RDB_PASSWORD"),
	}

	DB, err := sql.Open(c.SQLDriver, c.DbUser+":"+c.DbPass+"@tcp(mysql:3306)/"+c.DbName+"?charset=utf8mb4&parseTime=True")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	// migration
	// migrationファイルを作成する場合はコメントを外してSchema.WriteToを実行する
	// file, err := os.Create("migration")
	// if err != nil {
	// 	log.Fatalf("failed creating migration file: %v", err)
	// }
	// defer func(file *os.File) {
	// 	if err := file.Close(); err != nil {
	// 		log.Printf("failed closing migration file: %v", err)
	// 	}
	// }(file)
	// 起動直後migrationが失敗するので、失敗したら1秒待って再度実行する
	client := NewClient()
	count := 0
	for {
		err := client.Schema.Create(context.Background())
		if err != nil {
			count++
			time.Sleep(1 * time.Second)
			log.Printf("migration failed count: %d", count)
			if count > 30 {
				log.Fatalf("failed creating schema resources: %v", err)
			}
			continue
		}
		break
	}

	// seed
	seed.CarSeed(client)
}

func NewClient() *ent.Client {
	drv := entsql.OpenDB("mysql", DB)
	client := ent.NewClient(ent.Driver(drv))

	// defer func(Client *ent.Client) {
	// 	if err := Client.Close(); err != nil {
	// 		log.Printf("failed closing ent client: %v", err)
	// 	}
	// }(Client)

	// デバッグモードを利用
	env := os.Getenv("ENV")
	if env != "staging" && env != "production" {
		client = client.Debug()
	}

	return client
}
