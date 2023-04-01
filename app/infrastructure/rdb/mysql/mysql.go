package mysql

import (
	atlas "ariga.io/atlas/sql/migrate"
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jpdel518/go-ent/ent"
	"github.com/jpdel518/go-ent/ent/migrate"
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

func NewClient() *ent.Client {
	c := config{
		SQLDriver: os.Getenv("RDB_DRIVER"),
		DbName:    os.Getenv("RDB_NAME"),
		DbUser:    os.Getenv("RDB_USER"),
		DbPass:    os.Getenv("RDB_PASSWORD"),
	}

	client, err := ent.Open(c.SQLDriver, c.DbUser+":"+c.DbPass+"@tcp(mysql:3306)/"+c.DbName+"?charset=utf8mb4&parseTime=True")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	// defer func(client *ent.Client) {
	// 	if err := client.Close(); err != nil {
	// 		log.Printf("failed closing ent client: %v", err)
	// 	}
	// }(client)

	// デバッグモードを利用
	env := os.Getenv("ENV")
	if env != "staging" && env != "production" {
		client = client.Debug()
	}

	return client
}

func InitDatabase() {
	// デバッグモードを利用の場合は差分ファイルを作成
	ctx := context.Background()
	client := NewClient()
	env := os.Getenv("ENV")
	if env != "staging" && env != "production" {
		c := config{
			SQLDriver: os.Getenv("RDB_DRIVER"),
			DbName:    os.Getenv("RDB_NAME"),
			DbUser:    os.Getenv("RDB_USER"),
			DbPass:    os.Getenv("RDB_PASSWORD"),
		}

		// ローカルのent/migrateディレクトリに差分(migration)ファイルを作成
		dir, err := atlas.NewLocalDir("ent/migrate/migrations")
		if err != nil {
			log.Fatalf("failed creating atlas migration directory: %v", err)
		}
		// Migrate diff options.
		opts := []schema.MigrateOption{
			schema.WithDir(dir),                          // provide migration directory
			schema.WithMigrationMode(schema.ModeInspect), // provide migration mode
			schema.WithDialect(dialect.MySQL),            // Ent dialect to use
			schema.WithFormatter(atlas.DefaultFormatter),
		}
		var migrationName string
		if len(os.Args) != 2 {
			log.Println("default migration name is used. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
			migrationName = "create_schema"
		} else {
			migrationName = os.Args[1]
		}
		// 起動直後migrationが失敗するので、失敗したら1秒待って再度実行する
		count := 0
		for {
			// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
			err = migrate.NamedDiff(ctx, "mysql://"+c.DbUser+":"+c.DbPass+"@mysql:3306/"+c.DbName, migrationName, opts...)
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

		// 開発環境ではmigrationファイルではなく、entのauto migrateを利用する
		defer func(client *ent.Client) {
			if err := client.Close(); err != nil {
				log.Printf("failed closing ent client: %v", err)
			}
		}(client)
		if err := client.Schema.Create(ctx); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	}

	// seed
	seed.CarSeed(client)
}
