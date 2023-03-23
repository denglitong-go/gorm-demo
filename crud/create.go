package crud

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

// ShowCreate show db creation cases
func ShowCreate() {
	initDB()
	defer closeDB()
	pingDB()

	// createRecord()
	// createRecordWithSelectedFields()
	// batchInsert()
	// createFromMap()
	// createFromMapWithCustomStruct()
	createFromCustomStruct()
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

func createFromMapWithCustomStruct() {
	locationBytes, _ := json.Marshal(Location{X: 10, Y: 20})
	album := map[string]interface{}{
		"Title":    "Giant Steps",
		"Artist":   "John Coltrane",
		"Price":    34.98,
		"Location": locationBytes,
	}

	result := db.Table(albumTable).Create(&album)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Println("createFromMapWithCustomStruct album:", album)
}

func createFromCustomStruct() {
	album := Album{
		Title:    "Giant Steps",
		Artist:   "John Coltrane",
		Price:    34.98,
		Location: Location{X: 20, Y: 30},
	}

	result := db.Table(albumTable).Create(&album)
	if result.Error != nil {
		log.Fatal(result)
	}

	log.Println("createFromCustomStruct album:", album)
}
