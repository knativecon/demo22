package main

import (
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
            //log.Printf("proxying request to /knative-sources/slack-direct")
            //r.Host = "github-adapter.knative-sources.127.0.0.1.sslip.io"
            //r.URL.Path = "/default/slack-direct"
            //
            //p.ServeHTTP(w, r)

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
