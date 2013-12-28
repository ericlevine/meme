package writer

import (
  "io"
  "image/png"
  "github.com/ericlevine/meme/render"
)

func WriteMemePNG(r io.Reader, w io.Writer, top, bottom string) error {
  i, err := png.Decode(r)
  if err != nil { return err }
  meme, err := render.CreateMeme(i, top, bottom)
  if err != nil { return err }
  err = png.Encode(w, meme)
  if err != nil { return err }
  return nil
}
