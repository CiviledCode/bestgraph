# bestgraph 📈

[![Go Reference](https://pkg.go.dev/badge/github.com/civiledcode/bestgraph.svg)](https://pkg.go.dev/github.com/civiledcode/bestgraph)
[![Go Report Card](https://goreportcard.com/badge/github.com/civiledcode/bestgraph)](https://goreportcard.com/report/github.com/civiledcode/bestgraph)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/civiledcode/bestgraph.svg)](https://github.com/civiledcode/bestgraph/releases)

Welcome to **bestgraph**, the best candlestick graphing library you'll *ever* need! (*Disclaimer*: It’s not the best, but it’s definitely the best for me, and that's what really counts, right?)

Whether you're obsessed with Go, fascinated by market trends, or just want smooth grid navigation, **bestgraph** has you covered with "blazing" fast hardware-accelerated rendering, endless customization options, and a help page for those rare moments when you need it. 😏

## Why bestgraph?

- **Hardware Accelerated Rendering**: Why wait? We've harnessed the power of your GPU so your charts load faster.
- **Grid-Based Navigation**: Navigate like a pro! Our intuitive controls make zooming and panning feel like second nature. 🧭
- **High Customization**: Like everything in life, your charts should be exactly how you want them. Customize all the things.
- **Relatively Good Performance**: Well, it's *relatively* good. Don't push it.
- **Written in Go**: I mean, who doesn’t love Go? Go is life. Go is love. 🐹
- **Help Page**: Yes, you heard it. We even have a help page. No one's using it, but it's there.

## Still To-Do (Because, let’s face it, it’s not *actually* finished):

- Time intervals. For when you want to see more (or less) of the market's chaos. ⏳
- Callback functions in the renderer so you can render even more cool stuff (because who doesn't like extra stuff?).
- Similar viewport chunking for a variety of data structures. (Fancy, right? I don't know what it means either, but it's important.)
- Optimizations. (There’s probably a lot of room for improvement, but hey, it works for now, okay?)
- API for drawing consolidation lines. Because sometimes the market just needs to take a break.
- API for placing long/short position indicators. So you can feel like a real trader.
- API for drawing indicator lines. More lines, more power. 📊

## Getting Started


You'll need Go 1.19 or later because, you know, we're modern like that.

[Download Go here](https://go.dev/doc/install)

### Installing a C Compiler

Ebitengine requires a C compiler (we're not that fancy yet). Here's how you can install one:
```
sudo apt install gcc
```

### Dependencies
Debian / Ubuntu
```
sudo apt install libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
```
Fedora
```
sudo dnf install mesa-libGL-devel mesa-libGLES-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel alsa-lib-devel pkg-config
```
Solus
```
sudo eopkg install libglvnd-devel libx11-devel libxrandr-devel libxinerama-devel libxcursor-devel libxi-devel libxxf86vm-devel alsa-lib-devel pkg-config
```
Arch
```
sudo pacman -S mesa libxrandr libxcursor libxinerama libxi pkg-config
```
Alpine
```
sudo apk add alsa-lib-dev libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev mesa-dev pkgconf
```
Void
```
sudo xbps-install libXxf86vm-devel pkg-config
```

## Example Code
Here's a quick example to get you started with bestgraph. Because who reads all this without diving into code?

```go
func main() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	c := &chart.Chart{
		SymbolName:        "TEST/USDT",
		PipsHeight:        24,
		WickWidth:         2,
		Bearish:           color.RGBA{200, 0, 0, 255},
		Bullish:           color.RGBA{0, 200, 0, 255},
		PipLine:           color.RGBA{20, 20, 20, 100},
		PriceColor:        color.White,
		PriceFont:         s,
		PriceFontSize:     12,
		PricePerPip:       .0001,
		PipPriceIncrement: .0001,
		PipPriceMinimum:   .0001,
		HudFont:           s,
		HudFontSize:       20,
		HudTextColor:      color.White,
	}

	candles := chart.CandlesFromCSVFile("test.csv")
	slices.Reverse(candles)
	c.Candles = candles
	c.ZeroPrice = c.Candles[0].Open
	renderer := CreateChartRender(c, []float32{30, 60, 120, 240, 480, 960}, 1, s)
	err = renderer.StartChart()
	if err != nil {
		panic(err)
	}
}

```

## Contributing

Feel free to make a PR if you want to contribute.