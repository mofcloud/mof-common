/*
 * Copyright (c) 2022 The Mof Authors
 */

package pio

import (
	"reflect"
	"sort"
)

func MapKeysString(in interface{}) []string {
	val := reflect.ValueOf(in)

	res := make([]string, 0)

	if val.Kind() == reflect.Map {
		for _, e := range val.MapKeys() {
			v := val.MapIndex(e)
			switch t := v.Interface().(type) {
			case string:
				res = append(res, t)
			}
		}
	}

	sort.Strings(res)

	return res
}

func MapKeysInt64(in interface{}) []int64 {
	val := reflect.ValueOf(in)

	res := make([]int64, 0)

	if val.Kind() == reflect.Map {
		for _, e := range val.MapKeys() {
			v := val.MapIndex(e)
			switch t := v.Interface().(type) {
			case int64:
				res = append(res, t)
			}
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})

	return res
}

func MapKeysInt(in interface{}) []int {
	val := reflect.ValueOf(in)

	res := make([]int, 0)

	if val.Kind() == reflect.Map {
		for _, e := range val.MapKeys() {
			v := val.MapIndex(e)
			switch t := v.Interface().(type) {
			case int:
				res = append(res, t)
			}
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})

	return res
}
