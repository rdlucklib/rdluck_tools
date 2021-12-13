package redis

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

const (
	redis_conn  = "127.0.0.1:6379" //redis链接地址
	remove_keys = "key_*"          //要删除的key
)

//删除Redis Keys
func RemoveAllKeys() {
	conn, err := redis.Dial("tcp", redis_conn)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	vals, err := redis.Strings(conn.Do("KEYS", remove_keys))
	conn.Send("MULTI")
	fmt.Println(len(vals))
	for k, v := range vals {
		fmt.Println(k, v)
		conn.Send("DEL", vals[k])
	}
	fmt.Println(conn.Do("EXEC"))
}
