
package main

import (
        "fmt"
        "log"
        "net/http"
        "os"

        "github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

func incrementHandler(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving request")

        conn := redisPool.Get()
        defer conn.Close()

        log.Println("Obtained connection")
        counter, err := redis.Int(conn.Do("INCR", "visits"))
        if err != nil {
                http.Error(w, "Error incrementing visitor counter", http.StatusInternalServerError)
                return
        }
        fmt.Fprintf(w, "Visitor number: %d\n", counter)
}

func main() {
        redisHost := os.Getenv("REDISHOST")
        redisPort := os.Getenv("REDISPORT")
        redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

		log.Println("redisAddr: ", redisAddr)

		const maxConnections = 10
        redisPool = &redis.Pool{
                MaxIdle: maxConnections,
                Dial:    func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
        }

        http.HandleFunc("/", incrementHandler)

        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
        }
        log.Printf("Listening on port %s", port)
        if err := http.ListenAndServe(":"+port, nil); err != nil {
                log.Fatal(err)
        }
}