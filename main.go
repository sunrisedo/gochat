package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sunrise/tdexchat/chat"
)

var (
	PROGRAM     string
	VERSION     string
	BUILD       string
	COMMIT_SHA1 string

	gPort   int  // web server port
	gWsPort int  // websocket server port
	gIsHelp bool // show help info

	gEmotionNums [50]int

	//gRegEscape = regexp.MustCompile(`<script[\s\S]*?>[\s\S]*?</script>`)
)

func init() {
	fmt.Printf("-----------------------------------------------------\nProgram: %s\nVersion: %s\nBuild: %s\nCommit_sha1: %s\n-----------------------------------------------------\n", PROGRAM, VERSION, BUILD, COMMIT_SHA1)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	for i := 0; i < 50; i++ {
		gEmotionNums[i] = i
	}

	flag.IntVar(&gPort, "p", 8000, "web server port")
	flag.IntVar(&gWsPort, "wp", 8002, "websocket server port")
	flag.BoolVar(&gIsHelp, "h", false, "show help")
	flag.BoolVar(&gIsHelp, "help", false, "show help")
}

func main() {
	flag.Parse()

	if gIsHelp {
		flag.Usage()
		return
	}

	wsMux := http.NewServeMux()

	routerWeb()
	routerWebsocket(wsMux)
	fmt.Printf("rest server port:%d\nwebsocket server port:%d\n", gPort, gWsPort)
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", gPort), nil))
	}()
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", gWsPort), wsMux))
	}()

	select {}
}

func routerWeb() {
	http.HandleFunc("/", handleIndex)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.Handle("/upload/", http.StripPrefix("/upload", http.FileServer(http.Dir("upload"))))
}

func routerWebsocket(mux *http.ServeMux) {
	go chat.NewServer().Listen(mux)
}

var gFuncMap = template.FuncMap{
	"op": operate,
}

func operate(op string, a, b int) string {
	var result int

	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		result = a / b
	}

	return strconv.Itoa(result)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html").Delims("<{", "}>").Funcs(gFuncMap)
	t, err := t.ParseFiles("view/index.html")
	if err != nil {
		log.Println("handleIndex:", err)
		http.NotFound(w, r)
		return
	}
	t.Execute(w, map[string]interface{}{
		"emotionNums": gEmotionNums,
		"wsPort":      gWsPort,
	})
}

//@Deprecated
//func escapeBody(body []byte) []byte {
//    indexss := gRegEscape.FindAllIndex(body, -1)
//    if len(indexss) == 0 {
//        return body
//    }
//    var buffer bytes.Buffer
//    var i int
//    for _, indexs := range indexss {
//        fmt.Println(i, indexs[0], indexs[1])
//        buffer.Write(body[i:indexs[0]])
//        buffer.WriteString(html.EscapeString(string(body[indexs[0]:indexs[1]])))
//        i = indexs[1]
//    }
//    return buffer.Bytes()
//}
