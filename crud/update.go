package crud

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

func SaveAllFields() {
	var user User
	db.First(&user)

	user.Name = "somebody"
	user.Age = 12
	// update users set name = 'somebody', age = 12, ...with all other original fields value
	// where id = xx
	db.Table(albumTable).Save(&user)

	// insert into users (name, age) value ('newbody', 100)
	db.Save(&User{Name: "newbody", Age: 100})
	// update users set name = xx, age = xx where id = 1
	db.Save(&User{ID: 1, Name: "newbody", Age: 100})

	// Don't use Save with Model, it's an undefined behavior
}

func UpdateSingleColumn() {
	// update users set name = 'hello', updated_at='' where active=true
	db.Model(&User{}).Where("active = ?", true).Update("name", "hello")

	user := User{ID: 111}
	// update users set name = 'hello' where id = 111
	db.Model(&user).Update("name", "hello")
	// update users set name = 'hello' where id = 111 and active = true
	db.Model(&user).Where("active = ?", true).Update("name", "hello")
}

func UpdatesMultipleColumns() {
	user := User{ID: 111}

	// Update attributes with struct, will only update non-zero fields
	// update users set name = 'hello', age=18 where id = 111
	db.Model(&user).Updates(User{Name: "hello", Age: 18})

	db.Model(&user).Updates(map[string]interface{}{
		"name": "hello",
		"age":  18,
	})
}

func UpdateSelectedFields() {
	user := User{ID: 111}

	// update users set name = 'hello' where id = 111
	db.Model(&user).Select("name").Updates(map[string]interface{}{
		"name": "hello", "age": 18, "active": false,
	})

	// UPDATE users SET age=18, active=false WHERE id=111;
	db.Model(&user).Omit("name").Updates(map[string]interface{}{
		"name": "hello", "age": 18, "active": false,
	})

	// update users set name = 'hello', age=0 where id = 111
	db.Model(&user).Select("name", "age").Updates(map[string]interface{}{
		"name": "hello", "age": 0,
	})

	// select all fields to update
	// update users set name = 'hello', age=0, active=false where id = 111
	db.Model(&user).Select("*").Updates(map[string]interface{}{
		"name": "hello", "age": 0, "active": false,
	})

	// update users set name = 'hello', age=0 where id = 111
	db.Model(&user).Select("*").Omit("active").Updates(map[string]interface{}{
		"name": "hello", "age": 0, "active": false,
	})
}

// BeforeUpdate Update Hookds
// func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
// 	log.Println("before update interceptor")
// 	return
// }

// If we haven't specified a record primary key, GORM will perform a batch update
func BatchUpdates() {
	// update users set name = 'hello', age=18 where role = admin
	db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})

	// update users set name = 'hello', age=18 where id = [1,2,3]
	db.Table(albumTable).Where("id in ?", []int64{1, 2, 3}).Updates(map[string]interface{}{
		"name": "hello", "age": 18,
	})
}

func BlockGlobalUpdates() {
	// if you perform a batch update without any conditions, GORM will return error = gorm.ErrMissingWhereClause
	db.Model(&User{}).Update("name", "hello")

	// update users set name = 'hello' where 1=1
	db.Model(&User{}).Where("1 = 1").Update("name", "hello")

	// update users set name = 'hello'
	db.Exec("update users set name = ?", "hello")

	// update users set name = 'hello'
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&User{}).Update("name", "hello")
}

// Advanced update

func UpdateWithSQLExpression() {
	type Product struct {
		ID       int64
		Price    float64
		Quantity int64
		// other fields
	}
	product := Product{ID: 3}

	// update products set price = price * 2 + 100 where id = 3
	db.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))
	db.Model(&product).Updates(map[string]interface{}{
		"price": gorm.Expr("price * ? + ?", 2, 100),
	})

	// update products set quantity = quantity - 1 where id = 3
	db.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
	// update products set quantity = quantity - 1 where id = 3 and quantity > 1
	db.Model(&product).Where("quantity > ?", 1).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
}

func UpdateFromSubQuery() {
	type Company struct{}
	user := User{ID: 111}
	// update users set company_name = (select name from companies where companies.id = users.company_id)
	db.Model(&user).Update("company_name",
		db.Model(&Company{}).Select("name").Where("companies.id = users.company_id"),
	)

	// update users set
	// 	company_name = (select name from companies where companies.id = users.company_id)
	// 	where name = 'someone'
	db.Table("users as u").Where("name = ?", "someone").Update(
		"company_name",
		db.Table("companies as c").Select("name").Where("c.id = u.company_id"),
	)
	db.Table("users as u").Where("name", "someone").Updates(map[string]interface{}{
		"company_name": db.Table("companies as c").Select("name").Where("c.id = u.company_id"),
	})
}

// WithoutHooksOrTimeTrack If you want to skip Hooks and don't track the update time when updating, you can use
// UpdateColumn, UpdateColumns, it works like Update, Updates
func WithoutHooksOrTimeTrack() {
	user := User{ID: 111}
	// update single column
	db.Model(&user).UpdateColumn("name", "someone")
	// update multiple columns
	db.Model(&user).UpdateColumns(User{Name: "someone", Age: 18})
	// update selected columns
	// update users set name = 'hello', age = 0 where id = 111
	db.Model(&user).Select("name", "age").UpdateColumns(User{
		Name:   "someone",
		Age:    18,
		Gender: "Male",
	})
}

func ReturningDataFromModifiedRows() {
	var users []User
	// update users set salary = salary * 2 where role = 'admin' returning *
	// users => []User{{ID:1,Name:"someone",Role:"admin",Salary:100},...}
	db.Model(&users).Clauses(clause.Returning{}).Where("role = ?", "admin").
		Update("salary", gorm.Expr("salary * ?", 2))

	// return specified columns
	// update users set salary = salary * 2 where role = 'admin' returning name
	// users => []User{{ID:0,Name:"someone",Role:"",Salary:100},...}
	db.Model(&users).Clauses(clause.Returning{Columns: []clause.Column{{Name: "name"}}}).
		Where("role = ?", "admin").
		Update("salary", gorm.Expr("salary * ?", 2))
}

// Check Field has changed

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	// if Role changed
	if tx.Statement.Changed("Role") {
		return errors.New("role not allowed to change")
	}
	if !tx.Statement.Changed() {
		log.Println("not any fields change")
	}
	return nil
}

// Change Updating Values.
// To change updating values in Before Hooks, you should use SetColumn unless it's a full
// update with Save.
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if pw, err := bcrypt.GenerateFromPassword(u.Password, 0); err == nil {
		tx.Statement.SetColumn("EncryptedPassword", pw)
	}
	if tx.Statement.Changed("Code") {
		u.Age += 20
		tx.Statement.SetColumn("Age", u.Age)
	}
	return nil
}
