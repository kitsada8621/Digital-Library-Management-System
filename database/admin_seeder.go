package database

import (
	"dlms/models"
	"dlms/services"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func EnsureSeederAdminData() error {

	userService := services.UserService()
	roleService := services.RoleService()

	admin, err := roleService.FindByName("admin")
	if err != nil {
		fmt.Println("Err roleAdmin: ", err.Error())
		return err
	}

	author, err := roleService.FindByName("author")
	if err != nil {
		fmt.Println("Err roleAuthor: ", err.Error())
		return err
	}

	member, err := roleService.FindByName("user")
	if err != nil {
		fmt.Println("Err roleUser: ", err.Error())
		return err
	}

	defaultAccount := []models.User{
		models.User{
			ID:        primitive.NewObjectID(),
			Roles:     []primitive.ObjectID{admin.ID},
			Name:      "admin admin",
			Phone:     "0999999999",
			Email:     "admin@exmaple.com",
			Username:  "admin",
			Password:  "Dev123456!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		models.User{
			ID:        primitive.NewObjectID(),
			Roles:     []primitive.ObjectID{author.ID},
			Name:      "author author",
			Phone:     "0999999999",
			Email:     "author@exmaple.com",
			Username:  "author",
			Password:  "Dev123456!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		models.User{
			ID:        primitive.NewObjectID(),
			Roles:     []primitive.ObjectID{member.ID},
			Name:      "user user",
			Phone:     "0999999999",
			Email:     "user@exmaple.com",
			Username:  "user",
			Password:  "Dev123456!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	users := []models.User{}
	for _, account := range defaultAccount {
		if _, err := userService.FindByUsername(account.Username); err != nil {
			passwordHashed, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
			if err != nil {
				continue
			}
			account.Password = string(passwordHashed)
			users = append(users, account)
		}
	}

	if len(users) > 0 {
		if err := userService.CreateMany(users); err != nil {
			return err
		}
	}

	return nil

}
