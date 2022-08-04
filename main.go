package main

import (
	"cacheServer/cacheServer"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	MessageCache := cacheServer.NewCache("MessageCache", time.Second*5)
	go Start(":8080", MessageCache)
	go MockWebSocketMessage(MessageCache)
	select {}
}

func MockWebSocketMessage(cache cacheServer.CacheServer) {
	rand.Seed(time.Now().UnixNano())
	var i int = 0
	for {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		var receiver string
		if i%2 == 0 {
			receiver = "HYC"
		} else {
			receiver = "ZMN"
		}
		message := "send message " + strconv.Itoa(i) + " times"
		cache.Set(receiver, message, time.Now())
		i = i + 1
	}
}

func Start(Port string, cache cacheServer.CacheServer) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", health)
	mux.HandleFunc("/syncCache", syncCache(cache))
	mux.HandleFunc("/getCache", getCache(cache))
	svr := &http.Server{Addr: Port, Handler: mux}
	err := svr.ListenAndServe()
	return err
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "server work")
}

func getCache(cache cacheServer.CacheServer) func(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("getCache Success\n")
	return func(w http.ResponseWriter, r *http.Request) {
		key := GetUrlArg(r, "key")
		fmt.Fprintf(w, "Cache Key:%s,%v \n", key, cache.Get(key))
	}
}

func syncCache(cache cacheServer.CacheServer) func(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("syncCache Success\n")
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "All Cache:%v \n", cache.GetAll())
	}
}

func GetUrlArg(r *http.Request, name string) string {
	var arg string
	values := r.URL.Query()
	arg = values.Get(name)
	return arg
}
