/*
 * Copyright (c) 2022 The Mof Authors
 */

package pio

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var (
	Flatter = NewFlatter(".")
)

func NewFlatter(seperator string) *flattenToAttributes {
	return &flattenToAttributes{
		Seperator: seperator,
	}
}

type flattenToAttributes struct {
	Seperator string
}

func (f *flattenToAttributes) ToAttributes(input interface{}) ([]*Attribute, error) {
	res := make([]*Attribute, 0)

	target := map[string]string{}

	switch input.(type) {
	case []interface{}, map[string]interface{}:
		if err := f.helper(true, input, target, ""); err != nil {
			return res, err
		}
	default:
		src := map[string]interface{}{}

		// convert to bytes
		b, err := json.Marshal(input)
		if err != nil {
			return res, err
		}

		// convert to map
		err = json.Unmarshal(b, &src)
		if err != nil {
			return res, err
		}

		if err := f.helper(true, src, target, ""); err != nil {
			return res, err
		}
	}

	for k, v := range target {
		res = append(res, &Attribute{
			Key:   k,
			Value: v,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return strings.Compare(res[i].Key, res[j].Key) < 1
	})

	return res, nil
}

func (f *flattenToAttributes) helper(leading bool, src interface{}, target map[string]string, keyPrefix string) error {
	switch src.(type) {
	case map[string]interface{}:
		for k, v := range src.(map[string]interface{}) {
			newKey := f.joinKey(leading, keyPrefix, k)
			if err := f.helper(false, v, target, newKey); err != nil {
				return err
			}
		}
	case []interface{}:
		for i, v := range src.([]interface{}) {
			newKey := f.joinKey(leading, keyPrefix, strconv.Itoa(i))
			if err := f.helper(false, v, target, newKey); err != nil {
				return err
			}
		}
	default:
		target[keyPrefix] = fmt.Sprintf("%v", src)
	}

	return nil
}

func (f *flattenToAttributes) joinKey(isTop bool, prefix, subkey string) string {
	key := prefix

	if isTop {
		key += subkey
	} else {
		key += f.Seperator + subkey
	}

	return key
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Dimension struct {
	Key    string            `json:"key"`
	Values []*DimensionValue `json:"values"`
}

type DimensionValue struct {
	Value     string `json:"value"`
	ValueName string `json:"valueName"`
}
