package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/stdevMac/Mybot/src/dbIntegration"
	"github.com/stdevMac/Mybot/src/parser"
	"github.com/yanzay/tbot"
	"log"
	"time"
)

var dbRedis redis.Conn

func init()  {

	pool := dbIntegration.NewPool()
	// get a connection from the pool (redis.Conn)
	dbRedis = pool.Get()
	// use defer to close the connection when the function completes
	defer dbRedis.Close()

	// call Redis PING command to test connectivity
	err := dbIntegration.Ping(dbRedis)
	if err != nil {
		fmt.Println(err)
	}

}

func main() {

	// Create new telegram bot server using token
	token := parser.GetToken("token.txt")
	bot, err := tbot.NewServer(token)
	if err != nil {
		log.Fatal(err)
	}

	// Use whitelist for Auth middleware, allow to interact only with user1 and user2
	whitelist := []string{"marcosmaceo"}
	bot.AddMiddleware(tbot.NewAuth(whitelist))

	// Handle with StartHandler function
	bot.HandleFunc("/start", StartHandler)

	// Handle recharge
	bot.HandleFunc("/re", RechargeHandler)

	// Handle Resume
	bot.HandleFunc("/resume", ResumeHandler)

	// Set default handler if you want to process unmatched input
	bot.HandleDefault(EchoHandler)


	fmt.Println("Server is Running!")
	// Start listening for messages
	err = bot.ListenAndServe()
	log.Fatal(err)
}

func RechargeHandler(message *tbot.Message) {
	sender := parser.GetUserPass("user-pass.txt")

	bodyMessage, err := parser.GetBodyMessage(message)
	if err != nil {
		log.Fatal(err)
	}

	sender.SendMail([]string{"marcos.maceo@nauta.cu"}, message.From.UserName,bodyMessage)

	message.Reply("Already made the recharge!")
	message.Reply(fmt.Sprintf("Subject: \n=> %s  Body: \n=> %s ", message.From.UserName, bodyMessage))

}

func StartHandler(message *tbot.Message) {
	// Handler can reply with several messages
	message.Replyf("Hello, %s!", message.From)
	time.Sleep(1 * time.Second)
	message.Reply("We are ready to recharge some people!!")
}

func EchoHandler(message *tbot.Message) {
	message.Reply(message.Text())
}

func ResumeHandler(message *tbot.Message) {
	resume, err := dbIntegration.GetResume(dbRedis, message.From.UserName)
	if err != nil {
		message.Reply("Hubo error obteniendo el resumen del usuario" + message.From.UserName)
	}

	message.Reply(resume)
}
