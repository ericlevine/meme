package writer

import (
  "flag"
  "fmt"
  "io"
  "image"
  "image/color"
  "image/gif"
  "math"
  "time"
  "github.com/ericlevine/meme/render"
)

var start time.Time

var (
  grays int
  grayBound int = (1 << 16) / 3
  whiteBound int = (1 << 17) / 3
)

func init() {
  flag.IntVar(&grays, "grays", 4, "Number of grays in gifs")
}

func WriteMemeGIF(r io.Reader, w io.Writer, top, bottom string) error {
  start = time.Now()
  i, err := gif.DecodeAll(r)
  fmt.Println("decoded.")
  fmt.Println(time.Now().Sub(start))
  if err != nil { return err }
  err = createGifMeme(i, top, bottom)
  fmt.Println("gif'd.")
  fmt.Println(time.Now().Sub(start))
  if err != nil { return err }
  err = gif.EncodeAll(w, i)
  fmt.Println("encoded.")
  fmt.Println(time.Now().Sub(start))
  if err != nil { return err }
  return nil
}

func createGifMeme(background *gif.GIF, topText, bottomText string) error {
  firstFrame := background.Image[0]

  overlayMeme, err := render.CreateMeme(
      image.NewRGBA(firstFrame.Bounds()), topText, bottomText)
  if err != nil { return err }

  completion := make(chan bool)
  for i, _ := range background.Image {
    go func(index int) {
      overlayMemeOnFrame(background.Image, index, overlayMeme)
      completion <- true
    }(i)
  }
  count := 0
  for count < len(background.Image) {
    <-completion
    count += 1
  }
  return nil
}

func overlayMemeOnFrame(pics []*image.Paletted, i int, overlay image.Image) {
  grayIndices := injectGrays(&pics[i].Palette)
  stride := uint32(0xffff / (grays - 1))
  bounds := pics[i].Bounds()
  for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
      color := overlay.At(x, y)
      red, _, _, alpha := color.RGBA()
      if alpha > 0 {
        colorIndex := grayIndices[red / stride]
        pics[i].Pix[(y-bounds.Min.Y) * pics[i].Stride + (x-bounds.Min.X)] = colorIndex
      }
    }
  }
}

func injectGrays(p *color.Palette) []uint8 {
  stride := 0xffff / (grays - 1)
  offsets := make([]uint8, grays)
  for i := 0; i < grays - 1; i++ {
    offsets[i] = injectColor(p, color.Gray16{uint16(stride * i)})
  }
  offsets[grays - 1] = injectColor(p, color.Gray16{0xffff})
  return offsets
}

func injectColor(p *color.Palette, target color.Color) uint8 {
  if len(*p) < 256 {
    *p = append(*p, target)
    return uint8(len(*p) - 1)
  } else {
    tr, tg, tb, ta := target.RGBA()
    bestScore, bestIndex := math.MaxFloat64, 0
    for i, candidate := range *p {
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
    (*p)[bestIndex] = target
    return uint8(bestIndex)
  }
}
