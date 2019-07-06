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
var whitelist []string

func init()  {

	// Create new telegram bot server using token
	token := parser.GetToken("token.txt")
	bot = tbot.New(token)
	client = bot.Client()

	whitelist = []string{"marcosmaceo"}
}


func main() {
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

	dbIntegration.SetBasic(dbRedis, "marcosmaceo")

	// Handle with StartHandler function
	bot.HandleMessage("/start", StartHandler)

	// Handle recharge
	bot.HandleMessage("/recharge", RechargeHandler)

	// Handle Resume
	bot.HandleMessage("/resume", ResumeHandler)

	bot.HandleMessage("", EchoHandler)

	// Setup Middleware
	bot.Use(stat)

	// Start listening for messages
	fmt.Println("Server is Running!")
	err = bot.Start()
	if err != nil {
		log.Print(err)
	}
}

func RechargeHandler(message *tbot.Message) {
	err := dbIntegration.Ping(dbRedis)
	if err != nil {
		client.SendMessage(message.Chat.ID,"Hubo error con la conexion a la base de datos" + message.From.Username)
		return
	}
	sender := parser.GetUserPass("user-pass.txt")

	response, err := parser.GetBodyMessage(message)
	if err != nil {
		log.Print(err)
		return
	}

	err = dbIntegration.SetRechargeInfo(dbRedis, message.From.Username, response)
	if err != nil {
		log.Print(err)
		return
	}

	bodyMessage := parser.PrettyPrint(response)

	sender.SendMail([]string{"marcos.maceo@nauta.cu"}, message.From.Username, bodyMessage)

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
		for _, word := range whitelist {
			if word == u.Message.From.Username {
				h(u)
				log.Printf("Handle time: %v", time.Now().Sub(start))
			}
		}

		log.Printf("Handle time: %v", time.Now().Sub(start))
		log.Printf("User not allowed at %v", time.Now().Sub(start))
	}
}

func EchoHandler(message *tbot.Message) {
	client.SendMessage(message.Chat.ID, message.Text)
}

func ResumeHandler(message *tbot.Message) {
	err := dbIntegration.Ping(dbRedis)
	if err != nil {
		client.SendMessage(message.Chat.ID,"Hubo error con la conexion a la base de datos" + message.From.Username)
		return
	}
	resume, err := dbIntegration.GetResume(dbRedis, message.From.Username)
	if err != nil {
		client.SendMessage(message.Chat.ID,"Hubo error obteniendo el resumen del usuario" + message.From.Username)
		return
	}

	client.SendMessage(message.Chat.ID,resume)
}
