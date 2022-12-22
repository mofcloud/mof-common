/*
 * Copyright (c) 2022 The Mof Authors
 */

package pio

import (
	"encoding/json"
	"fmt"
	"strings"
)

func PrintStructPretty(in interface{}) {
	bytes, _ := json.MarshalIndent(in, "", "  ")
	fmt.Println(string(bytes))
}

func StringP(in string) *string {
	return &in
}

func IntP(in int) *int {
	return &in
}

func BoolP(in bool) *bool {
	return &in
}

func Int32P(in int32) *int32 {
	return &in
}

func Uint64P(in int64) *uint64 {
	u := uint64(in)
	return &u
}

func Int64P(in int64) *int64 {
	u := int64(in)
	return &u
}

func DedupStringSlice(src []string) []string {
	m := make(map[string]bool)
	res := make([]string, 0)

	for i := range src {
		if _, exist := m[src[i]]; !exist {
			m[src[i]] = true
		}
	}

	for k, _ := range m {
		res = append(res, k)
	}

	return res
}

func JoinStringPtr(src []*string) string {
	tmp := make([]string, 0)

	for i := range src {
		if src[i] != nil && len(*src[i]) > 0 {
			tmp = append(tmp, *src[i])
		}
	}

	return strings.Join(tmp, ",")
}

func ContainsStringSlice(src []string, key string) (bool, int) {
	for i := range src {
		if src[i] == key {
			return true, i
		}
	}

	return false, -1
}

func RemoveSpace(src string) string {
	src = strings.TrimSpace(src)
	src = strings.ReplaceAll(src, " ", "")
	return src
}

func FromStringPointer(in *string, de string) string {
	if in == nil {
		return de
	}

	return *in
}

func FromBoolPointer(in *bool, de bool) bool {
	if in == nil {
		return de
	}

	return *in
}
