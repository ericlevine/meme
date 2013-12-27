package main

import (
  "bufio"
  "fmt"
  "image"
  "os"
  "github.com/ericlevine/meme/writer"
)

func main() {
  inFilename := "jason.gif"
  outFilename := "out.gif"
  top := "hello there my name is eric"
  bottom := "i think that you are really great"

  in, err := os.Open(inFilename)
  if err != nil { errorExit("Error opening input file.", err) }
  defer in.Close()

  outFile, err := os.Create(outFilename)
  if err != nil { errorExit("Error creating output file.", err) }
  defer outFile.Close()
  out := bufio.NewWriter(outFile)

  _, format, err := image.DecodeConfig(in)
  if err != nil { errorExit("Error decoding format.", err) }

  _, err = in.Seek(0, 0)
  if err != nil { errorExit("Error seeking back to beginning of file.", err) }

  if format == "gif" {
    writer.WriteMemeGIF(in, out, top, bottom)
  } else if format == "jpeg" {
    writer.WriteMemeJPEG(in, out, top, bottom)
  } else {
    fmt.Println("Unsupported image format.")
    os.Exit(1)
  }

  err = out.Flush()
  if err != nil { errorExit("Error writing file.", err) }
}

func errorExit(message string, err error) {
  fmt.Println(message)
  fmt.Println(err)
  os.Exit(1)
}
