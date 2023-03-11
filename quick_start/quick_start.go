package quick_start

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func (p Product) String() string {
	return fmt.Sprintf("{Code: %s, Price: %d}", p.Code, p.Price)
}

func ShowQuickStart() {
	sqliteDsn := "test.db"
	db, err := gorm.Open(sqlite.Open(sqliteDsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal(err)
	}

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	// find product with integer primary key
	db.First(&product, 1)
	log.Println("product 1:", product)

	// find product with integer code D42
	db.First(&product, "code = ?", "D42")
	log.Println("product D42:", product)

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)

	db.First(&product, 1)
	log.Println("after update product 1:", product)

	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 300, Code: "F42"}) // non-zero fields

	db.First(&product, 1)
	log.Println("after update product 1:", product)

	db.Model(&product).Updates(map[string]interface{}{
		"Price": 400,
		"Code":  "G42",
	})

	db.First(&product, 1)
	log.Println("after update product 1:", product)

	// Delete - delete product
	db.Delete(&product, 1)

	if err = os.Remove(sqliteDsn); err != nil {
		log.Fatal(err)
	}
}
