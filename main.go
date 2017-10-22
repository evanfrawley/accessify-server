package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/evanfrawley/accessify-server/handlers"
)

const defaultAddr = ":80"
//const defaultAddr = ":443"

func main() {
    addr := os.Getenv("ADDR")
    if len(addr) == 0 {
        addr = defaultAddr
    }
    //
    //tlsCert := os.Getenv("TLSCERT")
    //tlsKey := os.Getenv("TLSKEY")

    mux := http.NewServeMux()
    mux.HandleFunc("/v1/accessify", handlers.GetAllData)

    fmt.Printf("server is listening at https://%s...\n", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
    //log.Fatal(http.ListenAndServeTLS(addr, tlsCert, tlsKey, mux))
}
