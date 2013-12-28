package main

import (
  "bytes"
  "flag"
  "fmt"
  "net/http"
  "os"
  "github.com/ericlevine/meme/writer"
  "github.com/ericlevine/meme/upload"
)

var (
  port int
)

func init() {
  flag.IntVar(&port, "port", 8080, "Web server port.")
}

func memeHandler(w http.ResponseWriter, r *http.Request) {
  inFilename := r.PostFormValue("url")
  top := r.PostFormValue("top")
  bottom := r.PostFormValue("bottom")

  inResponse, err := http.Get(inFilename)
  if err != nil { panic(err.Error()) }
  in := inResponse.Body

  out := bytes.Buffer{}

  contentType := inResponse.Header["Content-Type"][0]
  outFilename := ""

  if contentType == "image/gif" {
    writer.WriteMemeGIF(in, &out, top, bottom)
    outFilename = "memes/out.gif"
  } else if contentType == "image/jpeg" {
    writer.WriteMemeJPEG(in, &out, top, bottom)
    outFilename = "memes/out.jpg"
  } else if contentType == "image/png" {
    writer.WriteMemePNG(in, &out, top, bottom)
    outFilename = "memes/out.png"
  } else {
    fmt.Println("Unsupported image format.")
    os.Exit(1)
  }

  url, err := upload.Write(out.Bytes(), outFilename, contentType)
  if err != nil { panic(err.Error()) }

  fmt.Fprintf(w, "%s\n", url)
}

func main() {
  http.HandleFunc("/", memeHandler)
  http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
