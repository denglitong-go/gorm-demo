package associations

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto Create/Update
// GORM will auto-save associations and its reference using Upset when creating/updating a record.

func CreateUser2(db *gorm.DB) {
	user := User2{
		Name: "someone",
		Languages: []Language{
			{Name: "ZH"},
			{Name: "EN"},
		},
	}

	db.Create(&user)
	// BEGIN TRANSACTION;
	// INSERT INTO "languages" (name) values ("ZH"),("EN");
	// INSERT INTO users (name) values ("someone");
	// INSERT INTO user_languages (user_id, language_id) values (111, 1), (111, 2);

	// If you want to update associations data, you can use the `FullSaveAssociations` mode.
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)
	// insert into languages (name) values ("ZH"),("EN") ON DUPLICATE KEY SET name=VALUES(name)

	// Skip Auto Create/Update
	// Use Select or Omit.

	// insert into users (name) values ("someone")
	db.Select("Name").Create(&user)

	// skip create languages when creating a user
	db.Omit("Languages").Create(&user)

	// skip all associations when creating a user
	db.Omit(clause.Associations).Create(&user)
}

// Select/Omit Association fields.

// Setup Association Mode
// var user User
// db.Model(&user).Association("Languages")
// `user` is the source model, it must contains primary key
// `Languages` is a relationship's field name
// If the above two requirements matched, the AssociationMode should be started
// db.Model(&user).Association("Languages").Error

// Find Associations
// db.Model(&user).Association("Languages").Find(&languages)
// codes := []string{"zh-CN","en-US"};
// db.Model(&user).Where("code IN ?", codes).Association("Languages").Find(&languages)

// Append Associations
// db.Model(&user).Association("Languages").Append([]Language{xx,xx})

// Replace Associations
// db.Model(&user).Association("Languages").Replace([]Language{xx,xx})

// Delete Associations
// db.Model(&user).Association("Languages").Delete([]Language{xx,xx})

// Clear Associations
// db.Model(&user).Association("Languages").Clear()

// Count Associations
// db.Model(&user).Association("Languages").Count()
// db.Model(&user).Where("code IN ?", codes).Association("Languages").Count()

// Batch Data
// users := []User{user1, user2, user3}

// Find all roles for all users
// db.Model(&users).Association("Role").Find(&roles)

// Delete User A from all user's team
// db.Model(&users).Association("Team").Delete(&userA)

// For Append, Replace with batch data, the length of the arguments needs to be equal to the data's length or else
// it will return an error.

// Append userA to user1's team, append userB to user2's team,
// append {userA,userB,userC} to user3's team.
// db.Model(&users).Association("Team").Append(&userA, &userB, &[]User{userA,userB,userC})

// Delete with Select

// delete user's account when deleting user
// db.select("Account").Delete(&user)
// db.select("Account").Delete(&users) // batch delete

// delete user's Orders, CreditCards relations when deleting user
// db.select("Orders", "CreditCards").Delete(&user)

// delete user's has one/many/many2many relations when deleting user
// db.Select(clause.Associations).Delete(&user)

// Associations will only be deleted if the deleting records' primary key is not zero,
// GORM will use those primary keys as conditions to delete selected associations.

// DOESN'T WORK
// will delete all user with name = 'someone', but those users' account won't be deleted.
// db.Select("Account").Where("name = ?", "someone").Delete(&User{})

// will delete user with name = 'someone' and id = 1, and user1's account will be deleted.
// db.Select("Account").Where("name = ?", "someone").Delete(&User{ID: 1})

// will delete user with id = 1, and user1's account will be deleted.
// db.Select("Account").Delete(&User{ID: 1})

// Association Tags
// foreignKey
// references
// polymorphic
// polymorphicValue
// many2many
// joinForeignKey
// joinReferences
// constraint
