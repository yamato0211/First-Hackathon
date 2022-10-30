package db

import (
	"first-hackathon/graph/model"
	"first-hackathon/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Psql *gorm.DB
)

func InitDB() (err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", utils.DbHost, utils.DbUser, utils.DbPass, utils.DbName, utils.DbPort)
	Psql, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	if err = Psql.AutoMigrate(&model.User{}); err != nil {
		return
	}

	return
}
