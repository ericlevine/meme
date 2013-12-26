package main

import (
  "strings"
  "image"
  "image/png"
  "image/jpeg"
  "os"
  "bufio"
  "code.google.com/p/draw2d/draw2d"
)

var (
  yPadding = 20.0
  xMargin  = 20.0
  yMargin  = 20.0
)

type memeLine struct {
  x, y, scale float64
  text string
}

type memeText struct {
  top, bottom []memeLine
}

type memeContext struct {
  height, width float64
  gc draw2d.GraphicContext
}

type textBounds struct {
  left, top, right, bottom float64
}

func main() {
  inImage, err := openImage("image.jpeg")
  if err != nil { os.Exit(1) }

  i, err := createMeme(inImage, "hello there my name is eric", "i think that you are really great")
  if err != nil { os.Exit(1) }

  writeImage("out1.png", i)
  if err != nil { os.Exit(1) }
}

func createMeme(background image.Image, topText, bottomText string) (image.Image, error) {
  topText = strings.ToUpper(topText)
  bottomText = strings.ToUpper(bottomText)

  bounds := background.Bounds()
  height := float64(bounds.Max.Y - bounds.Min.Y)
  width := float64(bounds.Max.X - bounds.Min.X)

  ctx := memeContext{
    height: height,
    width: width,
  }

  i := image.NewRGBA(image.Rect(0, 0, int(ctx.width), int(ctx.height)))
  initializeGraphicContext(&ctx, i)

  ctx.gc.DrawImage(background)

  lineLength := int(width / height * 18)
  text := memeText{
    top: topLines(splitString(topText, lineLength), &ctx),
    bottom: bottomLines(splitString(bottomText, lineLength), &ctx),
  }

  for _, line := range text.top {
    writeString(line.text, line.x, line.y, line.scale, &ctx)
  }

  for _, line := range text.bottom {
    writeString(line.text, line.x, line.y, line.scale, &ctx)
  }

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

func initializeGraphicContext(ctx *memeContext, im *image.RGBA) {
  draw2d.SetFontFolder(".")

  ctx.gc = draw2d.NewGraphicContext(im)

  ctx.gc.SetStrokeColor(image.Black)
  ctx.gc.SetFillColor(image.White)
  ctx.gc.SetFontData(draw2d.FontData{"impact", draw2d.FontFamilySans, 0})
  ctx.gc.SetFontSize(200)
}

func topLines(lines []string, ctx *memeContext) []memeLine {
  bounds := getBounds(lines, ctx)
  maxWidth := getMaxWidth(bounds)

  scale := (ctx.width - 2 * xMargin) / maxWidth

  lastTopOffset := yMargin - yPadding

  results := make([]memeLine, len(lines))

  for i, str := range lines {
    b := bounds[i]

    width := b.right - b.left
    leftOffset := xMargin + ((maxWidth - width) / 2 - b.left) * scale

    height := b.bottom - b.top
    topOffset := lastTopOffset + (height * scale + yPadding)

    lastTopOffset = topOffset

    results[i] = memeLine{
      x: leftOffset,
      y: topOffset,
      scale: scale,
      text: str,
    }
  }

  return results
}

func bottomLines(lines []string, ctx *memeContext) []memeLine {
  bounds := getBounds(lines, ctx)
  maxWidth := getMaxWidth(bounds)

  scale := (ctx.width - 2 * xMargin) / maxWidth

  lastTopOffset := ctx.height + (yMargin - yPadding)

  results := make([]memeLine, len(lines))

  for i := len(lines) - 1; i >= 0; i-- {
    str := lines[i]
    b := bounds[i]

    width := b.right - b.left
    leftOffset := xMargin + ((maxWidth - width) / 2 - b.left) * scale

    height := b.bottom - b.top
    topOffset := lastTopOffset - yPadding

    lastTopOffset = topOffset - height * scale

    results[i] = memeLine{
      x: leftOffset,
      y: topOffset,
      scale: scale,
      text: str,
    }
  }

  return results
}

func writeString(str string, leftOffset, topOffset, scale float64, ctx *memeContext) {
  ctx.gc.Restore()
  ctx.gc.Save()

  ctx.gc.Translate(leftOffset, topOffset)
  ctx.gc.Scale(scale, scale)

  ctx.gc.SetLineWidth(15.0)
  ctx.gc.StrokeString(str)

  ctx.gc.FillString(str)
}

func getBounds(lines []string, ctx *memeContext) []textBounds {
  bounds := make([]textBounds, len(lines))
  for i, s := range lines {
    left, top, right, bottom := ctx.gc.GetStringBounds(s)
    bounds[i] = textBounds{left, top, right, bottom}
  }
  return bounds
}

func getMaxWidth(bounds []textBounds) float64 {
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
