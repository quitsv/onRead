package Controllers

import (
	"context"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	redis "github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var context_redis = context.Background()

func setRedis(rdb *redis.Client, key string, value string, expiration int) {
	err := rdb.Set(context_redis, key, value, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func getRedis(rdb *redis.Client, key string) string {
	val, err := rdb.Get(context_redis, key).Result()

	if err != nil {
		log.Fatal(err)
	}
	return val
}

func RunTools() {
	db := Connect()
	defer db.Close()

	//goRedis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	setRedis(rdb, "Best Seller", "ayo Baca Buku yang sedang HOT!", 0)

	gocron := gocron.NewScheduler(time.UTC)

	gocron.Every(1).Month().Do(func() {
		query := "SELECT email, nama FROM pengguna where tipe = 0"

		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		var email string
		var nama string

		for rows.Next() {
			err = rows.Scan(&email, &nama)
			if err != nil {
				log.Fatal(err)
			}
			message := "Hey " + nama + "! " + getRedis(rdb, "Best Seller")
			go sendEmail(email, []byte(message))
		}
	})
	gocron.StartBlocking()
}

func sendEmail(emailPenerima string, message []byte) {
	// Configuration
	sender := getSender()
	to := []string{emailPenerima}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Confirmation
	println("Sending email to: " + emailPenerima)

	// Create authentication
	auth := smtp.PlainAuth("", sender.Email, sender.Passwd, smtpHost)

	// Send actual message
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender.Email, to, message)
	if err != nil {
		log.Fatal(err)
	}
	println("email sent")
}

func getSender() User {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal(envErr)
	}

	var sender User

	sender.Email = os.Getenv("EMAIL")
	sender.Passwd = os.Getenv("PASSWORD")

	return sender
}
