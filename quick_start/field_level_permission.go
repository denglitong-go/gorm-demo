package quick_start

import (
	"time"
)

// UserModelWithPermission allows you to change the field-level permission with tag,
// you can make a field to be read-only, write-only, create-only, update-only or ignored.
type UserModelWithPermission struct {
	ID           uint   `gorm:"<-:create"`          // allow read and create
	Name         string `gorm:"<-:update"`          // allow read and update
	Email        string `gorm:"<-"`                 // allow read and write (create and update)
	Age          uint   `gorm:"<-:false"`           // allow read, disable write permission
	Birthday     string `gorm:"->"`                 // ready-only (disable write permission unless it configured)
	MemberNumber string `gorm:"->;<-:create"`       // allow read and create
	ActivatedAt  uint   `gorm:"->;false;<-:create"` // create only (disable read from db)
	CreatedAt    uint   `gorm:"-"`                  // ignore this field when write and read with struct
	// Ignored fields won't be created when using GORM Migrator to create table.
	UpdatedAt uint `gorm:"-:all"`       // ignore this field when write, read and migrate with struct
	DeletedAt uint `gorm:"-:migration"` // ignore this field when migrate with struct
}

// UserModelTrackingTimeUnix GORM use CreatedAt, UpdatedAt to track creating/updating time by convention,
// and GORM will set the current time (time.Now().Local()) when creating/updating if the fields are defined.
// To use fields with a different name, you can configure those fields with tag `autoCreateTime`, `autoUpdateTime`.
// If you prefer to save UNIX (milli/nano) seconds instead of time, you can simply change the field's
// data type from time.Time to int.
type UserModelTrackingTimeUnix struct {
	ID        uint
	CreatedAt time.Time // GORM will set to current time if it's zero on creating
	UpdatedAt int       // GORM will set to current time unix seconds on updating or if it's zero on creating
	DeletedAt int64     `gorm:"autoUpdateTime:nano"`  // GORM will set to unix nanoseconds as deleting time
	Created   int64     `gorm:"autoUpdateTime:milli"` // GORM will set to unix milliseconds as creating time
	Updated   int64     `gorm:"autoCreateTime"`       // // GORM will set to unix seconds as updating time
}

type Author struct {
	Name  string
	Email string
}

type BlogEmbeddedModel struct {
	ID     uint
	Author Author `gorm:"embedded"` // For a normal struct field, you can embed it with the tag.
	Upvote uint
}

// BlogEmbeddedModelFull equals to BlogEmbeddedModel
type BlogEmbeddedModelFull struct {
	ID     uint
	Name   string
	Email  string
	Upvote uint
}

type BlogEmbeddedPrefixModel struct {
	ID     uint
	Author Author `gorm:"embedded;embeddedPrefix:author_"`
	Upvote uint
}

// BlogEmbeddedPrefixModelFull equals to BlogEmbeddedPrefixModel
type BlogEmbeddedPrefixModelFull struct {
	ID          uint
	AuthorName  string
	AuthorEmail string
	Upvote      uint
}

// UserModelWithFieldTags Tags are optional to use when declaring models.
// Tags are case-insensitive, however, camlCase is preferred.
type UserModelWithFieldTags struct {
	ID           uint   `gorm:"column:id;primaryKey;autoIncrement"`               // column db name; column as primaryKey; column auto increment
	Name         string `gorm:"type:string;unique;not null;uniqueIndex:idx_name"` // column data type; column as unique; column as NOT NULL
	Email        string `gorm:"serializer:json"`                                  // specifies serializer for how to serialize and deserialize data into db
	Age          uint   `gorm:"default:0;check:age<200"`                          // column default value; column checks
	Birthday     time.Time
	Weight       float64   `gorm:"scale:2;precision:2"`                     // column scale; column precision 整数部分 2 位，小数部分 2 位
	MemberNumber string    `gorm:"size:256"`                                // specifies column data size/length
	CreatedAt    time.Time `gorm:"index:idx_created_at"`                    // column as index, use same index name for two fields will create composite indexes
	UpdatedAt    time.Time `gorm:"<-"`                                      // field's write permission, <-:create create-only, <-:update update-only, <-:false no write permission, <- create and update permission
	DeletedAt    time.Time `gorm:"->"`                                      // field's read permission, ->:false no read permission
	ArchivedAt   time.Time `gorm:"-"`                                       // ignore this field, - ignore read/write; -:migration ignore migration; -:all ignore read/write/migrate
	ActivatedAt  time.Time `gorm:"comment:activate when accept invitation"` // add comment for field when migration
}
