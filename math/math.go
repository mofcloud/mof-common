/*
 * Copyright (c) 2022 The Mof Authors
 */

package pmath

import (
	"math"
	"strconv"
)

func RoundFloat64(in float64, len int) float64 {
	if len < 1 {
		len = 2
	}

	di := math.Pow10(len)
	return math.Round(in*di) / di
}

func RoundFloat64FromString(in string, len int) (float64, error) {
	raw, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0.00, err
	}

	if len < 1 {
		len = 2
	}

	di := math.Pow10(len)
	return math.Round(raw*di) / di, nil
}
