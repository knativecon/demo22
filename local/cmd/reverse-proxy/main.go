package main

import (
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
            r.Host = "github-adapter.knative-sources.127.0.0.1.sslip.io"
            r.URL.Path = "/knative-sources/knativecon-issues"

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
