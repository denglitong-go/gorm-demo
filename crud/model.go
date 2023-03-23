package crud

import (
	"database/sql"
	gomysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	dbUser      = "DB_USER"
	dbPassword  = "DB_PASSWORD"
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
	albumTable = "album"

	db      *gorm.DB
	mysqlDB *sql.DB
)

func initDB() {
	// Get a database handle
	var err error
	// root:12345678@tcp(127.0.0.1:3306)/recordings?allowNativePasswords=false&checkConnLiveness=false&maxAllowedPacket=0
	log.Println("config data source name", mysqlConfig.FormatDSN())
	db, err = gorm.Open(mysql.Open(mysqlConfig.FormatDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	mysqlDB, err = db.DB()
	if err != nil {
		log.Fatal(err)
	}
}

func closeDB() {
	if err := mysqlDB.Close(); err != nil {
		log.Fatal(err)
	}
}

func pingDB() {
	if err := mysqlDB.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("GORM connected Mysql!")
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	// Any zero value like 0, '', false won't be saved into the database for those fields
	// defined default value, you might want to use pointer type or Scanner/Valuer to avoid this.
	Price    float32  `gorm:"default:1.00"`
	Location Location `gorm:"type:json;serializer:json"`
}

// Location Create from customized data type
type Location struct {
	X int
	Y int
}

// // Scan implements the sql.Scanner interface
// func (loc *Location) Scan(v interface{}) error {
// 	if v == nil {
// 		fmt.Println("Scan value is nil")
// 		return nil
// 	}
//
// 	// Scan a value into struct from database driver
// 	mysqlEncoding, ok := v.([]byte)
// 	if !ok {
// 		return errors.New(fmt.Sprintf("did not scan: expected []byte but was: %T", v))
// 	}
// 	fmt.Println("mysqlEncoding:", mysqlEncoding)
// 	// TODO Have no idea how to convert the mysql Point representation of bytes to Go struct.
// 	return nil
// }
//
// func (loc Location) GormDataType() string {
// 	return "geometry"
// }
//
// func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
// 	return clause.Expr{
// 		SQL:  "ST_PointFromText(?)",
// 		Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y)},
// 	}
// }
