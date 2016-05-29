package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	//"path/filepath"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
)

const usage = `Serves static files across multiple folders and mountpoints 
Usage: servefiles -http=[ADDRESS] [MOUNTPOINT] [FOLDER] [MOUNTPOINT] [FOLDER]...
    
Parameters:
    -http           the host and port to serve on. Default 127.0.0.1:44444
    [MOUNTPOINT]    mountpoint in the webserver address
    [FOLDER]        which folder to mount on the specified mountpoint
    
Example:
    servefiles -addr="127.0.0.1:9898" "/web/" "~/web" "/media/" "~/media"
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}

	var addr string
	flag.StringVar(&addr, "http", "127.0.0.1:44444", "the host and port to serve on")
	flag.Parse()

	args := flag.Args()
	if len(args)%2 != 0 || len(args) == 0 {
		log.Fatalf("error: missing mountpoint or folder in arguments:\n%v\n", args)
	}

	ns := vfs.NewNameSpace()

	for j := 0; j < len(args); j = j + 2 {
		mountpoint := args[j]
		mountpath := args[j+1]

		// this checks the real file exists
		if _, err := os.Stat(args[j+1]); err != nil {
			log.Fatalf("error: could not stat path %s: %s\n", args[j+1], err)
		}
		log.Printf("Binding %s to %s\n", mountpath, mountpoint)
		ns.Bind(mountpoint, vfs.OS(mountpath), "/", vfs.BindAfter)
	}
	log.Printf("Serving on: %s\n", addr)
	log.Println(http.ListenAndServe(addr, logger(http.FileServer(httpfs.New(ns)))))
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s\t%s\t%s", r.RemoteAddr, r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
