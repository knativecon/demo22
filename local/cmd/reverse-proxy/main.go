package main

import (
    "bytes"
    "io"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    remote, err := url.Parse("http://github-adapter.knative-sources.127.0.0.1.sslip.io")
    if err != nil {
        panic(err)
    }

    handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
        return func(w http.ResponseWriter, r *http.Request) {
            // we need to buffer the body if we want to read it here and send it
            // in the request.
            body, err := io.ReadAll(r.Body)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            // you can reassign the body if you need to parse it as multipart
            r.Body = io.NopCloser(bytes.NewReader(body))

            // create a new url from the raw RequestURI sent by the client
            url := "http://github-adapter.knative-sources.127.0.0.1.sslip.io/default/slack-direct"

            proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))

            // We may want to filter some headers, otherwise we could just use a shallow copy
            // proxyReq.Header = req.Header
            proxyReq.Header = make(http.Header)
            for h, val := range r.Header {
                proxyReq.Header[h] = val
            }

            resp, err := http.DefaultClient.Do(proxyReq)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadGateway)
                return
            }
            defer resp.Body.Close()

            log.Printf("proxying request to /knative-sources/slack-persisted")
            r.Host = "github-adapter.knative-sources.127.0.0.1.sslip.io"
            r.URL.Path = "/default/slack-persisted"

            p.ServeHTTP(w, r)
        }
    }

    proxy := httputil.NewSingleHostReverseProxy(remote)
    http.HandleFunc("/", handler(proxy))
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
}
