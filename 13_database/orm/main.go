package main

import (
	"bytes"
	"fmt"
	"orm/env"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// type Animal struct {
// 	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
// 	Name      string    `gorm:"type:varchar(20);not null;"`
// 	Emoji     string    `gorm:"type:varchar(10);not null"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

type Teste struct {
	gorm.Model
}

func returnUrl(url string) string {
	urlBytes := []byte(url)
	b := make([]byte, len(urlBytes))
	i := 0
	for _, v := range urlBytes {
		if v != 0 {
			b[i] = v
			i++
		}
	}

	return string(b[:bytes.IndexByte(b, 0)])
}

func setup() *gorm.DB {
	envs, err := env.SourceEnv()
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=America/Sao_Paulo",
		envs["HOST"], envs["USER"], envs["PASSWORD"], envs["DATABASE"], envs["PORT"])
	//dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Brazil/Sao_Paulo"
	// db, err := gorm.Open(postgres.New(postgres.Config{
	// 	DriverName: "pgx",
	// 	DSN:        returnUrl(dsn),
	// }), &gorm.Config{})
	db, err := gorm.Open(postgres.Open(returnUrl(dsn)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Migrate the schema
	err = db.AutoMigrate(&Teste{})
	if err != nil {
		panic(err)
	}
	return db
}

func Insert(db *gorm.DB, i int) {
	for c := 0; c < i; c++ {
		ob := Teste{}
		err := db.Create(&ob).Error
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	db := setup()

	start := time.Now()
	Insert(db, 1)
	fmt.Println(time.Since(start))
	// Insert(db, Animal{Name: "Dog", Emoji: "\U0001F431"})
	// db.Select(&Animal{}, "Id")
	// var a Animal
	// db.Find(a)
}
