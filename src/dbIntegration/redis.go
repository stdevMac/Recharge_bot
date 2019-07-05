package dbIntegration


import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/stdevMac/Mybot/src/parser"
	"time"
)

func NewPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
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
		money = append(amount, resp.Amount)
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

	// SET object
	_, err = c.Do("SET", user, jsonUsers)
	if err != nil {
		return err
	}

	return nil
}

func GetInfoUsers(c redis.Conn, username string) (InfoUser, error) {

	s, err := redis.String(c.Do("GET", username))
	if err == redis.ErrNil {
		fmt.Println("Requested user does not have records")
	} else if err != nil {
		return InfoUser{}, err
	}

	usr := InfoUser{}
	err = json.Unmarshal([]byte(s), &usr)

	fmt.Printf("%+v\n", usr)

	return usr, err

}