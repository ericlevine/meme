package render

import (
  "flag"
  "image"
  "io/ioutil"
  "strings"
  "code.google.com/p/draw2d/draw2d"
  "code.google.com/p/freetype-go/freetype"
)

var (
  fontFile string
  fontData = draw2d.FontData{"impact", draw2d.FontFamilySans, 0}
  fontLoaded = false

  yPadding = 20.0
  xMargin  = 20.0
  yMargin  = 20.0
  maxScale = 0.3
)

type memeContext struct {
  height, width float64
  gc draw2d.GraphicContext
}

type TextBounds struct {
  left, top, right, bottom float64
}

func init() {
  flag.StringVar(&fontFile, "font", "impactsr.ttf", "Font file.")
}

func loadFont() error {
  fontBytes, err := ioutil.ReadFile(fontFile)
  if err != nil { return err }
  font, err := freetype.ParseFont(fontBytes)
  if err != nil { return err }
  draw2d.RegisterFont(fontData, font)
  fontLoaded = true
  return nil
}

func CreateMeme(background image.Image, topText, bottomText string) (image.Image, error) {
  if !fontLoaded {
    err := loadFont()
    if err != nil { return nil, err }
  }

  topText = strings.ToUpper(topText)
  bottomText = strings.ToUpper(bottomText)

  bounds := background.Bounds()
  height := float64(bounds.Max.Y - bounds.Min.Y)
  width := float64(bounds.Max.X - bounds.Min.X)

  ctx := memeContext{
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

func initializeGraphicContext(ctx *memeContext, im *image.RGBA) {
  draw2d.SetFontFolder(".")

  ctx.gc = draw2d.NewGraphicContext(im)

  ctx.gc.SetStrokeColor(image.Black)
  ctx.gc.SetFillColor(image.White)
  ctx.gc.SetFontData(fontData)
  ctx.gc.SetFontSize(200)
}

func writeTop(lines []string, ctx *memeContext) {
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

func writeBottom(lines []string, ctx *memeContext) {
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

func writeString(str string, leftOffset, topOffset, scale float64, ctx *memeContext) {
  ctx.gc.Restore()
  ctx.gc.Save()

  ctx.gc.Translate(leftOffset, topOffset)
  ctx.gc.Scale(scale, scale)

  ctx.gc.SetLineWidth(15.0)
  ctx.gc.StrokeString(str)

  ctx.gc.FillString(str)
}

func getBounds(lines []string, ctx *memeContext) []TextBounds {
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
