package crud

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/soft_delete"
	"time"
)

func DeleteARecord() {
	user := User{ID: 10}
	// delete from users where id = 10
	db.Delete(user)

	// delete from users where id = 10 and name = 'someone'
	db.Where("name = ?", "someone").Delete(&user)
}

func DeleteWithPrimaryKey() {
	// delete from users where id = 10
	db.Delete(&User{}, 10)
	db.Delete(&User{}, "10")
	// delete from users where id in [1,2,3]
	db.Delete(&User{}, []int{1, 2, 3})
}

// Delete Hooks
// BeforeDelete, AfterDelete
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Role == "admin" {
		return errors.New("admin user not allowed to delete")
	}
	return nil
}

// BatchDelete If the specified value has no primary value, GORM will perform a batch delete,
// it will delete all matched records
func BatchDelete() {
	// delete from users where name like '%someone%'
	db.Where("name like ?", "%someone%").Delete(&User{})
	db.Delete(&User{}).Where("name like ?", "%someone%")

	var users = []User{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}
	// delete from users where id in [1,2,3]
	db.Delete(&users)
	// delete from users where id in [1,2,3] and name like '%someone%'
	db.Delete(&users).Where("name like ?", "%someone%")
}

// BlockGlobalDelete
// If you perform a batch delete without any conditions, GORM won't run it,
// and will return ErrMissingWhereClause error.
// You have to use some conditions or use raw SQL, or enable AllowGlobalUpdate mode.
func BlockGlobalDelete() {
	// will get gorm.ErrMissingWhereClause
	db.Delete(&User{})
	// will get gorm.ErrMissingWhereClause
	db.Delete(&[]User{
		{Name: "someone1"},
		{Name: "someone2"},
	})

	// delete from users where 1 = 1
	db.Where("1 = 1").Delete(&User{})
	db.Exec("delete from users")
	// delete from users
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})
}

func ReturningDataFromDeletedRows() {
	var users []User
	// delete from users where role = 'admin' returning *
	// users => []User{{ID:1,Name:'someone',Role:'admin',...}...}
	db.Clauses(clause.Returning{}).Where("role = ?", "admin").Delete(&users)

	// delete from users where admin = 'admin' returning id
	// users => []User{{ID:1,...},{ID:2,...},...}
	db.Clauses(clause.Returning{Columns: []clause.Column{{Name: "ID"}}}).
		Where("role = ?", "admin").Delete(&users)
}

// SoftDelete If your model includes a `gorm.DeletedAt` field (which is included in gorm.model),
// it will get soft delete ability automatically!
// when calling Delete, the record won't be removed from the database, but GORM will set
// the `DeletedAt`'s value to the current time, and the data is not findable
// with normal Query methods anymore.
func SoftDelete() {
	user := User{ID: 111}
	// update users set deleted_at = '2023-04-02 09:51' where id = 111
	db.Delete(&user)

	// batch delete
	// delete from users where age = 20
	db.Where("age = ?", 20).Delete(&User{})
	// soft deleted records will be ignored when querying
	// select * from users where age = 20 and
	db.Where("age = 20").Find(&user)
}

func FindSoftDeletedRecords() {
	var user User
	// select * from users where age = 20
	db.Unscoped().Where("age = 20").Find(&user)
}

func DeletePermanently() {
	user := User{ID: 111}
	// delete from users where id = 111
	db.Unscoped().Delete(&user)
}

// DeleteFlag By default, gorm.Model users *time.Time as the value for the DeletedAt field,
// and it provides other data formats support with plugin gorm.io/plugin/soft_delete
func DeleteFlag() {
	type User struct {
		ID uint
		// when creating unique composite index for the DeletedAt field,
		// you must use other data format like unit second/flag
		// with plugin `gorm.io/plugin/soft_delete`'s help
		Name string `gorm:"uniqueIndex:udx_name"`
		// unix second as delete flag
		DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:udx_name"`
	}

	// query: select * from users where deleted_at = 0
	// update: update users set deleted_at = $currentUnixSecond where id = 1

	type UserMilli struct {
		ID        uint
		Name      string
		DeletedAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
		// DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`
	}
	// query: select * from users where deleted_at = 0 and id = 1
	// update: users set deleted_at = $currentUnixMillis where id = 1

	type User10Flag struct {
		ID        uint
		Name      string
		DeletedAt soft_delete.DeletedAt `gorm:"softDelete:flag"`
	}
	// query: select * from users where deleted_at = 0 and id = 1
	// update: users set deleted_at = 1 where id = 1

	type UserMixedFlag struct {
		ID        uint
		Name      string
		DeletedAt time.Time
		IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
		// IsDel     soft_delete.DeletedAt `gorm:"softDelete:milli,DeletedAtField:DeletedAt"`
		// IsDel     soft_delete.DeletedAt `gorm:"softDelete:nano,DeletedAtField:DeletedAt"`
	}
	// query: select * from users where is_del = 0
	// update: update users set is_del = 1, deleted_at = $currentUnixSecond where id = 1
}
