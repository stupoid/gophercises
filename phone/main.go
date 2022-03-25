package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PhoneNumber struct {
	gorm.Model
	Number string `gorm:"unique"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("phone.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&PhoneNumber{})

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	numbers := []PhoneNumber{}
	for scanner.Scan() {
		numbers = append(numbers, PhoneNumber{Number: normalize(scanner.Text())})
	}
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "number"}},
		DoUpdates: clause.AssignmentColumns([]string{"number"}),
	}).CreateInBatches(numbers, 100)
}

func normalize(number string) string {
	var sb strings.Builder
	for _, r := range number {
		if unicode.IsDigit(r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
