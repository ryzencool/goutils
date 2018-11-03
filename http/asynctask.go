package http

import (
	"math/rand"
	"time"
)

type task struct {
	funcs func(url string, params map[string]interface{}) interface{}
	url string
	params map[string]interface{}
}

func timeopt() int {
	time.Sleep(2 * time.Second)
	return rand.Int()
}

func asyncTask(done <-chan interface{}, ts ...task) []interface{} {
	ch := make(chan interface{})
	responses := make([]interface{}, 0)
	for _, t := range ts {
		go func(t task) {
			res := t.funcs(t.url, t.params)
			ch <- res
		}(t)
	}

loop:
	for {
		select {
		case <-done:
			return nil
		case r := <-ch:
			responses = append(responses, r)
			if len(responses) == len(ts) {
				break loop
			}
		case <-time.After(5 * time.Second):
			return nil
		}
	}
	return responses
}
