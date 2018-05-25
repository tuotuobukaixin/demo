package main



import (
"fmt"

"github.com/garyburd/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "49.4.10.247:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()


	username, err := redis.String(c.Do("GET", "mykey"))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get mykey: %v \n", username)
	}
}