package quick_start

import (
	"database/sql"
	gomysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	dbAddr       = "127.0.0.1:3306"
	dbName       = "recordings"
	dbUser       = "DB_USER"
	dbPassword   = "DB_PASSWORD"
	dbDriverName = "mysql"
	gormDB       *gorm.DB
	sqlDB        *sql.DB
)

// ShowGORMConnectToMysql show GORM connecting to gormDB.
// DB_USER=root DB_PASSWORD=12345678 go run .
func ShowGORMConnectToMysql() {
	cfg := &gomysql.Config{
		User:   os.Getenv(dbUser),
		Passwd: os.Getenv(dbPassword),
		Net:    "tcp",
		Addr:   dbAddr,
		DBName: dbName,
		// to handle time.Time correctly
		ParseTime: true,
		// To specify charset=utf8mb4_unicode_520_ci to fully support UTF-8; default value: utf8mb4_general_ci.
		// All these collations are for the UTF-8 character encoding.
		// The differences are in how text is sorted and compared.
		Collation: "utf8mb4_unicode_520_ci",
	}
	log.Println("Connecting data source name:", cfg.FormatDSN())

	var err error

	gormDB, err = gorm.Open(
		mysql.New(mysql.Config{
			DSN:        cfg.FormatDSN(),
			DriverName: dbDriverName,
		}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	setupSqlDB()
	defer close()
	ping()
}

func setupSqlDB() {
	var err error
	sqlDB, err = gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
}

func close() {
	if err := sqlDB.Close(); err != nil {
		log.Fatal(err)
	}
}

func ping() {
	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Printf("GORM %s Connected!\n", dbDriverName)
}
