package main 

import (
	"context"
    "fmt"
    "io"
    "log"
	"math/rand"
    "net/http"
    "os"
	"time"

	"cloud.google.com/go/pubsub"
)

var projectID string
var topicID string
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	projectID = os.Getenv("GOOGLE_PROJECT_ID")
	if projectID == "" {
		log.Fatal("Please provide a project ID")
	}

	topicID = os.Getenv("TOPIC")
	if topicID == "" {
		log.Fatal("Please provide a topic ID")
	}

	rand.Seed(time.Now().UnixNano())

	log.Println("project ID: ", projectID, " topic ID: ", topicID,
		" rand: ", randString(5))
}

func main() {
    http.HandleFunc("/", handler)
    log.Printf("Listening on port 8080")
    err := http.ListenAndServe(":8080", nil)
    log.Fatal(err)
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.Println("Serving request")

	if (r.Method != "GET") {
		fmt.Fprintln(w, "Sorry, I only understand GET requests")
		return
	}

	publish(w, projectID, topicID, "request-" + randString(5))
}    

func publish(w io.Writer, projectID, topicID, msg string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	t.PublishSettings.NumGoroutines = 1

	result := t.Publish(ctx, &pubsub.Message{Data: []byte(msg)})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
			return fmt.Errorf("Get: %w", err)
	}
	log.Println(msg + " published")
	
	fmt.Fprintf(w, "Published a message; msg ID: %v; msg: %s\n", id, msg)
	return nil
}

func randString(n int) string {
	b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}