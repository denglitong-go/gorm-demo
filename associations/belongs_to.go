package associations

import "gorm.io/gorm"

// A belongs to association sets up a one-to-one connection with another model, such that
// each instance of the declaring model "belongs to" one instance of the other model.

type Company struct {
	// ID int
	gorm.Model
	CompanyRefer int
	Code         string
	Name         string
}

// User belongs to Company, CompanyID is the foreign key.
// There is both a CompanyID and a Company. By default, the CompanyID is implicitly used to
// create a foreign key relationship between the User and Company tables, and thus must be
// included in the User struct in order to fill the Company inner struct.
type User struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   Company
}

type UserOverrideForeignKey struct {
	gorm.Model
	Name         string
	CompanyRefer int
	Company      Company `gorm:"foreignKey:CompanyRefer"` // use CompanyRefer as foreign key
}

type UserOverrideReferences struct {
	gorm.Model
	Name       string
	CompanyKey string
	Company    Company `gorm:"references:Code"` // assign Company.Code to User.CompanyKey
}

// UserSpecifyReference GORM usually guess the relationship as `has one` if override foreign key name
// already exists in owner's type, we need to specify references in the belongs to relationship.
type UserSpecifyReference struct {
	gorm.Model
	Name      string
	CompanyID string
	// Company already exists CompanyID field, we need to specify the references field
	// to keep association to be `belongs to`, otherwise, GORM will guess it as `has one`
	Company Company `gorm:"references:CompanyID"`
}

type UserForeignKeyConstraints struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   Company `gorm:"constraints:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
