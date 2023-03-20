package common

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

// User model
// Models are normal structs with basic Go types, pointers/alias of them
// or custom types implementing Scanner and Valuer interfaces.
type User struct {
	// GORM prefers convention to configuration.
	// By default, GORM uses ID as primary key,
	// pluralizes struct name to snake_cases as column name,
	// and uses CreatedAt, UpdatedAt to track creating/updating time.
	ID           uint
	Name         string
	Email        *string
	Age          uint
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Model GORM defined a gorm.Model struct, which includes fields
// ID, CreatedAt, UpdatedAt, DeletedAt
type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// UserModel equals to User
type UserModel struct {
	gorm.Model
	Name        string
	Email       *string
	Age         uint
	Birthday    *time.Time
	ActivatedAt time.Time
}

type Scanner interface {
	// Scan assigns a value from a database drier.
	//
	// The src value will be of one of the following types:
	// 	int64
	//  float64
	//  bool
	//  string
	// 	time.Time
	//  nil - for NULL values
	//
	// An error should be returned if the value cannot be stored
	// without loss of information.
	//
	// Reference types such as []type are only valid
	// until the next call to Scan and should not be retained.
	// Their underlying memory is owned by the driver.
	// If retention is necessary, copy their values before the next call to Scan.
	Scan(src any) error
}

type Valuer interface {
	// Value returns a driver Value.
	// Value must not panic.
	//
	// Types implementing Valuer interface are able to convert themselves to a driver Value.
	Value() (Value, error)
}

// Value is a value that drivers must be able to handle.
type Value any
