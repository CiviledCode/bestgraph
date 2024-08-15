package bestgraph

import (
	"fmt"
	"image/color"

	"github.com/civiledcode/bestgraph/chart"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	HUD_SIZE    = 128
	FOOTER_SIZE = 48
)

type ChartRender struct {
	Ctx             *chart.Chart
	Candles         []chart.Candle
	PressedKeys     []ebiten.Key
	ScreenWidth     int
	ScreenHeight    int
	KeyTick         uint64
	ShouldClose     bool
	CandleIntervals []float32
	CurrentInterval int
	TickRate        uint64
	CandlesPerMove  int
	DefaultFont     *text.GoTextFaceSource
}

func (g *ChartRender) StartChart() error {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowFloating(false)
	ebiten.SetWindowTitle("BestGraph - " + g.Ctx.SymbolName)

	if err := ebiten.RunGame(g); err != nil {
		return err
	}

	return nil
}

func (g *ChartRender) Update() error {
	if g.ShouldClose {
		return ebiten.Termination
	}

	g.PressedKeys = inpututil.AppendPressedKeys(g.PressedKeys[:0])
	return nil
}

func (g *ChartRender) Draw(screen *ebiten.Image) {
	// uint64 max
	if g.KeyTick == 18446744073709551615 {
		g.KeyTick = 0
	} else {
		g.KeyTick++
	}

	for _, key := range g.PressedKeys {
		if key == ebiten.KeyH {
			g.drawHelp(screen)
			return
		}

		if key == ebiten.KeyQ {
			g.ShouldClose = true
			return
		}
	}

	if g.KeyTick%g.TickRate == 0 {
		for _, key := range g.PressedKeys {
			if g.KeyTick%(g.TickRate*5) == 0 {
				if key == ebiten.KeyO {
					if g.CurrentInterval < len(g.CandleIntervals)-1 {
						fmt.Println("Zoom Out")
						g.CurrentInterval++
						g.Ctx.SetWidthCandles(g.CandleIntervals[g.CurrentInterval])
					}
				} else if key == ebiten.KeyI {
					if g.CurrentInterval > 0 {
						fmt.Println("Zoom In")
						g.CurrentInterval--
						g.Ctx.SetWidthCandles(g.CandleIntervals[g.CurrentInterval])
					}
				}

				if key == ebiten.KeyJ {
					fmt.Println("Pip Price Increase")
					g.Ctx.PricePerPip += g.Ctx.PipPriceIncrement

				} else if key == ebiten.KeyK {
					if (g.Ctx.PricePerPip - g.Ctx.PipPriceIncrement) >= g.Ctx.PipPriceMinimum {
						fmt.Println("Pip Price Decrease")
						g.Ctx.PricePerPip -= g.Ctx.PipPriceIncrement
					} else {
						g.Ctx.PricePerPip = g.Ctx.PipPriceMinimum
					}
				}

				if key == ebiten.KeyN {
					fmt.Println("Candles Per Move Increased")
					g.CandlesPerMove++
				} else if key == ebiten.KeyM {
					if g.CandlesPerMove > 1 {
						fmt.Println("Candles Per Move Decreased")
						g.CandlesPerMove--
					}
				}
			}

			if key == ebiten.KeyW || key == ebiten.KeyArrowUp {
				g.Ctx.CurrentPip++
			} else if key == ebiten.KeyS || key == ebiten.KeyArrowDown {
				g.Ctx.CurrentPip--
			} else if key == ebiten.KeyA || key == ebiten.KeyArrowLeft {
				if g.Ctx.CurrentCandle-g.CandlesPerMove >= 0 {
					g.Ctx.CurrentCandle -= g.CandlesPerMove
				}
			} else if key == ebiten.KeyD || key == ebiten.KeyArrowRight {
				g.Ctx.CurrentCandle += g.CandlesPerMove
			}

		}
	}

	//g.Ctx.DrawCandle(screen, g.Ctx.Candles[0], g.Ctx.GetCandleXPixel(g.Ctx.CurrentCandle))
	g.Ctx.DrawPipPoints(screen, float32(g.ScreenWidth), float32(g.ScreenHeight))
	g.Ctx.DrawCandles(screen)
	g.Ctx.DrawHud(screen, float32(g.ScreenWidth), float32(g.ScreenHeight))
}

func (g *ChartRender) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.ScreenHeight = outsideHeight
	g.ScreenWidth = outsideWidth
	// Excluse header in useable size.
	useableHeight := outsideHeight - HUD_SIZE - FOOTER_SIZE
	// Set the chart size on resize.
	g.Ctx.SetHeightPixels(float32(int(useableHeight/int(g.Ctx.PipsHeight))) * g.Ctx.PipsHeight)
	g.Ctx.SetWidthPixels(float32(int(outsideWidth/int(g.Ctx.CandlesShown))) * g.Ctx.CandlesShown)
	return outsideWidth, outsideHeight
}

func (g *ChartRender) drawHelp(screen *ebiten.Image) {
	helpMessage :=
		`[H] - Show this message
[D or →] Increment Current Candle
[A or ←] Decrement Current Candle
[W or ↑] Increment Current Pip
[S or ↓] Decrement Current Pip
[I] Decrease Candles Shown
[O] Increase Candles Shown
[J] Increase Price Per Pip
[K] Decrease Price Per Pip
[N] Increase Candles Per Move
[M] Decrease Candles Per Move
[Q] Close Chart
`

	op := &text.DrawOptions{}
	op.LineSpacing = 20
	op.GeoM.Translate(10, 10)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, helpMessage, &text.GoTextFace{
		Source: g.DefaultFont,
		Size:   18,
	}, op)
}

func CreateChartRender(chart *chart.Chart, candleIntervals []float32, startingInterval int, font *text.GoTextFaceSource) *ChartRender {
	chartRenderer := &ChartRender{
		Ctx:             chart,
		CandlesPerMove:  1,
		CandleIntervals: candleIntervals,
		CurrentInterval: startingInterval,
		DefaultFont:     font,
		TickRate:        5,
	}

	// Hardcoded for now...
	chart.ViewportY = HUD_SIZE

	if chart.PipsHeight == 0 {
		chart.PipsHeight = 24
	}

	if chart.CandlesShown == 0 {
		chart.CandlesShown = candleIntervals[startingInterval]
	}

	return chartRenderer
}
