package main

import (
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

type MemeContext struct {
  height, width float64
  gc draw2d.GraphicContext
}

type TextBounds struct {
  left, top, right, bottom float64
}

func main() {
  inImage, err := openImage("image.jpeg")
  if err != nil { os.Exit(1) }

  topStrs := []string{"HELLO THERE", "MY NAME IS ERIC"}
  bottomStrs := []string{"I THINK THAT", "YOU ARE REALLY GREAT"}

  bounds := inImage.Bounds()
  height := float64(bounds.Max.Y - bounds.Min.Y)
  width := float64(bounds.Max.X - bounds.Min.X)

  ctx := MemeContext{
    height: height,
    width: width,
  }

  i := image.NewRGBA(image.Rect(0, 0, int(ctx.width), int(ctx.height)))
  initializeGraphicContext(&ctx, i)

  ctx.gc.DrawImage(inImage)

  writeTop(topStrs, &ctx)
  writeBottom(bottomStrs, &ctx)

  writeImage("out1.png", i)
  if err != nil { os.Exit(1) }
}

func createMeme(inFile, outFile, topText, bottomText string) error {
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
