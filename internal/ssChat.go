package internal

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8081", "http service address")
var wss = NewWSServer()

func NewSSChatServer() *SSChatServer {
	return &SSChatServer{}
}

type SSChatServer struct {
}

func (s SSChatServer) Run() {
	flag.Parse()
	files := http.FileServer(http.Dir("../internal/web"))

	wss.ListenBroadcast()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", wss.wsHandle)
	http.Handle("/web/", http.StripPrefix("/web/", files))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "../internal/web/home.html")
}

//package internal
//
//import (
//	"flag"
//	"log"
//	"net/http"
//)
//
//func NewSSChatServer() SSChatServer {
//	return SSChatServer{}
//}
//
//type SSChatServer struct {
//}
//
//var addr = flag.String("addr", ":8081", "http service address")
//
//func (s SSChatServer) Run() {
//	flag.Parse()
//	hub := newHub()
//	go hub.run()
//	http.HandleFunc("/", serveHome)
//	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
//		serveWs(hub, w, r)
//	})
//	err := http.ListenAndServe(*addr, nil)
//	if err != nil {
//		log.Fatal("ListenAndServe: ", err)
//	}
//}
//
//func serveHome(w http.ResponseWriter, r *http.Request) {
//	log.Println(r.URL)
//	if r.URL.Path != "/" {
//		http.Error(w, "Not found", http.StatusNotFound)
//		return
//	}
//	if r.Method != http.MethodGet {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//	http.ServeFile(w, r, "../internal/home.html")
//}
