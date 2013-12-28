package main

import (
  "bufio"
  "flag"
  "fmt"
  "image"
  "os"
  "github.com/ericlevine/meme/writer"
)

var (
  inFilename, outFilename, top, bottom string
)

func init() {
  flag.StringVar(&inFilename, "in", "", "Input image file.")
  flag.StringVar(&outFilename, "out", "", "Output image file.")
  flag.StringVar(&top, "top", "", "Top meme text.")
  flag.StringVar(&bottom, "bottom", "", "Bottom meme text.")
}

func main() {
  flag.Parse()

  if inFilename == "" || outFilename == "" {
    fmt.Println("-in and -out flags are required")
    os.Exit(1)
  }

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
    err = writer.WriteMemeGIF(in, out, top, bottom)
    if err != nil { errorExit("Could not create meme.", err) }
  } else if format == "jpeg" {
    err = writer.WriteMemeJPEG(in, out, top, bottom)
    if err != nil { errorExit("Could not create meme.", err) }
  } else if format == "png" {
    err = writer.WriteMemePNG(in, out, top, bottom)
    if err != nil { errorExit("Could not create meme.", err) }
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
