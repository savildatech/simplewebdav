package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/webdav"
)

func main() {
	dir := os.Getenv("WEBDAV_DIR")
	portStr := os.Getenv("WEBDAV_PORT")
	user := os.Getenv("WEBDAV_USER")
	pass := os.Getenv("WEBDAV_PASS")
	readonlyStr := os.Getenv("WEBDAV_READONLY")

	if dir == "" || portStr == "" || user == "" || pass == "" {
		log.Fatal("Required env vars: WEBDAV_DIR, WEBDAV_PORT, WEBDAV_USER, WEBDAV_PASS")
	}

	readonly := strings.ToLower(readonlyStr) == "true"

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

		if readonly {
			switch r.Method {
			case "GET", "HEAD", "PROPFIND", "OPTIONS":
				// allow
			default:
				http.Error(w, "Read-only mode", http.StatusForbidden)
				return
			}
		}

		handler.ServeHTTP(w, r)
	})

	addr := ":" + portStr
	log.Printf("Serving %s on http://localhost%s (readonly: %v)", dir, addr, readonly)
	log.Fatal(http.ListenAndServe(addr, nil))
}
