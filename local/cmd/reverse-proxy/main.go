package main

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "net/http"
)

const baseURL = "http://github-adapter.knative-sources.127.0.0.1.sslip.io"

func forward(r *http.Request, body []byte, path string) {
    go func() {
        // create a new url from the raw RequestURI sent by the client
        log.Printf("proxying request to %s\n", path)

        url := fmt.Sprintf("%s/%s", baseURL, path)

        proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
        proxyReq.Header = r.Header

        resp, err := http.DefaultClient.Do(proxyReq)
        if err != nil {
            log.Printf("error: %v", err)
            return
        }
        defer resp.Body.Close()
    }()
}

func main() {

    handler := func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        r.Body.Close()

        forward(r, body, "default/slack-direct")
        forward(r, body, "default/slack-persisted")
        forward(r, body, "default/broker-demo21")
        forward(r, body, "default/broker-demo22")
    }

    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
}
