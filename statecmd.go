package goscene

import (
	"encoding/json"
	"strconv"
)

type StateCmd struct {
	err error
	val string
}

func NewStateCmd(val string, err error) *StateCmd {
	return &StateCmd{
		val: val,
		err: err,
	}
}

func (cmd *StateCmd) SetVal(val string) {
	cmd.val = val
}

func (cmd *StateCmd) SetErr(err error) {
	cmd.err = err
}

func (cmd *StateCmd) Result() (string, error) {
	return cmd.val, cmd.err
}

func (cmd *StateCmd) Val() string {
	return cmd.val
}

func (cmd *StateCmd) Err() error {
	return cmd.err
}

func (cmd *StateCmd) Int() (int, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.Atoi(cmd.Val())
}

func (cmd *StateCmd) Uint64() (uint64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseUint(cmd.Val(), 10, 64)
}

func (cmd *StateCmd) Float64() (float64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseFloat(cmd.Val(), 64)
}

func (cmd *StateCmd) Bool() (bool, error) {
	if cmd.err != nil {
		return false, cmd.err
	}
	return strconv.ParseBool(cmd.Val())
}

func (cmd *StateCmd) Bytes() ([]byte, error) {
	return []byte(cmd.Val()), cmd.Err()
}

func (cmd *StateCmd) ScanJSON(v any) {
	data, _ := cmd.Bytes()
	json.Unmarshal(data, v)
}
