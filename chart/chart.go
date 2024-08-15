package chart

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// The chart is essentially a grid at a base level. Each price point is a row while each time interval is a column.
// When making that distiction, it's intuitive to make a system that can traverse said grid in a positonal fashion.
// The chart structure acts as a medium between traversing and drawing the chart.
type Chart struct {
	SymbolName string
	Candles    []Candle

	// STYLE

	// Colors for candles.
	Bearish color.Color
	Bullish color.Color

	PipLine    color.Color
	PriceColor color.Color

	HudTextColor color.Color

	PriceFont     *text.GoTextFaceSource
	PriceFontSize float64

	HudFont     *text.GoTextFaceSource
	HudFontSize float64

	WickWidth float32

	// How many pixels within the candle on the right and left should blank space exist.
	CandlePadding float32

	//The X and Y coordinates on the screen for the top left corner of the viewport.
	ViewportX float32
	ViewportY float32

	// HEIGHT

	// How many pips are fit within the chart viewport.
	PipsHeight float32
	// How many pixels the chart viewport is.
	ViewportHeight float32
	// How many pixels each pip is. This should not be manually edited.
	PixelsPerPip float32

	// WIDTH
	// How many candles are shown at once in the viewport.
	CandlesShown float32
	// The width of the viewport.
	ViewportWidth float32
	// The width of the candles. This should not by manually edited.
	CandleWidth float32

	// POSITION
	// These positions represent the cell in the top left corner of the viewport.
	// The current candle that is being shown.
	CurrentCandle int
	// The current pip level.
	CurrentPip int

	// PRICE

	// ZeroPrice is the open price of the first candle.
	ZeroPrice float32
	// PricePerPip is how much the price must increase for a pip.
	PricePerPip       float32
	PipPriceIncrement float32
	PipPriceMinimum   float32
}

func (c *Chart) SetHeightPips(pips float32) {
	c.PipsHeight = pips
	c.PixelsPerPip = c.ViewportHeight / pips
}

func (c *Chart) SetHeightPixels(pixels float32) {
	c.ViewportHeight = pixels
	c.PixelsPerPip = pixels / c.PipsHeight
}

func (c *Chart) SetWidthCandles(candles float32) {
	c.CandlesShown = candles
	c.CandleWidth = c.ViewportWidth / candles
}

func (c *Chart) SetWidthPixels(pixels float32) {
	c.ViewportWidth = pixels
	c.CandleWidth = pixels / c.CandlesShown
}

func (c *Chart) DrawPipPoints(screen *ebiten.Image, screenWidth, screenHeight float32) {
	vector.DrawFilledRect(screen, 0, c.ViewportY, screenWidth, 1, c.PipLine, true)
	for i := float32(1); i < c.PipsHeight; i++ {
		pricePoint := ((float32(c.CurrentPip) - i) * c.PricePerPip) + c.ZeroPrice
		vector.DrawFilledRect(screen, c.ViewportX, c.ViewportY+(i*c.PixelsPerPip), c.ViewportWidth, 1, c.PipLine, true)
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(c.ViewportWidth)-48, float64(c.ViewportY+(i*c.PixelsPerPip)))
		op.ColorScale.ScaleWithColor(c.PriceColor)
		text.Draw(screen, strconv.FormatFloat(float64(pricePoint), 'f', 4, 32), &text.GoTextFace{
			Source: c.PriceFont,
			Size:   c.PriceFontSize,
		}, op)
	}
}

func percentDifference(a, b float64) float64 {
	return (math.Abs(a-b) / ((a + b) / 2)) * 100
}

