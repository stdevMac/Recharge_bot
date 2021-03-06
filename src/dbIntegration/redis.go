package dbIntegration

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/stdevMac/Recharge_bot/src/parser"
	"strconv"
	"time"
)

var atackers = "Atackers"

func NewPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// Ping tests connectivity for dbIntegration (PONG should be returned)
func Ping(c redis.Conn) error {
	// Send PING command to Redis
	// PING command returns a Redis "Simple String"
	// Use dbIntegration.String to convert the interface type to string
	s, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}

// InfoUser is a struct for keep user recharge info
type InfoUser struct {
	Username  string `json:"username"`
	Numbers     []string `json:"number"`
	Dates     []string `json:"dates"`
	Amount     []string `json:"amount"`
	Money     []string `json:"money"`
}

// SetRechargeInfo Update in redis information of recharger
func SetRechargeInfo(c redis.Conn, user string, responseParser []parser.ResponseParser) error {

	var numbers []string
	var dates []string
	var amount []string
	var money []string

	for _, resp := range responseParser {
		numbers = append(numbers, resp.Number)
		dates = append(dates, time.Now().Format(time.RFC850))
		amount = append(amount, resp.Amount)
		money = append(money, resp.Money)
	}

	usr := InfoUser{
		Username:  	user,
		Numbers:    numbers,
		Dates: 		dates,
		Amount:		amount,
		Money:		money,
	}

	// Update previous record of the user
	infoUsers, err := GetInfoUsers(c, user)
	if err == redis.ErrNil {
		infoUsers.Username = user
	} else if err != nil {
		fmt.Println("Jone!!")
		return err
	}

	infoUsers.Numbers = append(infoUsers.Numbers, usr.Numbers...)
	infoUsers.Dates = append(infoUsers.Dates, usr.Dates...)
	infoUsers.Amount = append(infoUsers.Amount, usr.Amount...)
	infoUsers.Money = append(infoUsers.Money, usr.Money...)

	// serialize InfoUser object to JSON
	jsonUsers, err := json.Marshal(infoUsers)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", infoUsers)

	// SET object
	_, err = c.Do("SET", user, jsonUsers)
	if err != nil {
		return err
	}

	return nil
}

// SetBasic Update in redis first information
func SetBasic(c redis.Conn, user string) error {

	basic := InfoUser{
		Username:  	user,
		Numbers:    []string{},
		Dates: 		[]string{},
		Amount:		[]string{},
		Money:		[]string{},
	}

	// serialize InfoUser object to JSON
	jsonUsers, err := json.Marshal(basic)
	if err != nil {
		return err
	}

	// SET object
	_, err = c.Do("SET", user, jsonUsers)
	if err != nil {
		return err
	}

	return nil
}

type Attackers struct {
	badUsers []string `json:"badUsers"`
	numberAttacks []int `json:"numberAttacks"`
} 

// SetAttacker Update in redis first information
func SetAttacker(c redis.Conn, user string) error {

	atk, err := GetAttackers(c)
	if err != nil {
		fmt.Println(err)
		return err
	}
	hasRegisterAttack := false
	for i, badUser := range atk.badUsers {
		if badUser == user {
			hasRegisterAttack = true
			atk.numberAttacks[i] += 1	
		}
	}
	if !hasRegisterAttack {
		atk.numberAttacks = append(atk.numberAttacks, 1)
		atk.badUsers = append(atk.badUsers, user)
	}

	// serialize InfoUser object to JSON
	jsonUsers, err := json.Marshal(atk)
	if err != nil {
		return err
	}

	// SET object
	_, err = c.Do("SET", atackers, jsonUsers)
	if err != nil {
		return err
	}

	return nil
}

func GetInfoUsers(c redis.Conn, username string) (InfoUser, error) {

	s, err := redis.String(c.Do("GET", username))
	if err == redis.ErrNil {
		fmt.Println("Requested user does not have records")
		return InfoUser{},err
	} else if err != nil {
		return InfoUser{}, err
	}
	usr := InfoUser{}
	err = json.Unmarshal([]byte(s), &usr)

	return usr, err

}

func GetAttackers(c redis.Conn) (Attackers, error) {

	s, err := redis.String(c.Do("GET", atackers))
	if err == redis.ErrNil {
		fmt.Println("Requested user does not have records")
	} else if err != nil {
		return Attackers{}, err
	}

	atk := Attackers{}
	err = json.Unmarshal([]byte(s), &atk)

	fmt.Printf("%+v\n", atk)

	return atk, err

}

func GetResume(c redis.Conn, username string) (string, error) {

	s, err := redis.String(c.Do("GET", username))
	if err == redis.ErrNil {
		fmt.Println("Requested user does not have records")
	} else if err != nil {
		return "", err
	}

	usr := InfoUser{}
	err = json.Unmarshal([]byte(s), &usr)

	fmt.Printf("%+v\n", usr)

	return prettyFormat(usr), nil

}
func prettyFormat(infoUser InfoUser) string {

	initial := "Usted ha realizado "

	var response string
	var numberOfRecharges = 0

	for i := 0; i < len(infoUser.Numbers); i++ {
		tmp, err := strconv.Atoi(infoUser.Amount[i])
		if err != nil {
			fmt.Println(
				fmt.Sprintf("Error converting the amount of recharges for user by the user %s => ",
					infoUser.Username) + err.Error())
			continue
		}
		numberOfRecharges += tmp
		response += fmt.Sprintf("# %s -> %s recargas -> %s -> con fecha %s. \n",
			infoUser.Numbers[i], infoUser.Amount[i], infoUser.Money[i], infoUser.Dates[i])
	}

	return initial + strconv.Itoa(numberOfRecharges) + " recargas.\n\n" + response
}