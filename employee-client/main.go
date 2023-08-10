package main 

import (
    "context"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"

    "cloud.google.com/go/compute/metadata"
    "google.golang.org/api/idtoken"
)

func main() {
    employeeApiUrl := os.Getenv("EMPLOYEE_API")
    if employeeApiUrl == ""  {
        log.Fatal("Employ API URL is required")
    }

    http.HandleFunc("/", handler)
    log.Printf("Listening on port 8080")
    err := http.ListenAndServe(":8080", nil)
    log.Fatal(err)
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.Println("Serving request")

    employeeApiUrl := os.Getenv("EMPLOYEE_API")
    
    if metadata.OnGCE()  {
        log.Println("OnGCE: true")
        // employeeApiUrl is the audience
        client, err := idtoken.NewClient(context.Background(), employeeApiUrl)
        if err != nil {
            log.Fatal("Failed to get idtoken")
        }
        resp, err := client.Get(employeeApiUrl + "/api/help")
        if err != nil {
            log.Fatal("Failed calling Employee API")
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Fatal("Failed read response from Employee API")
        }

        fmt.Println("Response from Employee API: ", string(body))
    
        w.WriteHeader(http.StatusOK)
        w.Write(body)
    } else {
        resp, err := http.Get(employeeApiUrl + "/api/help")
        if err != nil {
            log.Fatal("Failed calling Employee API")
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Fatal("Failed read response from Employee API")
        }

        fmt.Println("Response from Employee API: ", string(body))
    
        w.WriteHeader(http.StatusOK)
        w.Write(body)
    }
}    