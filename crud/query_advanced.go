package crud

import (
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
	"log"
)

type User struct {
	ID       uint
	Name     string
	Age      uint
	Gender   string
	Password []byte
	Role     string

	// if you don't want to include gorm.Model, you can enable the soft delete feature like:
	DeletedAt gorm.DeletedAt

	// hundreds of fields
}

type APIUser struct {
	ID   uint
	Name string
}

func smartSelectFields() {
	var apiUser APIUser
	// select id, name automatically when querying
	// SELECT id, name FROM users LIMIT 10
	db.Table(albumTable).Find(&apiUser).Limit(10)

	var users []User
	// QueryFields mode will select by all fields' name for current model
	// select all columns
	db.Session(&gorm.Session{QueryFields: true}).Find(&users)
}

func forUpdate() {
	var users []User
	// SELECT * FROM users FOR UPDATE
	db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
}

func subQuery() {
	type Order struct {
		// fields
	}
	var orders []Order
	// select * from orders where amount > (select avg(amount) from orders)
	db.Where("amount > (?)",
		db.Table("orders").Select("AVG(amount)"),
	).Find(&orders)

	subQuery := db.Select("AVG(age)").Where("name like ?", "var%").Take("users")
	// select AVG(age) as average from users group by name having AVG(age) > (select AVG(age) where name like 'var%')
	db.Select("AVG(age) as average").Group("name").Having("AVG(age) > (?)", subQuery)
}

// GORM allows you using subquery in FROM clause with the method Table
func fromSubQuery() {
	var user User
	// select * from (select name, age from users) as u where age = 18
	db.Table("(?) as u", db.Model(&User{}).Select("name", "age")).Where("age = ?", 18).Find(&user)

	subQuery1 := db.Model(&User{}).Select("name")
	subQuery2 := db.Model(&User{}).Select("name")
	// select * from (select name from users) as u, (select name from pets) as p
	db.Table("(?) as u, (?) as p", subQuery1, subQuery2).Find(&user)
}

func groupConditions() {
	type Pizza struct {
		// fields
	}
	// select * from pizzas where
	//  (pizza = "pepperoni") and (size = "small" or size = "medium")
	// or
	//  (pizza = "hawaiian" and size = "xlarge")
	db.Where(
		db.Where("pizza = ?", "pepperoni").
			Where(
				db.Where("size = ?", "small").Or("size = ?", "medium"),
			),
	).Or(
		db.Where("pizza = ?", "hawaiian").Where("size = ?", "xlarge"),
	).Find(&Pizza{})
}

// Selecting IN with multiple columns
func inWithMultipleColumns() {}

// Named Argument
func namedArgument() {
	var users []User
	db.Where("name1 = @name or name2 = @name",
		sql.Named("name", "somebody"),
		sql.Named("other", "value"),
	).First(&users)
}

// FindToMap GORM allows scanning results to map[string]interface{} or []map[string]interface{},
// don't forget to specify Model or Table.
func FindToMap() {
	result := map[string]interface{}{}
	db.Model(&User{}).First(&result, "id = ?", 1)
}

// FirstOrInit Get first matched record or initialize a new instance with given conditions,
// only works with struct or map conditions
func FirstOrInit() {
	var user User
	db.FirstOrInit(&user, User{Name: "non_existing"})

	db.Where(User{Name: "somebody"}).FirstOrInit(&user)

	db.FirstOrInit(&user, map[string]interface{}{"name": "somebody"})

	// Initialize struct with more attributes if record not found,
	// those Attrs won't be used to build the SQL query.

	// select * from users where name = 'non_existing' order by id limit 1
	// insert user -> User{name: "non_existing", age: 20}
	db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrInit(&user)
	db.Where(User{Name: "non_existing"}).Attrs("age", 20).FirstOrInit(&user)
	// select * from users where name = 'somebody' order by id limit 1
	// found user -> User{name: "somebody", age: 18}, will ignore Attrs.
	db.Where(User{Name: "somebody"}).Attrs("age", 20).FirstOrInit(&user)

	// Assign attributes to struct regardless it is found or not, those attributes won't
	// be used to build SQL query and the final data won't be saved into database.

	// User not found, initialize it with give conditions and Assign attruibutes
	// insert user -> User{name: "non_existing", age: 20}
	db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrInit(&user)
	// found user -> User{name: "somebody", age: 18}
	// will return -> User{name: "somebody", age: 20}
	db.Where(User{Name: "somebody"}).Assign(User{Age: 20}).FirstOrInit(&user)
}

