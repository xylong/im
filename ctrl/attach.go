package ctrl

import (
	"fmt"
	"im/util"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	os.MkdirAll("./mnt", os.ModePerm)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	uploadLocal(w, r)
}

func uploadLocal(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		file   multipart.File
		header *multipart.FileHeader
		f      *os.File
	)

	if file, header, err = r.FormFile("file"); err != nil {
		util.Fail(w, err.Error())
	}
	suffix := ".png"
	fileName := header.Filename
	if arr := strings.Split(fileName, "."); len(arr) > 1 {
		suffix = "." + arr[len(arr)-1]
	}
	if fileType := r.FormValue("filetype"); len(fileType) > 0 {
		suffix = fileType
	}
	name := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	if f, err = os.Create("./mnt/" + name); err != nil {
		util.Fail(w, err.Error())
		return
	}
	if _, err = io.Copy(f, file); err != nil {
		util.Fail(w, err.Error())
		return
	}
	util.Success(w, "/mnt/"+name, "")
}
