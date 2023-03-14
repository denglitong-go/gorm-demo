package quick_start

import (
	"database/sql"
	"fmt"
	gomysql "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"

	mysqlConfig = &gomysql.Config{
		User:   os.Getenv(dbUser),
		Passwd: os.Getenv(dbPassword),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
		// to handle time.Time correctly
		ParseTime: true,
		// To specify charset=utf8mb4_unicode_520_ci to fully support UTF-8; default value: utf8mb4_general_ci.
		// All these collations are for the UTF-8 character encoding.
		// The differences are in how text is sorted and compared.
		Collation: "utf8mb4_unicode_520_ci",
	}

	postgreSQLDSN = fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=recordings port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv(dbUser), os.Getenv(dbPassword),
	)

	mysqlDB    *sql.DB
	postgresDB *sql.DB
)

// ShowGORMConnectToMysql show GORM connecting to gormDB.
// DB_USER=root DB_PASSWORD=12345678 go run .
func ShowGORMConnectToMysql() {
	log.Println("Connecting data source name:", mysqlConfig.FormatDSN())

	// gormDB, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{})
	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:        mysqlConfig.FormatDSN(),
			DriverName: "mysql",
		}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	setupMysqlDB(gormDB)
	defer closeMysqlDB()
	pingMysqlDB()
}

func setupMysqlDB(gormDB *gorm.DB) {
	var err error
	mysqlDB, err = gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
}

func closeMysqlDB() {
	if err := mysqlDB.Close(); err != nil {
		log.Fatal(err)
	}
}

func pingMysqlDB() {
	if err := mysqlDB.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("GORM connected Mysql!")
}

// ShowGORMConnectToPostgreSQL show GORM connecting to PostgreSQL.
// DB_USER=root DB_PASSWORD=12345678 go run .
func ShowGORMConnectToPostgreSQL() {
	log.Println("Connecting data source name:", postgreSQLDSN)
	// gormDB, err := gorm.Open(postgres.Open(cfg.FormatDSN()), &gorm.Config{})
	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:        postgreSQLDSN,
			DriverName: "postgres",
		}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	setupPostgresDB(gormDB)
	defer closePostgresDB()
	pingPostgresDB()
}

func setupPostgresDB(gormDB *gorm.DB) {
	var err error
	postgresDB, err = gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
}

func closePostgresDB() {
	if err := postgresDB.Close(); err != nil {
		log.Fatal(err)
	}
}

func pingPostgresDB() {
	if err := postgresDB.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("GORM connected Postgres!")
}
