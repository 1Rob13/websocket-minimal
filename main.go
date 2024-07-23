package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-ping/ping"
	"nhooyr.io/websocket"
)

type PingProbe struct {
	Name  string
	Stats ping.Statistics
}

func main() {

	var (
		PingStats = PingProbe{Name: "Test", Stats: ping.Statistics{}}
	)

	go func() {
		for {
			PingStats.Stats = Ping()
			time.Sleep(time.Second * 15)
		}
	}()

	http.HandleFunc("/ping-stats", func(w http.ResponseWriter, r *http.Request) {

		jsonBytes, err := json.Marshal(PingStats.Stats)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		reader := bytes.NewReader(jsonBytes)
		b, err := io.ReadAll(reader)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	})

	http.HandleFunc("/ping-ws", func(w http.ResponseWriter, r *http.Request) {

		opts := websocket.AcceptOptions{InsecureSkipVerify: true}

		conn, err := websocket.Accept(w, r, &opts)

		if err != nil {
			fmt.Errorf("cant connect websocket %v", err)
			return
		}

		go func() {
			for {

				b, err := json.Marshal(PingStats)

				if err != nil {
					fmt.Errorf("cant send to websocket %v", err)
					return
				}
				conn.Write(context.Background(), 1, b)
				time.Sleep(time.Second * 16)
			}
		}()

	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
func Ping() ping.Statistics {

	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		panic(err)
	}

	pinger.Count = 1

	pinger.Run()
	fmt.Println(pinger.Statistics())

	return *pinger.Statistics()

}
