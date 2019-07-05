package parser

import (
	"bufio"
	"github.com/stdevMac/Mybot/src/sendMail"
	"github.com/yanzay/tbot"
	"log"
	"os"
)

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

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
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

	s, e := Readln(scanner)
	if e != nil {
		log.Fatal(err)
	}
	sender.User = s
	s, e = Readln(scanner)
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

	s, e := Readln(scanner)
	if e != nil {
		log.Fatal(err)
	}

	return s
}