// FirstOrCreate Get first matched record or create a new one with given conditions
// (only works with struct, map conditions), RowsAffected returns created/updated record's count.
func FirstOrCreate() {
	var user User
	// Create struct with more attributes if record not found, those Attrs won't be used to build SQL query.
	db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
	// user -> User{ID: 112, Name: "non_existing", Age: 20}
	// result.RowsAffected // => 1

	// Found user with `name` = `somebody`
	db.Where(User{Name: "somebody"}).FirstOrCreate(&user)
	// user -> User{ID: 111, Name: "somebody", "Age": 18}
	// result.RowsAffected // => 0

	// User not found, initialize it with give conditions and Assign attributes
	db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
	// user -> User{ID: 112, Name: "non_existing", Age: 20}

	// Found user with `name` = `somebody`
	db.Where(User{Name: "somebody"}).Assign(User{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'somebody' ORDER BY id LIMIT 1;
	// UPDATE users SET age=20 WHERE id = 111;
	// user -> User{ID: 111, Name: "somebody", Age: 20}
}

func OptimizerAndIndexHints() {
	var user User
	// SELECT * /*+ MAX_EXECUTION_TIME(10000) */ FROM `users`
	db.Clauses(hints.New("MAX_EXECUTION_TIME(1000)")).Find(&user)

	// SELECT * FROM `users` USE INDEX (`idx_user_name`)
	db.Clauses(hints.UseIndex("idx_user_name")).Find(&user)

	// select * from users force index for join (idx_user_name, idx_user_id)
	db.Clauses(hints.UseIndex("idx_user_name", "idx_user_id").ForJoin()).Find(&user)
}

func Iteration() {
	rows, err := db.Model(&User{}).Where("name = ?", "somebody").Rows()
	log.Fatal(err)
	defer rows.Close()

	for rows.Next() {
		var user User
		db.ScanRows(rows, &user)
		// ...
	}
}

func FindInBatches() {
	// batch size 100
	type Result struct {
		// fields
	}
	var results []Result
	db.Where("processed = ?", false).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			// batch processing found records
			log.Println(result)
		}

		tx.Save(&results)
		log.Println(tx.RowsAffected, batch)

		// returns error will stop future batches
		return nil
	})
}

// AfterFind Query Hooks
func (u User) AfterFind(tx *gorm.DB) (err error) {
	log.Println("after find", u.ID)
	return
}

// Pluck Query single column from database and scan into a slice,
// if you want to query multiple columns, use Select with Scan/Find instead
func Pluck() {
	var ages []int64
	db.Model(&User{}).Pluck("age", &ages)

	var names []string
	// distinct pluck
	db.Model(&User{}).Distinct().Pluck("name", &names)

	type NameAge struct {
		Name string
		Age  int64
	}
	var nameAges []NameAge
	db.Model(&User{}).Select("name", "age").Scan(&nameAges)
	db.Model(&User{}).Select("name", "age").Find(&nameAges)
}

func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	return db.Where("amount > ?", 1000)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode_sign = ?", "CreditCard")
}

func PaidWithCreditCrash(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode_sign = ?", "Crash")
}

// If you need to accept runtime arguments, then you can return func
func OrderStatus(status []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status IN (?)", status)
	}
}

// Scopes allows you to specify commonly-used queries which can be referenced as method calls
func Scopes() {
	type Order struct{}
	var orders []Order
	// Find all credit card orders and amount greater than 1000
	db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)

	// Find all crash orders and amount greater than 1000
	db.Scopes(AmountGreaterThan1000, PaidWithCreditCrash).Find(&orders)

	// Find all paid, shipped orders that amount greater than 1000
	db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
}

func Count() {
	var count int64
	db.Model(&User{}).Where("name = ?", "somebody").Count(&count)

	db.Table(albumTable).Count(&count)

	db.Table(albumTable).Distinct().Count(&count)

	db.Table(albumTable).Select("count(distinct(name))").Count(&count)

	// count with group by
	users := []User{
		{Name: "name1"},
		{Name: "name2"},
		{Name: "name3"},
		{Name: "name3"},
	}
	log.Println("all users", users)
	// select count(1) from users group by name
	// count => 3
	db.Model(&User{}).Group("name").Count(&count)
}
