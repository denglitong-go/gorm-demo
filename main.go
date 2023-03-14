// Package main
// go get -u gorm.io/gorm
// go get -u gorm.io/driver/sqlite
package main

import "go-orm/quick_start"

func main() {
	quick_start.ShowQuickStart()
	quick_start.ShowGORMConnectToMysql()
	quick_start.ShowGORMConnectToPostgreSQL()
}
