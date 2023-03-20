package crud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"time"
)

const (
	sqliteDsn = "test_create.db"
)

var (
	db *gorm.DB
)

// User equals to User
type User struct {
	gorm.Model
	Name        string
	Email       string
	Age         uint
	Birthday    time.Time
	ActivatedAt time.Time
	Location    Location
}

// Location Create from customized data type
//
type Location struct {
	X, Y int
}

// Scan implements the sql.Scanner interface
func (loc *Location) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	bytes, ok := v.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("Failed to unmarshal JSONB value: %v", v))
	}
	location := Location{}
	err := json.Unmarshal(bytes, &location)
	*loc = location
	return err
}

func (loc Location) GormDataType() string {
	return "geometry"
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		// Not supported in sqlite
		SQL:  "ST_PointFromText(?)",
		Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y)},
	}
}

func ShowCreate() {
	initDB()
	defer closeDB()

	initSchema()
	createRecord()
	createRecordWithSelectedFields()
	batchInsert()
	createFromMap()
	createFromSQLExpressionOrContextValuer()
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(sqliteDsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

func closeDB() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal(err)
	}
	if err = os.Remove(sqliteDsn); err != nil {
		log.Fatal(err)
	}
}

func initSchema() {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
}

func createRecord() {
	user := User{Name: "Litong", Age: 18, Birthday: time.Now()}

	result := db.Create(&user)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Printf("createRecord RowsAffected: %d, user.ID: %d\n", result.RowsAffected, user.ID)
}

func createRecordWithSelectedFields() {
	user := User{Name: "Litong", Age: 18, Birthday: time.Now()}

	// INSERT INTO `user_models` (`name`, `age`) VALUES ("Litong", 18)
	result := db.Select("Name", "Age").Create(&user)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Printf("RowsAffected: %d, user.ID: %d\n", result.RowsAffected, user.ID)

	// insert into table with other fields except that `id`, `name`, `age`
	result = db.Omit("ID", "Name", "Age").Create(&user)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Printf("createRecordWithSelectedFields RowsAffected: %d, user.ID: %d", result.RowsAffected, user.ID)
}

func batchInsert() {
	var users = []User{
		{Name: "user1"},
		{Name: "user2"},
		{Name: "user3"},
	}

	// result := db.Create(&users)
	result := db.CreateInBatches(&users, len(users)/2)
	if result.Error != nil {
		log.Fatal(result)
	}

	for _, user := range users {
		log.Println("batchInsert", user.ID)
	}

	// Skip hooks
	result = db.Session(&gorm.Session{SkipHooks: true}).Omit("ID").CreateInBatches(&users, len(users)/2)
	if result.Error != nil {
		log.Fatal(result)
	}

	for _, user := range users {
		log.Println("batchInsert skip hooks", user.ID)
	}
}

// BeforeCreate Create Hooks
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().ID()
	fmt.Println("BeforeCreate UUID.ID() is", id)
	return
}

func createFromMap() {
	user := map[string]interface{}{
		"Name": "Litong",
		"Age":  18,
	}

	result := db.Model(&User{}).Create(&user)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Println("createFromMap user:", user)

	// when creating from map, hooks won't be invoked, association won't be saved
	// and primary key values won't be backfilled.
	users := []map[string]interface{}{
		{"Name": "user3", "Age": 19},
		{"Name": "user4", "Age": 20},
	}

	result = db.Model(&User{}).Create(&users)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Println("createFromMap users:", users)
}

func createFromSQLExpressionOrContextValuer() {
	// create from map
	// INSERT INTO `users` (`name`, `location`) VALUES ("Litong", ST_PointFromText(POINT(100 100)))
	result := db.Model(&User{}).Create(map[string]interface{}{
		"Name": "Litong",
		"Location": clause.Expr{
			SQL:  "ST_PointFromText(?)",
			Vars: []interface{}{"POINT(100 100)"},
		},
	})
	if result.Error != nil {
		log.Fatal(result)
	}
	log.Println(result)
}
