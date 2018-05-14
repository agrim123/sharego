package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("method:", r.Method)
		if r.Method == "GET" {
			crutime := time.Now().Unix()
			h := md5.New()
			io.WriteString(h, strconv.FormatInt(crutime, 10))
			token := fmt.Sprintf("%x", h.Sum(nil))

			t, _ := template.ParseFiles("upload.html")
			t.Execute(w, token)
		} else {
			r.ParseMultipartForm(32 << 20)
			file, handler, err := r.FormFile("file")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				return
			}
			defer f.Close()
			io.Copy(f, file)
		}
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		allfiles, _ := ioutil.ReadDir("./uploads")
		type File struct {
			Name string
			Path string
		}

		var Files []*File

		for _, f := range allfiles {
			file := &File{f.Name(), "./uploads/" + f.Name()}
			Files = append(Files, file)
		}

		t, _ := template.ParseFiles("list.html")
		t.Execute(w, Files)

	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
