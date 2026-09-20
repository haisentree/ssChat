package internal

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

var addr = flag.String("addr", ":8080", "http service address")
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

	// 启动提示：打印运行地址和端口
	displayAddr := *addr
	if strings.HasPrefix(displayAddr, ":") {
		displayAddr = "localhost" + displayAddr
	}
	log.Printf("ssChat 服务已启动，访问地址: http://%s （按 Ctrl+C 停止）", displayAddr)

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