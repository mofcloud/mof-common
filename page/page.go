package page

import (
	"encoding/base64"
	"encoding/json"
)

// ************** Page **************

const defaultPageSize = 100

// Page MOF style page
type Page struct {
	PageNum  int `yaml:"pageNum" json:"pageNum"`
	PageSize int `yaml:"pageSize" json:"pageSize"`
}

// Encode Page with base64
func (p *Page) Encode() string {
	bytes, _ := json.Marshal(p)

	return base64.StdEncoding.EncodeToString(bytes)
}

// DecodeToPage decode with base64
func DecodeToPage(str string) (*Page, error) {
	if str == "" {
		return &Page{
			PageNum:  1,
			PageSize: defaultPageSize,
		}, nil
	}

	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	res := &Page{}
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
