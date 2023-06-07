package emath

import (
	"math"
	"math/big"
)

func BigIntQuoDecimal(amount *big.Int, decimals int) float64 {
	b := new(big.Float).SetInt(amount)
	r, _ := new(big.Float).Quo(b, big.NewFloat(math.Pow10(decimals))).Float64()
	return r
}

func BigFloatQuoDecimal(b *big.Float, decimal float64) float64 {
	r, _ := new(big.Float).Quo(b, big.NewFloat(math.Pow(10, decimal))).Float64()
	return r
}

func BigIntToFloat64(b *big.Int) float64 {
	f := new(big.Float).SetInt(b)
	result, _ := f.Float64()
	return result

}
