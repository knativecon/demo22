# Local Development

To receive GitHub events and forward them to your GitHub adapter running in Kind, 
do the following:
* in one terminal, run `go run ./cmd/reverse-proxy/main.go`. This reverse proxy forward requests
  on port 8080 to all GitHub sources running in the Kind cluster.
* In another terminal, run `ssh -R 80:localhost:8080 localhost.run`. This creates a tunnel to port 8080
  * Look for the last console line and copy the URL. (e.g. 8ff5e4119dff9d.lhr.life tunneled with tls termination, `https://8ff5e4119dff9d.lhr.life`) 
  * See [localhost.run](https://localhost.run) for more information
* In your GitHub repository, go to `settings` -> `webhooks` and add a new webhook using the copied
  URL. Change `Content type` to `application/json`. Leave the secret empty, and select individual events. 
  Select `Issues` and `Issue comment`. Uncheck `push`.

You are all set. 

