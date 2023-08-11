package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type PubSubNotification struct {
	Message struct {
		Attributes   map[string]interface{} `json:"attributes"`
		MessageID    string                 `json:"messageId"`
		PublishTime  time.Time              `json:"publishTime"`
		// Data is a base64 encoded
		Data string `json:"data"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func (notif PubSubNotification) decodeMsg() (string, error) {
	decoded, err1 := base64.StdEncoding.DecodeString(notif.Message.Data)
	if err1 != nil {
		return "", fmt.Errorf("decoding notification message base64 encoded data %q: %v", notif.Message.Data, err1)
	}
	return string(decoded), nil
}

func main() {
	http.HandleFunc("/", process)
	log.Printf("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
func process(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving request")
	if r.Method == "GET" {
		fmt.Fprintln(w, "I am ready to process POST requests from pub/sub subscription")
		return
	}

	log.Println("Reading POST data")
	body, err1 := ioutil.ReadAll(r.Body)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Error reading POST data: %v\n", err1)
		return
	}

	log.Println("POST data =", string(body))
	var notification PubSubNotification
	r.Body.Close()

	if len(bytes.TrimSpace(body)) == 0 {
		log.Println("Empty request body. Expecting a PubSub message.")
		return
	}

	err2 := json.Unmarshal(body, &notification)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling POST data: %v. Could not parse %q\n", err2, string(body))
		return
	}

	decoded, err1 := base64.StdEncoding.DecodeString(notification.Message.Data)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "decoding base64 encoded data %q: %v\n", notification.Message.Data, err1)
		return
	}

	log.Println("received: ", string(decoded))
}
