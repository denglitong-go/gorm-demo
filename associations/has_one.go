package associations

import "gorm.io/gorm"

// The only difference between hasOne and belongsTo is where the foreign key column is located.
// Let's say you have two entities: User and an Account.
// If the users table has the account_id column then a User belongsTo Account.
// 	(And the Account either hasOne or hasMany Users)
// But if the users table does not have the account_id column, and instead the accounts table
// has the user_id column, then User hasOne or hasMany Accounts.

// A `has one` association sets up a one-to-one connection with another model, but with
// somewhat different semantics and consequences. This association indicates that each
// instance of a model contains or processes one instance of another model.

// UserModel has one CreditCard, UserID is the foreign key
type UserModel struct {
	gorm.Model
	CreditCard CreditCard
}

// CreditCard belongs to UserModel, UserID is the foreign key
type CreditCard struct {
	gorm.Model
	Number   string
	UserID   uint
	UserName string
}

func GetAllUsers(db *gorm.DB) ([]UserModel, error) {
	var users []UserModel
	err := db.Model(&UserModel{}).Preload("CreditCard").Find(&users).Error
	return users, err
}

// UserModelOverrideForeignKey For a has one relationship, a foreign key field must also
// exist, the owner will save the primary key of the model belongs to it into this field.
type UserModelOverrideForeignKey struct {
	gorm.Model
	CreditCard CreditCard `gorm:"foreignKey:UserName"`
}

// UserModelOverrideReferences When save User model, save User.name to CreditCard.UserName
// as to set up the foreign key. By default, the owned entity will save the has one model's
// primary key into a foreign key.
type UserModelOverrideReferences struct {
	gorm.Model
	Name       string     `gorm:"index"`
	CreditCard CreditCard `gorm:"foreignKey:UserName;references:name"`
}

// Polymorphism associations

type Toy struct {
	gorm.Model
	Name      string
	OwnerID   uint
	OwnerType string
}

type Dog struct {
	gorm.Model
	Name string
	Toy  Toy `gorm:"polymorphism:Owner;"`
}

// db.Create(&Dog{Name:"dog1", Toy: Toy{Name:"toy1"}})
// INSERT INTO dogs (`name`) VALUES ('dog1')
// INSERT INTO toys (`name`,`owner_id`,`owner_type`) VALUES ('toy1',1,'dogs')

type Cat struct {
	gorm.Model
	Name string
	Toy  Toy `gorm:"polymorphism:Owner;polymorphicValue:cat_toy;"`
}

// db.Create(&Cat{Name:"cat1", Toy: Toy{Name:"toy2"}})
// INSERT INTO cats (`name`) VALUES ('cat1')
// INSERT INTO toys (`name`,`owner_id`,`owner_type`) VALUES ('toy2',1,'cat_toy')

type UserSelfReferentialHasOne struct {
	gorm.Model
	Name      string
	ManagerID *uint
	Manager   *UserSelfReferentialHasOne
}

type UserForeignKeyConstrains struct {
	gorm.Model
	CreditCard CreditCard `gorm:"constrains:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
