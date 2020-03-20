package core

func GetAmountOutsV3(addr string, amount uint64, fees uint64) map[string]interface{} {
	var sum uint64
	rs := make(map[string]interface{})
	outs, fAmount := QueryOutput.GetUsable(addr)
	rs["info"] = 0
	rs["outs"] = outs
	iAmount, _ := NewAmount(fAmount)
	if iAmount < amount {
		rs["info"] = 2
		return rs
	}
	amount += fees
	length := 0
	for i, out := range outs {
		sum += out.Amount
		if sum >= amount {
			rs["outs"] = outs[0 : i+1]
			length = len(outs[0 : i+1])
			break
		}
	}
	if length > 20 {
		rs["info"] = 1
	}
	return rs
}
