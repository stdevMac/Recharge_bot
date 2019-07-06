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
var bot *tbot.Server
var client *tbot.Client

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
	// Create new telegram bot server using token
	token := parser.GetToken("token.txt")
	bot = tbot.New(token)
	client = bot.Client()
}


func main() {
	// Use whitelist for Auth middleware, allow to interact only with user1 and user2
	//whitelist := []string{"marcosmaceo"}

	// Handle with StartHandler function
	bot.HandleMessage("/start", StartHandler)

	// Handle recharge
	bot.HandleMessage("/re", RechargeHandler)

	// Handle Resume
	bot.HandleMessage("/resume", ResumeHandler)

	// Start listening for messages
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is Running!")
}

func RechargeHandler(message *tbot.Message) {
	sender := parser.GetUserPass("user-pass.txt")

	bodyMessage, err := parser.GetBodyMessage(message)
	if err != nil {
		log.Fatal(err)
	}

	sender.SendMail([]string{"marcos.maceo@nauta.cu"}, message.From.Username,bodyMessage)

	client.SendMessage(message.Chat.ID,"Already made the recharge!")
	client.SendMessage(message.Chat.ID,fmt.Sprintf("Subject: \n=> %s  Body: \n=> %s ", message.From.Username, bodyMessage))

}

func StartHandler(message *tbot.Message) {
	// Handler can reply with several messages
	client.SendMessage(message.Chat.ID,"Hello, %s!")
	time.Sleep(1 * time.Second)
	client.SendMessage(message.Chat.ID,"We are ready to recharge some people!!")
}

func stat(h tbot.UpdateHandler) tbot.UpdateHandler {
	return func(u *tbot.Update) {
		start := time.Now()
		h(u)
		log.Printf("Handle time: %v", time.Now().Sub(start))
	}
}

func EchoHandler(message *tbot.Message) {
	client.SendMessage(message.Chat.ID, message.Text)
}

func ResumeHandler(message *tbot.Message) {
	resume, err := dbIntegration.GetResume(dbRedis, message.From.Username)
	if err != nil {
		client.SendMessage(message.Chat.ID,"Hubo error obteniendo el resumen del usuario" + message.From.Username)
		return
	}

	client.SendMessage(message.Chat.ID,resume)
}
