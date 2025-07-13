package main

import (
	"log"
	"net/http"
	"os"
//	"strconv"

	"golang.org/x/net/webdav"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatal("Usage: program dir port username password")
	}

	dir := os.Args[1]
	portStr := os.Args[2]
	//port, err := strconv.Atoi(portStr)
	//if err != nil {
	//	log.Fatal("Invalid port")
	//}
	user := os.Args[3]
	pass := os.Args[4]

	handler := &webdav.Handler{
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != user || p != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="webdav"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	})

	addr := ":" + portStr
	log.Printf("Serving %s on http://localhost%s", dir, addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
