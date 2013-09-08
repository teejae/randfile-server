package main

// Serves a file of random data over http.

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
)

const defaultFileSize = 100 * 1024 * 1024 // 100 MB

type randServer struct{}

func (s *randServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	form := r.Form
	sizeParam := form.Get("size")
	sizeInt, err := strconv.ParseInt(sizeParam, 10, 0)
	if err != nil {
		sizeInt = defaultFileSize
	}

	size := int(sizeInt)

	h := w.Header()
	h.Add("Content-Description", "File Transfer")
	h.Add("Content-Type", "application/octet-stream")
	h.Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"random-%v.bin\"", size))
	h.Add("Content-Transfer-Encoding", "binary")
	h.Add("Connection", "Keep-Alive")
	h.Add("Expires", "0")
	h.Add("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	h.Add("Pragma", "public")
	h.Add("Content-Length", fmt.Sprint(size))
	b := make([]byte, 8)
	e := binary.LittleEndian
	random := rand.New(rand.NewSource(0))
	for i := 0; i < size/8; i++ {
		u := (uint64(random.Uint32()) | (uint64(random.Uint32()) << 4))
		e.PutUint64(b, u)
		w.Write(b)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.Handle("/rand", &randServer{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
