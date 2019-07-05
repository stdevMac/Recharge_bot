package parser

import (
	"bufio"
	"fmt"
	"github.com/stdevMac/Mybot/src/sendMail"
	"github.com/yanzay/tbot"
	"log"
	"os"
)

// ResponseParser
type ResponseParser struct {
	Number string
	Amount string
	Money string
}

// parse the message from the user and convert to ResponseParser
func parse(message *tbot.Message) ([]ResponseParser, error) {
	var response []ResponseParser

	response = append(response, ResponseParser{"53892073","1","20"})

	return response, nil
}

// GetBodyMessage parse the message and return string with the body of the mail
func GetBodyMessage(message *tbot.Message) (string, error) {
	var response string

	responseParsers, err := parse(message)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(responseParsers); i++ {
		response += fmt.Sprintf("%s,%s,%s\n", responseParsers[i].Number, responseParsers[i].Amount, responseParsers[i].Money)
	}

	return response, nil
}

// readLine returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func readLine(r *bufio.Reader) (string, error) {
	var (isPrefix bool = true
		err error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln),err
}

func GetUserPass(fileName string) sendMail.Sender {
	var sender sendMail.Sender

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewReader(file)

	s, e := readLine(scanner)
	if e != nil {
		log.Fatal(err)
	}
	sender.User = s
	s, e = readLine(scanner)
	if e != nil {
		log.Fatal(err)
	}
	sender.Password = s

	return sender
}

func GetToken(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewReader(file)

	s, e := readLine(scanner)
	if e != nil {
		log.Fatal(err)
	}

	return s
}