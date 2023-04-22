package associations

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID uint
	Price  float64
}

type User3 struct {
	gorm.Model
	Name   string
	Orders []Order
}

func FindUserPreloadOrders(db *gorm.DB) {
	var users []User3
	db.Preload("Orders").Find(&users)
	// select * from users;
	// select * from orders where user_id in (1,2,3,...);
}

// Joins Preloading

// Preload All

// Preload with conditions

// Custom Preloading SQL

// Nested Preloading

// Embedded Preloading
