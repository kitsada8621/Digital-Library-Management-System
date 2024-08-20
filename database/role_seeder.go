package database

import (
	"dlms/models"
	"dlms/services"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func EnsureSeederRoleData() error {
	roleService := services.RoleService()
	rolesMaster := []models.Role{
		models.Role{
			ID:        primitive.NewObjectID(),
			RoleName:  "admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		models.Role{
			ID:        primitive.NewObjectID(),
			RoleName:  "author",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		models.Role{
			ID:        primitive.NewObjectID(),
			RoleName:  "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	roles := []models.Role{}
	for _, row := range rolesMaster {
		if _, err := roleService.FindByName(row.RoleName); err != nil {
			roles = append(roles, row)
		}
	}

	fmt.Printf("total seeder role: %d\n", len(roles))
	if len(roles) > 0 {
		if err := roleService.CreateRoleMany(roles); err != nil {
			fmt.Println("\n\nErr migrate role: ", err.Error())
			return err
		}
	}
	return nil
}
