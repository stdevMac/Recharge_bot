package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/stdevMac/Recharge_bot/src/dbIntegration"
	"github.com/stdevMac/Recharge_bot/src/parser"
	"github.com/yanzay/tbot"
	"log"
	"time"
)

var dbRedis redis.Conn
var bot *tbot.Server
var client *tbot.Client
var whitelist []string

func init() {

	// Create new telegram bot server using token
	token := parser.GetFileFirstLine("token.txt")
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

	//dbIntegration.SetBasic(dbRedis, "marcosmaceo")

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
		client.SendMessage(message.Chat.ID, "Hubo error con la conexion a la base de datos"+message.From.Username)
		log.Printf("Couldn't access to database: %s", err)
		return
	}
	fmt.Println(message.Chat.ID)
	sender := parser.GetUserPass("user-pass.txt")

	response, err := parser.GetBodyMessage(message)
	if err != nil {
		client.SendMessage(message.Chat.ID, "Hubo error con interpretando el mensaje, revise que cumple con los parametros...")
		log.Printf("Couldn't get body message: %s", err)
		return
	}

	err = dbIntegration.SetRechargeInfo(dbRedis, message.From.Username, response)
	if err != nil {
		client.SendMessage(message.Chat.ID, "Hubo error actualizando sus datos...contacte a @marcosmaceo")
		log.Printf("Couldn't update user info: %s", err)
		log.Print(err)
		return
	}

	bodyMessage := parser.PrettyPrint(response)

	sendTo := parser.GetFileFirstLine("send_to.txt")

	err = sender.SendMail(bodyMessage, sendTo)
	if err != nil {
		_, err = client.SendMessage(message.Chat.ID, "Su recarga no se pudo procesar...contacta a @marcosmaceo")
		log.Printf("Couldn't send email for recharge, send error: %s", err)
		return
	}

	client.SendMessage("677517973", fmt.Sprintf("El usuario \"@%s\" realizo el siguiente pedido: \n %s", message.Chat.Username, message.Text))
	client.SendMessage(message.Chat.ID, "Su recarga esta siendo procesada...")
	log.Printf("Recharge from https://t.me/%s is processing....")
}

func StartHandler(message *tbot.Message) {
	// Handler can reply with several messages
	client.SendMessage(message.Chat.ID, fmt.Sprintf("Hola, %s %s!", message.Chat.FirstName, message.Chat.LastName))
	time.Sleep(1 * time.Second)
	client.SendMessage(message.Chat.ID, fmt.Sprintf("%s %s(%s)\nMediante este bot se pueden realizar recargas....Unicamente si te encuentras entre los usuarios que pueden.....El proceder es el siguiente:\n	Para recargar:\n```\n/recharge #1, #2, #3\n```\n\nDonde son lo siguiente:\n#1 -> Numero a recargar\n#2 -> Cantidad de recargas a dicho numero\n#3 -> Cantidad de dinero a transferir\n\nAdemas pueden ser multiples:\n```\n/recharge #1, #2, #3\n#4, #5, #6\n#7, #8, #9\n#10, #11, #12\n```\n...\nDonde #4, #7 y #10 son los numeros a recargar, #5, #8 y #11 son la cantidad de recargas a dichos numeros y #6, #9 y #12 son las cantidades con las que se recargaran\n\nSi se quiere ver el resumen de recargas que ha realizado, solo poner el siguiente comando\n`/resume`\n\nPara cualquier duda escribir a @marcosmaceo. Gracias", message.Chat.FirstName, message.Chat.LastName, message.Chat.Username))
	log.Printf("Finish start for https://t.me/%s", message.Chat.Username)
}

func stat(h tbot.UpdateHandler) tbot.UpdateHandler {
	return func(u *tbot.Update) {
		start := time.Now()
		for _, word := range whitelist {
			if word == u.Message.From.Username {
				h(u)
				log.Printf("Handle time: %v", time.Now().Sub(start))
				return
			}
		}
		err := dbIntegration.SetAttacker(dbRedis, u.Message.From.Username)
		if err != nil {
			log.Print("Couldn't set attacker")
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
		client.SendMessage(message.Chat.ID, "Se produjo un error con la conexion a la base de datos"+message.From.Username)
		log.Printf("Couldn't access to database: %s", err)
		return
	}
	var username string
	if len(message.Text) > len("/resume") {
		runes := []rune(message.Text)
		username = string(runes[len("/resume"):len(message.Text)])
	} else {
		username = message.Chat.Username
	}

	resume, err := dbIntegration.GetResume(dbRedis, username)
	if err != nil {
		client.SendMessage(message.Chat.ID, "Se produjo un error obteniendo el resumen del usuario "+message.From.Username)
		log.Printf("Couldn't get resumen for %s", message.Chat.Username)
		return
	}

	client.SendMessage(message.Chat.ID, resume)
	log.Printf("Resumen requested for %s", message.Chat.Username)
}

func AddRechargerHandler(message *tbot.Message) {

}
