package database

import (
	"fmt"
	"log"

	pb "github.com/SirishaGopigiri/sample-grpc-server/user"
	"gorm.io/gorm"
)

type User struct {
	Name  string `gorm:"primaryKey"`
	Email string
	Age   int
}

func CreateUser(db *gorm.DB, pb_user *pb.User) error {
	user := &User{
		Name:  pb_user.Name,
		Email: pb_user.Email,
		Age:   int(pb_user.Age),
	}
	result := db.Create(user)
	if result.Error != nil {
		log.Fatalf("Error inserting data: %v", result.Error)
		return result.Error
	} else {
		fmt.Println("Data inserted successfully! Name:", user.Name)
	}
	return nil
}

func GetUser(db *gorm.DB) ([]*pb.User, error) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		log.Printf("Unable to get results, %v", result.Error)
		return nil, result.Error
	}
	pb_users := []*pb.User{}
	for _, user := range users {
		pb_user := &pb.User{
			Name:  user.Name,
			Email: user.Email,
			Age:   int32(user.Age),
		}
		pb_users = append(pb_users, pb_user)
	}
	return pb_users, nil
}
