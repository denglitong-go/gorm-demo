package crud

import (
	"gorm.io/gorm"
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
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	log.Println("before update interceptor")
	return
}

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
