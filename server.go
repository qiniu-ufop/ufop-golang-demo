package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

// HTTPGetMaxSize 最大处理的文件长度
const HTTPGetMaxSize = 2 * 1024 * 1024

func httpGet(url string) (body []byte, err error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("httpGet error: %s", err)
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(http.MaxBytesReader(nil, res.Body, HTTPGetMaxSize))
	if err != nil {
		return nil, fmt.Errorf("httpGet read body error: %s", err)
	}
	return
}

func handler(rw http.ResponseWriter, req *http.Request) {
	var err error
	defer func() {
		if err != nil {
			http.Error(rw, err.Error(), 500)
		}
	}()

	defer req.Body.Close()

	var body []byte

	url := req.URL.Query().Get("url")
	if url != "" {
		body, err = httpGet(url)
		if err != nil {
			log.Println("handler http get error:", err.Error())
		}
	} else {
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("handler body read error: %s\n", err.Error())
			return
		}
	}

	tpl, err := template.New("res").Parse(`
req body:
-------------
{{.body}}

time: {{.time}}
		`)
	if err != nil {
		panic("")
	}
	brw := bufio.NewWriter(rw)
	length := len(body)
	tpl.Execute(brw, map[string]interface{}{
		"body":   string(body),
		"length": length,
		"time":   time.Now(),
	})
	brw.Flush()
}

func health(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("ok"))
}

func main() {
	port := os.Getenv("PORT_HTTP")
	if port == "" {
		port = "9100"
	}
	http.HandleFunc("/handler", handler)
	http.HandleFunc("/health", health)
	log.Fatalln(http.ListenAndServe("0.0.0.0:"+port, nil))
}
