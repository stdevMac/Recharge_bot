package main

import (
	"fmt"
	"github.com/yanzay/tbot"
	"log"
	"time"
)

func main() {
	token := "598068140:AAFqLPtPxSgggUTDFCRmrdZ-T7E2IRY9edY"
	// Create new telegram bot server using token
	bot, err := tbot.NewServer(token)
	if err != nil {
		log.Fatal(err)
	}

	// Use whitelist for Auth middleware, allow to interact only with user1 and user2
	whitelist := []string{"marcosmaceo"}
	bot.AddMiddleware(tbot.NewAuth(whitelist))

	// Yo handler works without slash, simple text response
	bot.Handle("yo", "YO!")

	// Handle with HiHandler function
	bot.HandleFunc("/hi", HiHandler)

	// Handle recharge
	bot.HandleFunc("/re", RechargeHandler)

	// Set default handler if you want to process unmatched input
	bot.HandleDefault(EchoHandler)


	fmt.Println("Server is Running!")
	// Start listening for messages
	err = bot.ListenAndServe()
	log.Fatal(err)
}

func RechargeHandler(message *tbot.Message) {

}

func HiHandler(message *tbot.Message) {
	// Handler can reply with several messages
	message.Replyf("Hello, %s!", message.From)
	time.Sleep(1 * time.Second)
	message.Reply("We are ready to recharge some people!!")
}

func EchoHandler(message *tbot.Message) {
	message.Reply(message.Text())
}
