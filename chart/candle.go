package chart

import (
	"os"

	"github.com/gocarina/gocsv"
)

type Candle struct {
	High  float32 `csv:"high"`
	Low   float32 `csv:"low"`
	Open  float32 `csv:"open"`
	Close float32 `csv:"close"`
}

func CandlesFromCSVFile(fileName string) []Candle {
	var candles []Candle
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	gocsv.UnmarshalFile(f, &candles)
	return candles
}
