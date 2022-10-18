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
            //body, err := io.ReadAll(r.Body)
            //if err != nil {
            //    log.Println(err)
            //    return
            //}
            //r2 := r.Clone(r.Context())
            //
            //r.Body = io.NopCloser(bytes.NewReader(body))
            //r2.Body = io.NopCloser(bytes.NewReader(body))
            //
            //r.ParseForm()
            //
            //rec := httptest.NewRecorder()
            //
            //log.Printf("proxying request to /knative-sources/slack-direct")
            //r.Host = "github-adapter.knative-sources.127.0.0.1.sslip.io"
            //r.URL.Path = "/default/slack-direct"
            //
            //p.ServeHTTP(rec, r2)

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
