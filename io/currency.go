/*
 * Copyright (c) 2022 The Mof Authors
 */

package pio

import (
	"fmt"
	"github.com/mofcloud/mof-common/math"
	curLib "golang.org/x/text/currency"
	"strings"
)

func GetCurrency(in string) (string, float64, error) {
	if len(in) < 1 {
		return "", 0.00, fmt.Errorf("failed to get currency from [%s]", in)
	}

	var cur string
	amount := 0.00

	switch in[0] {
	case '$':
		cur = curLib.USD.String()
		raw := strings.Replace(in[1:], ",", "", -1)
		if v, err := pmath.RoundFloat64FromString(raw, 2); err != nil {
			return "", 0.00, fmt.Errorf("failed to parse amount value from [%s]", in)
		} else {
			amount = v
		}
	case '¥':
		cur = curLib.CNY.String()
		raw := strings.Replace(in[1:], ",", "", -1)
		if v, err := pmath.RoundFloat64FromString(raw, 2); err != nil {
			return "", 0.00, fmt.Errorf("failed to parse amount value from [%s]", in)
		} else {
			amount = v
		}
	}

	return cur, amount, nil
}
