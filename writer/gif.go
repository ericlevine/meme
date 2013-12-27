package writer

import (
  "io"
  "image"
  "image/color"
  "image/gif"
  "math"
  "github.com/ericlevine/meme/render"
)

func WriteMemeGIF(r io.Reader, w io.Writer, top, bottom string) error {
  i, err := gif.DecodeAll(r)
  if err != nil { return err }
  // TODO: copy gif here
  err = createGifMeme(i, top, bottom)
  if err != nil { return err }
  err = gif.EncodeAll(w, i)
  if err != nil { return err }
  return nil
}

func createGifMeme(background *gif.GIF, topText, bottomText string) error {
  firstFrame := background.Image[0]

  overlayMeme, err := render.CreateMeme(
      image.NewRGBA(firstFrame.Bounds()), topText, bottomText)
  if err != nil { return err }

  for _, pic := range background.Image {
    injectColor(pic.Palette, image.White)
    injectColor(pic.Palette, image.Black)
    bounds := pic.Bounds()
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
      for x := bounds.Min.X; x < bounds.Max.X; x++ {
        color := overlayMeme.At(x, y)
        _, _, _, alpha := color.RGBA()
        if alpha > 0 {
          pic.Set(x, y, color)
        }
      }
    }
  }

  return nil
}

func injectColor(p color.Palette, target color.Color) {
  if len(p) < 256 {
    p = append(p, target)
  } else {
    tr, tg, tb, ta := target.RGBA()
    bestScore, bestIndex := math.MaxFloat64, 0
    for i, candidate := range p {
      cr, cg, cb, ca := candidate.RGBA()
      score := math.Pow(float64(tr - cr), 2)
      score += math.Pow(float64(tg - cg), 2)
      score += math.Pow(float64(tb - cb), 2)
      score += math.Pow(float64(ta - ca), 2)
      if score < bestScore {
        bestScore = score
        bestIndex = i
      }
    }
    p[bestIndex] = target
  }
}
