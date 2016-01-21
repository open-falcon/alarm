package cron

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/open-falcon/alarm/g"
	"github.com/open-falcon/common/model"
	"log"
	"time"
)

func ReadHighEvent() {
	queues := g.Config().Redis.HighQueues
	if len(queues) == 0 {
		return
	}

	for {
		event, err := popEvent(queues)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		consume(event, true)
	}
}

func ReadLowEvent() {
	queues := g.Config().Redis.LowQueues
	if len(queues) == 0 {
		return
	}

	for {
		event, err := popEvent(queues)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		consume(event, false)
	}
}

func popEvent(queues []string) (*model.Event, error) {

	count := len(queues)

	params := make([]interface{}, count+1)
	for i := 0; i < count; i++ {
		params[i] = queues[i]
	}
	// set timeout 0
	params[count] = 0

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	reply, err := redis.Strings(rc.Do("BRPOP", params...))
	if err != nil {
		log.Printf("get alarm event from redis fail: %v", err)
		return nil, err
	}

	if g.Config().Debug {
		log.Println("======>>>>redis brpoop:",params)
		log.Println( reply )
	}

	var event model.Event
	err = json.Unmarshal([]byte(reply[1]), &event)
	if err != nil {
		log.Printf("parse alarm event fail: %v", err)
		return nil, err
	}

	if g.Config().Debug {
		log.Println("======>>>>")
		log.Println(event.String())
	}

	// save in memory. display in dashboard
	g.Events.Put(&event)

	return &event, nil
}
