package parser

import (
	"bufio"
	"fmt"
	"github.com/stdevMac/Recharge_bot/src/sendMail"
	"github.com/yanzay/tbot"
	"log"
	"os"
	"strings"
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

	lines := strings.Split(message.Text[len("\recharge "):],"\n")

	for _, line := range lines {
		//line = line[:len(line) - 1]
		split := strings.Split(line, ",")
		if err := correctSplit(split) ; err != nil {
			return []ResponseParser{}, err
		}
		response = append(response, ResponseParser{split[0], split[1], split[2]})
	}

	return response, nil
}
func correctSplit(split []string) error {
	if len(split) != 3 {
		return fmt.Errorf("Mas de 3 argumentos")
	}
	return nil
}

// GetBodyMessage parse the message and return string with the body of the mail
func GetBodyMessage(message *tbot.Message) ([]ResponseParser, error) {

	responseParsers, err := parse(message)
	if err != nil {
		return []ResponseParser{}, err
	}
	return responseParsers, nil
}

func PrettyPrint(responseParsers []ResponseParser) string {
	var response string

	for i := 0; i < len(responseParsers); i++ {
		response += fmt.Sprintf("%s,%s,%s\n", responseParsers[i].Number, responseParsers[i].Amount, responseParsers[i].Money)
	}

	return response
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
		log.Print(err)
		return sendMail.Sender{}
	}
	defer file.Close()

	scanner := bufio.NewReader(file)

	s, e := readLine(scanner)
	if e != nil {
		log.Print(err)
		return sendMail.Sender{}
	}
	sender.User = s
	s, e = readLine(scanner)
	if e != nil {
		log.Print(err)
		return sendMail.Sender{}
	}
	sender.Password = s

	return sender
}

func GetFileFirstLine(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Print(err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewReader(file)

	s, e := readLine(scanner)
	if e != nil {
		log.Print(err)
		return ""
	}

	return s
}