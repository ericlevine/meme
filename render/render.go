package render

import (
  "flag"
  "image"
  "io/ioutil"
  "math"
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

type textBounds struct {
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
  writeLines(prepString(topText, lineLength), true, &ctx)
  writeLines(prepString(bottomText, lineLength), false, &ctx)

  return i, nil
}

func prepString(s string, lineLength int) []string {
  return splitString(strings.ToUpper(s), lineLength)
}

func initializeGraphicContext(ctx *memeContext, im *image.RGBA) {
  draw2d.SetFontFolder(".")

  ctx.gc = draw2d.NewGraphicContext(im)

  ctx.gc.SetStrokeColor(image.Black)
  ctx.gc.SetFillColor(image.White)
  ctx.gc.SetFontData(fontData)
  ctx.gc.SetFontSize(200)
}

func writeLines(lines []string, isTop bool, ctx *memeContext) {
  bounds := getBounds(lines, ctx)
  maxWidth := getMaxWidth(bounds)

  idealScale := (ctx.width - 2 * xMargin) / maxWidth
  scale := math.Min(idealScale, maxScale)

  nextTopOffset := 0.0
  if isTop {
    nextTopOffset = yMargin
  } else {
    // Start at the bottom edge
    nextTopOffset = ctx.height

    // Remove the bottom margin
    nextTopOffset -= yMargin

    // Remove the padding between the lines of text
    nextTopOffset -= yPadding * float64(len(lines) - 1)

    // Remove the height of the lines
    for _, b := range bounds {
      nextTopOffset -= (b.bottom - b.top) * scale
    }
  }

  for i, str := range lines {
    b := bounds[i]

    width := (b.right - b.left) * scale
    fullWidth := maxWidth * idealScale
    leftOffset := xMargin + ((fullWidth - width) / 2 - b.left * scale)

    topOffset := nextTopOffset

    height := b.bottom - b.top
    bottomOffset := topOffset + height * scale
    nextTopOffset = bottomOffset + yPadding

    writeString(str, leftOffset, bottomOffset, scale, ctx)
  }
}

func writeString(str string, leftOffset, bottomOffset, scale float64, ctx *memeContext) {
  ctx.gc.Restore()
  ctx.gc.Save()

  ctx.gc.Translate(leftOffset, bottomOffset)
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
