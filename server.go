package main

import (
  "bytes"
  "flag"
  "fmt"
  "net/http"
  "github.com/ericlevine/meme/writer"
  "github.com/ericlevine/meme/upload"
  "code.google.com/p/vitess/go/cache"
)

var (
  port int
  memeCache *cache.LRUCache = cache.NewLRUCache(50 * 1024 * 1024)
)

type memeCacheEntry struct {
  contentType string
  data []byte
}

func (m *memeCacheEntry) Size() int {
  return len(m.data) + len(m.contentType)
}

func init() {
  flag.IntVar(&port, "port", 8080, "Web server port.")
}

func memeHandler(w http.ResponseWriter, r *http.Request) {
  inFilename := r.PostFormValue("url")
  top := r.PostFormValue("top")
  bottom := r.PostFormValue("bottom")

  fmt.Printf("URL: %s\n", inFilename)
  fmt.Printf("Top text: %s\n", top)
  fmt.Printf("Bottom text: %s\n", bottom)

  rawCacheEntry, ok := memeCache.Get(inFilename)
  var cacheEntry *memeCacheEntry
  if ok {
    fmt.Println("Cache hit!")
    cacheEntry = rawCacheEntry.(*memeCacheEntry)
  } else {
    fmt.Println("Cache miss!")
    inResponse, err := http.Get(inFilename)
    if err != nil { panic(err.Error()) }
    inBuffer := bytes.Buffer{}
    inBuffer.ReadFrom(inResponse.Body)
    cacheEntry = &memeCacheEntry{
      contentType: inResponse.Header["Content-Type"][0],
      data: inBuffer.Bytes(),
    }
    memeCache.Set(inFilename, cacheEntry)
  }

  out := bytes.Buffer{}
  in := bytes.NewReader(cacheEntry.data)

  contentType := cacheEntry.contentType
  outFilename := ""

  fmt.Println("Generating meme.")
  if contentType == "image/gif" {
    err := writer.WriteMemeGIF(in, &out, top, bottom)
    if err != nil { panic(err.Error()) }
    outFilename = "memes/out.gif"
  } else if contentType == "image/jpeg" {
    err := writer.WriteMemeJPEG(in, &out, top, bottom)
    if err != nil { panic(err.Error()) }
    outFilename = "memes/out.jpg"
  } else if contentType == "image/png" {
    err := writer.WriteMemePNG(in, &out, top, bottom)
    if err != nil { panic(err.Error()) }
    outFilename = "memes/out.png"
  } else {
    panic("Unsupported image format.")
  }

  fmt.Println("Meme generated, uploading...")
  url, err := upload.Write(out.Bytes(), outFilename, contentType)
  if err != nil { panic(err.Error()) }

  fmt.Fprintf(w, "%s\n", url)
  fmt.Println("Completed request.")
}

func main() {
  flag.Parse()
  http.HandleFunc("/", memeHandler)
  http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
