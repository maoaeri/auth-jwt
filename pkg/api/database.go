package api

import (
	"fmt"
	"log"
	"myapp/pkg/helper"
	"myapp/pkg/model"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
)

func GetDB() *gorm.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var (
		DB_USER     = helper.GetEnvVar("DB_USER")
		DB_PASSWORD = helper.GetEnvVar("DB_PASSWORD")
		DB_NAME     = helper.GetEnvVar("DB_NAME")
	)

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	sqldb, err := connection.DB()
	if err != nil {
		log.Fatalln(err)
	}

	if err = sqldb.Ping(); err != nil {
		log.Fatalln(err)
	}

	q := `CREATE TABLE IF NOT EXISTS users (
		ID SERIAL PRIMARY KEY,
		Created_at TIMESTAMP,
		Updated_at TIMESTAMP,
		Deleted_at TIMESTAMP,
		Name VARCHAR(255),
		Email VARCHAR(255),
		Password VARCHAR(255),
		Role VARCHAR(255),
		Refresh_token VARCHAR(255)
	);`
	connection.Exec(q)
	fmt.Println("Connected to database")
	return connection
}

func CloseDB(connection *gorm.DB) {
	sqldb, err := connection.DB()
	if err != nil {
		log.Fatalln(err)
	}
	sqldb.Close()
}

func CreateUser(user *model.User) error {
	connection := GetDB()
	defer CloseDB(connection)
	result := connection.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAllUsers() (users []model.User) {
	connection := GetDB()
	defer CloseDB(connection)

	result := connection.Find(&users)
	if result.Error != nil {
		log.Fatalln("Error in fetching users")
		return nil
	}
	return users
}

func GetUser(email string) (user model.User) {
	connection := GetDB()
	defer CloseDB(connection)

	result := connection.Where("email = ?", email).First(&user)
	if result.Error != nil {
		log.Fatalln("Error in fetching user")
	}
	return user
}

func DeleteUser(email string) {
	connection := GetDB()
	defer CloseDB(connection)

	var user model.User
	user = GetUser(email)

	result := connection.Delete(&user)
	if result.Error != nil {
		log.Fatalln("Error in deleting user.")
	}
}

func DeleteAllUsers() {
	connection := GetDB()
	defer CloseDB(connection)

	var users []model.User
	users = GetAllUsers()

	result := connection.Delete(&users)
	if result.Error != nil {
		log.Fatalln("Error in deleting user.")
	}
}

/*func CheckDuplicateEmail(email string) (bool, error) {
	connection := GetDB()
	defer CloseDB(connection)
	err := connection.Where("email = ?", email).First(&dbuser).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		err := errors.New("Email already in use")
		return true, err
	}
	return false, nil
}*/
