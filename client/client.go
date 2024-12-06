package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserInfo struct {
	ID        int        `gorm:"id" json:"-"`
	CreatedAt time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"updated_at" json:"updated_at"`
}

func (u *UserInfo) TableName() string {
	return "user_info"
}

func main() {

	User := "root"
	Passwd := "password"
	Host := "localhost"
	Port := 4000
	DbName := "tourism_festival"

	constr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", User, Passwd, Host, Port, DbName)

	client, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       constr,
		DefaultStringSize:         1 << 10,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}))
	if err != nil {
		panic(fmt.Errorf("failed to connect to mysql: %v", err))
	}

	var user UserInfo

	err = client.Select("id,created_at,updated_at").First(&user).Error
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("%+v\n", user)
}
