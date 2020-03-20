package core

func GetLatestAccounts(idx int, size int) []AccountRecord {
	if idx == 0 {
		idx = 1
	}
	if size == 0 {
		size = 50
	}

	return List.QueryLatestAccounts(idx, size)
}
