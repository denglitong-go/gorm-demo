// Package main
// go get -u gorm.io/gorm
// go get -u gorm.io/driver/sqlite
package main

import (
	"go-orm/crud"
)

// DB_USER=root DB_PASSWORD=12345678 go run .
func main() {
	// quick_start.ShowQuickStart()
	// quick_start.ShowGORMConnectToMysql()
	// quick_start.ShowGORMConnectToPostgreSQL()
	// crud.ShowCreate()
	crud.ShowQuery()
}
