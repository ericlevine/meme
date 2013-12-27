package main

import (
  "strings"
  "image"
  "image/color"
  "image/png"
  "image/gif"
  "image/jpeg"
  "math"
  "os"
  "bufio"
  "code.google.com/p/draw2d/draw2d"
)

var (
  yPadding = 20.0
  xMargin  = 20.0
  yMargin  = 20.0
)

type MemeContext struct {
  height, width float64
  gc draw2d.GraphicContext
}

type TextBounds struct {
  left, top, right, bottom float64
}

func main() {
  //inImage, err := openImage("image.jpeg")
  //if err != nil { os.Exit(1) }

  //i, err := createMeme(inImage, "hello there my name is eric", "i think that you are really great")
  //if err != nil { os.Exit(1) }

  //writeImage("out1.png", i)
  //if err != nil { os.Exit(1) }

  inImage, err := openGif("ian.gif")
  if err != nil { os.Exit(1) }

  err = createGifMeme(inImage, "hello there my name is eric", "i think that you are really great")
  if err != nil { os.Exit(1) }

  writeGif("out1.gif", inImage)
  if err != nil { os.Exit(1) }
}

func createGifMeme(background *gif.GIF, topText, bottomText string) error {
  firstFrame := background.Image[0]

  overlayMeme, err := createMeme(image.NewRGBA(firstFrame.Bounds()), topText, bottomText)
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

func createMeme(background image.Image, topText, bottomText string) (image.Image, error) {
  topText = strings.ToUpper(topText)
  bottomText = strings.ToUpper(bottomText)

  bounds := background.Bounds()
  height := float64(bounds.Max.Y - bounds.Min.Y)
  width := float64(bounds.Max.X - bounds.Min.X)

  ctx := MemeContext{
    height: height,
    width: width,
  }

  i := image.NewRGBA(background.Bounds())
  initializeGraphicContext(&ctx, i)

  ctx.gc.DrawImage(background)

  lineLength := int(width / height * 18)
  writeTop(splitString(topText, lineLength), &ctx)
  writeBottom(splitString(bottomText, lineLength), &ctx)

  return i, nil
}

func openImage(name string) (image.Image, error) {
  f, err := os.Open(name)
  if err != nil { return nil, err }
  defer f.Close()
  i, err := jpeg.Decode(f)
  if err != nil { return nil, err }
  return i, nil
}

func openGif(name string) (*gif.GIF, error) {
  f, err := os.Open(name)
  if err != nil { return nil, err }
  defer f.Close()
  i, err := gif.DecodeAll(f)
  if err != nil { return nil, err }
  return i, nil
}

func writeImage(name string, i image.Image) error {
  f, err := os.Create(name)
  if err != nil { return err }
  defer f.Close()
  b := bufio.NewWriter(f)
  err = png.Encode(b, i)
  if err != nil { return err }
  err = b.Flush()
  if err != nil { return err }
  return nil
}

func writeGif(name string, i *gif.GIF) error {
  f, err := os.Create(name)
  if err != nil { return err }
  defer f.Close()
  b := bufio.NewWriter(f)
  err = gif.EncodeAll(b, i)
  if err != nil { return err }
  err = b.Flush()
  if err != nil { return err }
  return nil
}

func initializeGraphicContext(ctx *MemeContext, im *image.RGBA) {
  draw2d.SetFontFolder(".")

  ctx.gc = draw2d.NewGraphicContext(im)

  ctx.gc.SetStrokeColor(image.Black)
  ctx.gc.SetFillColor(image.White)
  ctx.gc.SetFontData(draw2d.FontData{"impact", draw2d.FontFamilySans, 0})
  ctx.gc.SetFontSize(200)
}

func writeTop(lines []string, ctx *MemeContext) {
  bounds := getBounds(lines, ctx)
  maxWidth := getMaxWidth(bounds)

  scale := (ctx.width - 2 * xMargin) / maxWidth

  lastTopOffset := yMargin - yPadding

  for i, str := range lines {
    b := bounds[i]

    width := b.right - b.left
    leftOffset := xMargin + ((maxWidth - width) / 2 - b.left) * scale

    height := b.bottom - b.top
    topOffset := lastTopOffset + (height * scale + yPadding)

    lastTopOffset = topOffset

    writeString(str, leftOffset, topOffset, scale, ctx)
  }
}

func writeBottom(lines []string, ctx *MemeContext) {
  bounds := getBounds(lines, ctx)
  maxWidth := getMaxWidth(bounds)

  scale := (ctx.width - 2 * xMargin) / maxWidth

  lastTopOffset := ctx.height + (yMargin - yPadding)

  for i := len(lines) - 1; i >= 0; i-- {
    str := lines[i]
    b := bounds[i]

    width := b.right - b.left
    leftOffset := xMargin + ((maxWidth - width) / 2 - b.left) * scale

    height := b.bottom - b.top
    topOffset := lastTopOffset - yPadding

    lastTopOffset = topOffset - height * scale

    writeString(str, leftOffset, topOffset, scale, ctx)
  }
}

func writeString(str string, leftOffset, topOffset, scale float64, ctx *MemeContext) {
  ctx.gc.Restore()
  ctx.gc.Save()

  ctx.gc.Translate(leftOffset, topOffset)
  ctx.gc.Scale(scale, scale)

  ctx.gc.SetLineWidth(15.0)
  ctx.gc.StrokeString(str)

  ctx.gc.FillString(str)
}

func getBounds(lines []string, ctx *MemeContext) []TextBounds {
  bounds := make([]TextBounds, len(lines))
  for i, s := range lines {
    left, top, right, bottom := ctx.gc.GetStringBounds(s)
    bounds[i] = TextBounds{left, top, right, bottom}
  }
  return bounds
}

func getMaxWidth(bounds []TextBounds) float64 {
  maxWidth := 0.0
  for _, b := range bounds {
    width := b.right - b.left
    if width > maxWidth {
      maxWidth = width
    }
  }
  return maxWidth
}

func splitString(s string, length int) []string {
  strs := make([]string, 0)
  split := findSplit(s, length)
  for split != -1 {
    strs = append(strs, s[:split])
    s = s[split + 1:]
    split = findSplit(s, length)
  }
  strs = append(strs, s)
  return strs
}

func findSplit(s string, length int) int {
  if len(s) < length { return -1 }
  best := -1
  for split := 0; split != best; split = strings.Index(s[best+1:], " ") + best + 1 {
    if split > length {
      if absInt(length - split) > absInt(length - best) {
        return best
      } else {
        return split
      }
    }
    best = split
  }
  return -1
}

func absInt(i int) int {
  if i < 0 { return -i }
  return i
}
