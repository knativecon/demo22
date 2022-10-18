package main

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "net/http"
)

const baseURL = "http://github-adapter.knative-sources.127.0.0.1.sslip.io"

func forward(r *http.Request, body []byte, path string) error {
    // create a new url from the raw RequestURI sent by the client
    log.Printf("proxying request to %s\n", path)

    url := fmt.Sprintf("%s/%s", baseURL, path)

    proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
    proxyReq.Header = r.Header

    resp, err := http.DefaultClient.Do(proxyReq)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

func main() {

    handler := func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        //
        //// you can reassign the body if you need to parse it as multipart
        //r.Body = io.NopCloser(bytes.NewReader(body))

        err = forward(r, body, "default/slack-persisted")
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadGateway)
            return
        }

        err = forward(r, body, "default/broker")
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadGateway)
            return
        }

        err = forward(r, body, "default/slack-direct")
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadGateway)
            return
        }
    }

    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
}
