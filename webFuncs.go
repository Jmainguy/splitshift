package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		currentTime := time.Now().Unix()
		h := md5.New()
		_, err := io.WriteString(h, strconv.FormatInt(currentTime, 10))
		if err != nil {
			t, _ := template.ParseFiles("error.gtpl")
			err = t.Execute(w, err.Error())
			if err != nil {
				panic(err)
			}
			return
		}
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, err := template.ParseFiles("upload.gtpl")
		if err != nil {
			t, _ := template.ParseFiles("error.gtpl")
			err = t.Execute(w, err.Error())
			if err != nil {
				panic(err)
			}
			return
		}
		err = t.Execute(w, token)
		if err != nil {
			t, _ := template.ParseFiles("error.gtpl")
			err = t.Execute(w, err.Error())
			if err != nil {
				panic(err)
			}
			return
		}
	} else {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			t, _ := template.ParseFiles("error.gtpl")
			err = t.Execute(w, err.Error())
			if err != nil {
				panic(err)
			}
			return
		}
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			t, _ := template.ParseFiles("error.gtpl")
			err = t.Execute(w, err.Error())
			if err != nil {
				panic(err)
			}
			return
		}

		// File has been downloaded by server
		// Time to process it and return the csv result

		resultLines, err := processFile(file)
		if err != nil {
			t, _ := template.ParseFiles("error.gtpl")
			err = t.Execute(w, err.Error())
			if err != nil {
				panic(err)
			}
			return
		}
		buf := new(bytes.Buffer)
		for _, line := range resultLines {
			buf.WriteString(line)
		}

		w.Header().Set("Content-Disposition", "attachment; filename=export.txt")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Transfer-Encoding", "chunked")

		_, err = buf.WriteTo(w)
		if err != nil {
			t, _ := template.ParseFiles("error.gtpl")
			err = t.Execute(w, err.Error())
			if err != nil {
				panic(err)
			}
			return
		}
	}
}
