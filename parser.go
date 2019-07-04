package main

import "github.com/yanzay/tbot"

type ResponseParser struct {
	number string
	amount int
	money int
}


func Parse(message *tbot.Message) []ResponseParser {
	var response []ResponseParser

	response = append(response, ResponseParser{"53892073",1,20})

	return response

}