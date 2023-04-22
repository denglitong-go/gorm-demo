package associations

import (
	"gorm.io/gorm"
	"time"
)

// User2 has and belongs to many languages, `user_languages` is the join table
// When using GORM AutoMigrate to create a table for User, GORM will create the join table automatically.
type User2 struct {
	gorm.Model
	Name      string
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	gorm.Model
	Name string
}

// Join table: user_languages
// 	foreign key: user_id, reference: users.id
// 	foreign key: language_id, reference: languages.id

// Back-Reference

type UserBackReference struct {
	gorm.Model
	Languages []*LanguagesBackReference `gorm:"many2many:user_languages;"`
}

type LanguagesBackReference struct {
	gorm.Model
	Name  string
	Users []*UserBackReference `gorm:"many2many:user_languages;"`
}

// Retrieve

func RetrieveUsersWithLanguagesEagerLoading(db *gorm.DB) ([]User2, error) {
	var users []User2
	err := db.Model(&User2{}).Preload("Languages").Find(&users).Error
	return users, err
}

func RetrieveLanguagesWithUsersEagerLoading(db *gorm.DB) ([]Language, error) {
	var languages []Language
	err := db.Model(&Language{}).Preload("Users").Find(&languages).Error
	return languages, err
}

type User2OverrideForeignKey struct {
	gorm.Model
	Languages []LanguageOverrideForeignKey `gorm:"many2many:user_languages;foreignKey:Refer;joinForeignKey:UserReferID;references:Refer;joinReferences:ProfileReferID"`
	Refer     uint                         `gorm:"index:;unique"`
}

type LanguageOverrideForeignKey struct {
	gorm.Model
	Name  string
	Refer uint `gorm:"index:;unique"`
}

// Which creates join table: user_languages
// 	foreign key: user_refer_id, reference: user.refer
// 	foreign key: language_refer_id, reference: language.refer

type UserSelfReferentialMany2Many struct {
	gorm.Model
	Friends []*User `gorm:"many2many:user_friends;"`
}

// which creates join table: user_friends
// 	foreign key: user_id, reference: users.id
//  foreign key: friend_id, reference: users.id

// Customize Join Table

type Person struct {
	gorm.Model
	Name      string
	Addresses []Address `gorm:"many2many:person_addresses"`
}

type Address struct {
	gorm.Model
	Name string
}

type PersonAddress struct {
	PersonID  uint `gorm:"primaryKey"`
	AddressID uint `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func RetrieveJoinTables(db *gorm.DB) ([]PersonAddress, error) {
	var personAddresses []PersonAddress
	err := db.SetupJoinTable(&Person{}, "Addresses", &PersonAddress{})
	if err != nil {
		return personAddresses, err
	}
	err = db.Find(&personAddresses).Error
	return personAddresses, err
}
