package associations

import "gorm.io/gorm"

// A has many association sets up a one-to-many connection with another model,
// unlike has one, the owner could have zero or many instances of models.

// UserM has many CreditCards, UserID is the foreign key
type UserM struct {
	gorm.Model
	CreditCards []CreditCard
}

type CreditCardM struct {
	gorm.Model
	Number     string
	UserID     uint
	UserRefer  uint
	UserNumber string
}

// GetAll Retrieve user list with eager loading credit cards
func GetAll(db *gorm.DB) ([]UserM, error) {
	var users []UserM
	err := db.Model(&UserM{}).Preload("CreditCards").Find(&users).Error
	return users, err
}

type UserMOverrideForeignKey struct {
	gorm.Model
	// CreditCardM.UserRefer = User.ID
	CreditCards []CreditCardM `gorm:"foreignKey:UserRefer"`
}

type UserMOverrideReferences struct {
	gorm.Model
	MemberNumber string
	// CreditCard.UserNumber = User.MemberNumber
	CreditCards []CreditCardM `gorm:"foreignKey:UserNumber;references=MemberNumber"`
}

// GORM supports polymorphism association for `has one` and `has many`, it will save
// owned entity's table name into polymorphic type's field, primary key value into the
// polymorphic field.

type ToyM struct {
	gorm.Model
	OwnerID   uint
	OwnerType string
}

type DogM struct {
	gorm.Model
	Toys []ToyM `gorm:"polymorphic:Owner;"`
}

// db.Create(&DogM{Name:"dog1",Toys:[]ToyM{{Name:"toy1"},{Name:"toy2"}}})
// INSERT INTO dogs (`name`) values ('dog1')
// INSERT INTO toys (`name`,`owner_id`,`owner_type`) values ('toy1', 1, 'dogs'), ('toy2', 1, 'dogs')

type UserMSelfReferentialHasMany struct {
	gorm.Model
	Name      string
	ManagerID *uint
	Teammates []UserMSelfReferentialHasMany `gorm:"foreignKey:ManagerID"`
}
