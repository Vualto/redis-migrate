package main

import (
	"flag"
	"fmt"
	"os"

	redis "github.com/gomodule/redigo/redis"
)

func main() {
	sourceURL := flag.String("source", "", "source redis connection URL")
	destURL := flag.String("destination", "", "destination redis connection URL")
	count := flag.Int("count", 100, "redis scan 'count'")

	flag.Parse()

	source, err := redis.DialURL(*sourceURL)
	if err != nil {
		panic("error connecting to source redis")
	}

	dest, err := redis.DialURL(*destURL)
	if err != nil {
		panic("error connecting to destination redis")
	}

	size, err := redis.Int(source.Do("DBSIZE"))
	if err != nil {
		panic(err.Error())
	}

	var cursor int
	var progress int
	var reply interface{}

	for {
		reply, err = source.Do("SCAN", cursor, "COUNT", *count)
		if err != nil {
			panic(err.Error())
		}

		for _, k := range reply.([]interface{})[1].([]interface{}) {
			progress++

			key, err := redis.String(k, nil)
			if err != nil {
				panic(err.Error())
			}
			val, err := source.Do("GET", key)
			if err != nil {
				panic(err.Error())
			}

			fmt.Printf("processing '%s', progress %d\n", key, progress)

			_, err = dest.Do("SET", key, val)
			if err != nil {
				panic(err.Error())
			}
		}

		cursor, err = redis.Int(reply.([]interface{})[0], nil)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("cursor: %d\n", cursor)

		if cursor == 0 {
			break
		}
	}

	destSize, err := redis.Int(dest.Do("DBSIZE"))
	if err != nil {
		panic(err.Error())
	}

	if destSize >= size {
		fmt.Println("completed.")
		os.Exit(0)
	}

	fmt.Println("something went wrong")
	os.Exit(1)
}
