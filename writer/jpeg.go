package writer

import (
  "io"
  "image/jpeg"
  "github.com/ericlevine/meme/render"
)

func WriteMemeJPEG(r io.Reader, w io.Writer, top, bottom string) error {
  i, err := jpeg.Decode(r)
  if err != nil { return err }
  meme, err := render.CreateMeme(i, top, bottom)
  if err != nil { return err }
  err = jpeg.Encode(w, meme, nil)
  if err != nil { return err }
  return nil
}
