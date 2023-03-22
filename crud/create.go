package crud

import (
	"database/sql"
	"fmt"
	gomysql "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

const (
	sqliteDsn = "test_create.db"
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

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// ShowCreate show db creation cases
func ShowCreate() {
	initDB()
	defer closeDB()
	pingDB()

	createRecord()
	createRecordWithSelectedFields()
	batchInsert()
	createFromMap()
}

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

func createRecord() {
	album := Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	}

	result := db.Table(albumTable).Create(&album)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Printf("createRecord RowsAffected: %d, user.ID: %d\n", result.RowsAffected, album.ID)
}

func createRecordWithSelectedFields() {
	album := Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	}

	// INSERT INTO `album` (`title`, `artist`, `price`)
	// VALUES ("The Modern Sound of Betty Carter", "Betty Carter", 49.99)
	result := db.Table(albumTable).Select("Title", "Artist", "Price").Create(&album)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Printf("RowsAffected: %d, user.ID: %d\n", result.RowsAffected, album.ID)
}

func batchInsert() {
	var albums = []Album{
		{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{Title: "Sarah Vaughan", Artist: "Sarah Vaughan", Price: 34.98},
	}

	// result := db.Create(&albums)
	batchSize := len(albums) / 2
	result := db.Table(albumTable).CreateInBatches(&albums, batchSize)
	if result.Error != nil {
		log.Fatal(result)
	}

	for _, album := range albums {
		log.Println("batchInsert", album.ID)
	}

	// // Skip hooks
	// result = db.Table(albumTable).Session(&gorm.Session{SkipHooks: true}).Omit("ID").CreateInBatches(&albums, len(albums)/2)
	// if result.Error != nil {
	// 	log.Fatal(result)
	// }
	//
	// for _, album := range albums {
	// 	log.Println("batchInsert skip hooks", album.ID)
	// }
}

// BeforeCreate Create Hooks
func (alb *Album) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().ID()
	fmt.Println("BeforeCreate UUID.ID() is", id)
	return
}

func createFromMap() {
	album := map[string]interface{}{
		"Title":  "Giant Steps",
		"Artist": "John Coltrane",
		"Price":  34.98,
	}

	result := db.Table(albumTable).Create(album)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Println("createFromMap album:", album)

	// when creating from map, hooks won't be invoked, association won't be saved
	// and primary key values won't be backfilled.
	albums := []map[string]interface{}{
		{"Title": "Jeru", "Artist": "Gerry Mulligan", "Price": 17.99},
		{"Title": "Sarah Vaughan", "Artist": "Sarah Vaughan", "Price": 34.98},
	}

	result = db.Table(albumTable).Create(albums)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Println("createFromMap albums:", albums)
}
