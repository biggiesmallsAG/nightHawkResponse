/*
 *@package  main
 *@file     redis_client.go
 *@author   Daniel Eden <daniel.eden@gmail.com>
 *@
 *@description  nighHawk Redis client PUB/SUB.
 */

package nightHawk

import (
 	"time"
 	"encoding/json"
 	"gopkg.in/redis.v4"
 	"fmt"
)

type WebSocketMsg struct {
    Time string
    Level string
    Message string
}

func RedisPublish(level string, message string, redis_cli bool) {
	if redis_cli {
    	msg := WebSocketMsg{
    		Time: time.Now().UTC().Format(Layout),
    		Level: level,
    		Message: message,
    	}

    	redis_net := fmt.Sprintf("%s:%d", REDIS_SERVER, REDIS_PORT)

    	client := redis.NewClient(&redis.Options{
    			Addr: redis_net,
    			Password: "",
    			DB: 0,
    		})

    	if pubsub, err := client.Subscribe(REDIS_CHAN); err == nil {
    		defer pubsub.Close()
	    	if data, err := json.Marshal(&msg); err == nil {
	    		client.Publish(REDIS_CHAN, string(data)).Result()
	    	}    	    		
    	}
	}
}