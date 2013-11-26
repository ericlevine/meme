package main

import (
  "fmt"
  "image"
  "image/png"
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

func main() {
  ctx := MemeContext{
    height: 1000.0,
    width: 1000.0,
  }

  i := image.NewRGBA(image.Rect(0, 0, int(ctx.height), int(ctx.width)))
  initializeGraphicContext(&ctx, i)

  left, top, right, bottom := ctx.gc.GetStringBounds("MY NAME IS ERIC")
  fmt.Println(left)
  fmt.Println(top)
  fmt.Println(right)
  fmt.Println(bottom)

  scale := (ctx.width - 2 * xMargin) / (right - left)
  fmt.Println(scale)

  ctx.gc.Translate(-left * scale + xMargin, (bottom - top) * scale + yMargin)
  ctx.gc.Scale(scale, scale)

  ctx.gc.SetStrokeColor(image.White)
  ctx.gc.SetLineWidth(15.0)
  ctx.gc.StrokeString("MY NAME IS ERIC")

  ctx.gc.FillString("MY NAME IS ERIC")
  ctx.gc.Translate(1, 0)

  ctx.gc.Translate(-left * scale + xMargin, (bottom - top) * scale + yMargin)

  //draw2d.Rect(ctx.gc, left, top, right, bottom)
  //ctx.gc.SetLineWidth(3.0)
  //ctx.gc.Stroke()

  filename := "out1.png"
  f, err := os.Create(filename)
  if err != nil { os.Exit(1) }
  defer f.Close()
  b := bufio.NewWriter(f)
  err = png.Encode(b, i)
  if err != nil { os.Exit(1) }
  err = b.Flush()
  if err != nil { os.Exit(1) }
  fmt.Println("File written.")
}

func initializeGraphicContext(ctx *MemeContext, im *image.RGBA) {
  draw2d.SetFontFolder(".")

  ctx.gc = draw2d.NewGraphicContext(im)
  ctx.gc.SetStrokeColor(image.White)
  ctx.gc.SetFillColor(image.Black)
  ctx.gc.SetFontData(draw2d.FontData{"impact", draw2d.FontFamilySans, 0})
  ctx.gc.SetFontSize(200)
}

func writeTopLine(xOffset, yOffset float64) {

}
