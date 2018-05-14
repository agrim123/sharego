package Handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.ServeFile(rw, r, "../templates/index.html")
}

func UploadHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		fp := path.Join("../templates", "upload.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(rw, token); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		f, err := os.OpenFile("../uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		http.Redirect(rw, r, "/", http.StatusSeeOther)
	}
}

func ListHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	allfiles, _ := ioutil.ReadDir("../uploads")

	type File struct {
		Name string
		Path string
	}

	var Files []*File

	for _, f := range allfiles {
		file := &File{f.Name(), "./uploads/" + f.Name()}
		Files = append(Files, file)
	}

	fp := path.Join("../templates", "list.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(rw, Files); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func UploadNameHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.ServeFile(rw, r, "../uploads/"+p.ByName("name"))
}