// TODO: Move this to ChartRender
func (c *Chart) DrawHud(screen *ebiten.Image, screenWidth, screenHeight float32) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(10, 16)
	op.ColorScale.ScaleWithColor(c.HudTextColor)
	s := fmt.Sprintf("Candle Idx: %d", c.CurrentCandle)
	text.Draw(screen, s, &text.GoTextFace{
		Source: c.HudFont,
		Size:   c.HudFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(10, 48)
	op.ColorScale.ScaleWithColor(c.HudTextColor)
	s = fmt.Sprintf("Pip Idx: %d", c.CurrentPip)
	text.Draw(screen, s, &text.GoTextFace{
		Source: c.HudFont,
		Size:   c.HudFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(10, 80)
	op.ColorScale.ScaleWithColor(c.HudTextColor)
	s = fmt.Sprintf("Candles: %d", len(c.Candles))
	text.Draw(screen, s, &text.GoTextFace{
		Source: c.HudFont,
		Size:   c.HudFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(180, 16)
	op.ColorScale.ScaleWithColor(c.HudTextColor)
	s = fmt.Sprintf("Pip Price: %0.4f", c.PricePerPip)
	text.Draw(screen, s, &text.GoTextFace{
		Source: c.HudFont,
		Size:   c.HudFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(180, 48)
	op.ColorScale.ScaleWithColor(c.HudTextColor)
	s = fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS())
	text.Draw(screen, s, &text.GoTextFace{
		Source: c.HudFont,
		Size:   c.HudFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(180, 80)
	op.ColorScale.ScaleWithColor(c.HudTextColor)
	s = fmt.Sprintf("Period: %0.2f", c.CandlesShown)
	text.Draw(screen, s, &text.GoTextFace{
		Source: c.HudFont,
		Size:   c.HudFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(180, 80)
	op.ColorScale.ScaleWithColor(c.HudTextColor)
	s = fmt.Sprintf("Period: %0.2f", c.CandlesShown)
	text.Draw(screen, s, &text.GoTextFace{
		Source: c.HudFont,
		Size:   c.HudFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(360, 16)
	op.ColorScale.ScaleWithColor(c.HudTextColor)
	text.Draw(screen, c.SymbolName, &text.GoTextFace{
		Source: c.HudFont,
		Size:   c.HudFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(360, 48)
	percent := 0.0
	if c.CurrentCandle < len(c.Candles) {
		percent = percentDifference(float64(c.Candles[c.CurrentCandle].Close), float64(c.ZeroPrice)) * 100
		if c.ZeroPrice > c.Candles[c.CurrentCandle].Close {
			op.ColorScale.ScaleWithColor(c.Bearish)
		} else {
			op.ColorScale.ScaleWithColor(c.Bullish)
		}
	}
	s = fmt.Sprintf("Zero Diff: %0.2f%s", percent, "%")
	text.Draw(screen, s, &text.GoTextFace{
		Source: c.HudFont,
		Size:   c.HudFontSize,
	}, op)

}

// Convert price into the pip increment relative open price of the first candle in the chart.
func (c *Chart) PriceToPipIndex(price float32) int {
	return int((price - c.ZeroPrice) / c.PricePerPip)
}

// Given the current position of the seek in the chart get the position within the viewport that a
// price point sits at.
func (c *Chart) PriceToYPixel(price float32, lowVal, highVal float32) float32 {
	pips := c.PriceToPipIndex(price)
	lowRange := c.CurrentPip - int(c.PipsHeight)
	if lowRange <= pips && pips <= c.CurrentPip {
		return c.ViewportY + (float32(c.CurrentPip-pips) * c.PixelsPerPip)
	} else {
		if lowRange > pips {
			return lowVal
		} else {
			return highVal
		}

	}
}

// Return the min and max prices in the viewport.
func (c *Chart) Range() (min float32, max float32) {
	max = (float32(c.CurrentPip) * c.PricePerPip) + c.ZeroPrice
	min = max - (c.PipsHeight * c.PricePerPip)
	return
}

func (c *Chart) GetCandlesInView() []Candle {
	inview := make([]Candle, int(c.CandlesShown))
	min, max := c.Range()

	var to int
	if c.CurrentCandle+int(c.CandlesShown) < len(c.Candles) {
		to = c.CurrentCandle + int(c.CandlesShown)
	} else {
		to = len(c.Candles)
	}

	var candle Candle
	for i := c.CurrentCandle; i < to; i++ {
		candle = c.Candles[i]
		if candle.High <= max && candle.High >= min {
			inview[i-c.CurrentCandle] = candle
		} else if candle.Low <= max && candle.Low >= min {
			inview[i-c.CurrentCandle] = candle
		}
	}
	return inview
}

func (c *Chart) GetCandleXPixel(candleIndex int) float32 {
	if candleIndex < c.CurrentCandle {
		return -1
	} else if candleIndex > c.CurrentCandle+int(c.CandlesShown) {
		return -2
	}

	return c.ViewportX + (float32(candleIndex-c.CurrentCandle) * c.CandleWidth)
}

func (c *Chart) DrawCandle(screen *ebiten.Image, candle Candle, candleX float32) {
	if candle.Close < candle.Open {
		// Bearish
		bottom := c.ViewportY + c.ViewportHeight
		wickBottom := c.PriceToYPixel(candle.Low, bottom, -1)
		if wickBottom == -1 {
			// Dont render if the lowest point is above the highest price point shown.
			return
		}

		wickTop := c.PriceToYPixel(candle.High, -1, c.ViewportY)
		if wickTop == -1 {
			// Dont render if the highest point is below the lowest price point shown.
			return
		}

		candleBottom := c.PriceToYPixel(candle.Open, bottom, c.ViewportY)
		candleTop := c.PriceToYPixel(candle.Close, bottom, c.ViewportY)
		if candleTop == candleBottom {
			candleTop -= 5
		}
		// Draw Wick.
		vector.DrawFilledRect(screen, candleX+((c.CandleWidth/2)-(c.WickWidth/2)), wickTop, 2, wickBottom-wickTop, c.Bearish, true)
		// Draw body.
		vector.DrawFilledRect(screen, candleX+c.CandlePadding, candleTop, c.CandleWidth-(2*c.CandlePadding), candleBottom-candleTop, c.Bearish, true)
	} else {
		// Bullish
		bottom := c.ViewportY + c.ViewportHeight
		wickBottom := c.PriceToYPixel(candle.Low, bottom, -1)
		if wickBottom == -1 {
			// Dont render if the lowest point is above the highest price point shown.
			return
		}

		wickTop := c.PriceToYPixel(candle.High, -1, c.ViewportY)
		if wickTop == -1 {
			// Dont render if the highest point is below the lowest price point shown.
			return
		}

		candleTop := c.PriceToYPixel(candle.Close, bottom, c.ViewportY)
		candleBottom := c.PriceToYPixel(candle.Open, bottom, c.ViewportY)
		if candleTop == candleBottom {
			candleTop -= 5
		}
		// Draw Wick.
		vector.DrawFilledRect(screen, candleX+((c.CandleWidth/2)-(c.WickWidth/2)), wickTop, 2, wickBottom-wickTop, c.Bullish, true)
		// Draw body.
		vector.DrawFilledRect(screen, candleX+c.CandlePadding, candleTop, c.CandleWidth-(2*c.CandlePadding), candleBottom-candleTop, c.Bullish, true)
	}
}

func (c *Chart) DrawCandles(screen *ebiten.Image) {
	if c.CurrentCandle >= len(c.Candles) {
		return
	}

	candles := c.GetCandlesInView()
	for i, candle := range candles {
		if candle.Open == 0 {
			continue
		}

		x := c.GetCandleXPixel(i + c.CurrentCandle)
		if x < 0 {
			continue
		}
		c.DrawCandle(screen, candle, x)
	}
}

func (c *Chart) DrawOffscreenCandleArrow(screen *ebiten.Image, candle Candle, candleX float32, top bool, bullish bool) {
	// TODO
}
