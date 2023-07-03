package state

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/liminin/goscene/model"
)

type Item struct {
	m   *model.Item
	err error
}

func NewItem(m *model.Item, err error) *Item {
	return &Item{
		m:   m,
		err: err,
	}
}

func (i *Item) Result() (string, error) {
	return i.Val(), i.err
}

func (i *Item) Val() string {
	return fmt.Sprint(i.m.Value)
}

func (i *Item) Err() error {
	return i.err
}

func (i *Item) Int() (int, error) {
	if i.err != nil {
		return 0, i.err
	}
	return strconv.Atoi(i.Val())
}

func (i *Item) Uint64() (uint64, error) {
	if i.err != nil {
		return 0, i.err
	}
	return strconv.ParseUint(i.Val(), 10, 64)
}

func (i *Item) Float64() (float64, error) {
	if i.err != nil {
		return 0, i.err
	}
	return strconv.ParseFloat(i.Val(), 64)
}

func (i *Item) Bool() (bool, error) {
	if i.err != nil {
		return false, i.err
	}
	return strconv.ParseBool(i.Val())
}

func (i *Item) Bytes() ([]byte, error) {
	return []byte(i.Val()), i.Err()
}

func (i *Item) ScanJSON(v any) {
	data, _ := i.Bytes()
	json.Unmarshal(data, v)
}
