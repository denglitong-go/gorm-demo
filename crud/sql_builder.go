package crud

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
	"time"
)

// RawSQL Query raw sql with Scan
func RawSQL() {
	type Result struct {
		ID   uint
		Name string
		Age  uint
	}
	var result Result
	db.Raw("select id, name, age from users wehre id = ?", 3).Scan(&result)
	db.Raw("select id, name, age from users wehre name = ?", "someone").Scan(&result)

	var age int
	db.Raw("select avg(age) from users where role = ?", "admin").Scan(&age)

	var users []User
	db.Raw("update users set name = ? where age = ? returning id, name", "someone", 3).Scan(&users)

	db.Exec("drop table users")
	db.Exec("update orders set shipped_at = ? where id in ?", time.Now(), []int64{1, 2, 3})
	db.Exec("update users set money = ? where name = ?", gorm.Expr("money * ? + ?", 2, 1), "someone")
}

func NamedArgument() {
	var user User
	db.Where("name1 = @name or name2 = @name", sql.Named("name", "someone")).Find(&user)
	db.Where("name1 = @name or name2 = @name", map[string]interface{}{
		"name": "someone",
	}).Find(&user)

	db.Raw("select * from users where name1 = @name or name2 = @name2",
		sql.Named("name", "someone"),
		sql.Named("name2", "otherone"),
	).Find(&user)

	db.Exec("update users set name1 = @name, name2 = @name2",
		sql.Named("name", "someone"),
		sql.Named("name2", "otherone"),
	)

	db.Raw("select * from users where name1 = @name or name2 = @name2", map[string]interface{}{
		"name":  "someone",
		"name2": "otherone",
	}).Find(&user)

	type NamedArgument struct {
		Name  string
		Name2 string
	}
	db.Raw("select * from users where name1 = @name or name2 = @name2", NamedArgument{
		Name:  "someone",
		Name2: "otherone",
	}).Find(&user)
}

// DryRunMode Generate SQL and its arguments without executing,
// can be used to prepare or test generated SQL.
func DryRunMode() {
	var user User
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
	fmt.Println(stmt.SQL.String()) // select * from users where id = 1
	fmt.Println(stmt.Vars)         // []interface{}[1]
}

// ToSQL Returns generated SQL without executing.
// GORM uses the database/sql's argument placeholders to construct the SQL statement,
// which will automatically escape arguments to avoid SQL injection, but generated SQL don't
// provide the safety guarantees, please only use it for debugging.
func ToSQL() {
	id, limit := 100, 10
	generatedSQL := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&User{}).Where("id > ?", id).Limit(limit).Order("age desc").Find(&[]User{})
	})
	// select * from users where id > 100 and deleted_at is null order by age desc limit 10
	fmt.Println(generatedSQL)
}

func RowAndRows() {
	var name string
	var age int
	// user GORM API to build sql
	row := db.Table("users").Where("name = ?", "someone").Select("name", "age").Row()
	row.Scan(&name, &age)

	// use raw sql
	row = db.Raw("select name, age from users where name = ?", "someone").Row()
	row.Scan(&name, &age)

	rows, _ := db.Model(&User{}).Where("name = ?", "some").Select("name", "age").Rows()
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&name, &age)
	}
}

func ScanRowsToStruct() {
	rows, _ := db.Model(&User{}).Where("name = ?", "some").Select("name", "age").Rows()
	defer rows.Close()
	var user User
	for rows.Next() {
		db.ScanRows(rows, &user)
	}
}

// Connection Run multiple SQL in same db tcp connection (not in transaction)
func Connection() {
	db.Connection(func(tx *gorm.DB) error {
		tx.Exec("SET my.role = ?", "admin")
		tx.Exec("other sql")
		return nil
	})
}

// Clauses GORM uses SQL builder generates SQL internally, for each operation,
// GORM creates a *gorm.Statement object, all GORM APIs add/change Clause for the Statement,
// at last, GORM generated SQL based on those clauses.
func Clauses() {
	var user User
	var limit = 1
	db.Clauses(clause.Select{Columns: []clause.Column{{Name: "*"}}}).First(&user)
	db.Clauses(clause.From{Tables: []clause.Table{{Name: clause.CurrentTable}}}).First(&user)
	db.Clauses(clause.Limit{Limit: &limit}).First(&user)
	db.Clauses(clause.OrderBy{Columns: []clause.OrderByColumn{
		{
			Column: clause.Column{
				Table: clause.CurrentTable,
				Name:  clause.PrimaryKey,
			},
		},
	}}).First(&user)
}

func ClauseBuilder() {
	var user User
	// For different databases, Clauses may generate different SQL,
	// which is supported because GORM allows database driver register Clause Builder
	// to replace the default one.
	db.Offset(10).Limit(5).First(&user)
}

func ClauseOptions() {
	var users []User
	// insert ignore into users (...) values (...),(...)
	db.Clauses(clause.Insert{Modifier: "IGNORE"}).Save(&users)
}

func StatementModifier() {
	// select * /*+ hint */ from users
	db.Clauses(hints.New("hint")).Find(&User{})
}
