package core

import (
	"github.com/Qitmeer/qitmeer/core/types"
)

func NewAmount(amount float64) (uint64, error) {
	iAmount, err := types.NewAmount(amount)
	if err != nil {
		return 0, err
	}
	return uint64(iAmount), nil
}
