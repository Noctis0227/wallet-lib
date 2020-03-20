package core

type IQueryOutput interface {
	GetUsable(addr string) ([]Output, float64)
	GetUnspent(addr string) ([]Output, float64)
	GetSpent(addr string) ([]Output, float64)
	GetMemoryUnspent(addr string) ([]Output, float64)
	GetMemorySpent(addr string) ([]Output, float64)
	GetLockUnspent(addr string) ([]Output, float64)
	GetBalance(addr string) ([]Output, float64)
	GetAllOuts(addr string) ([]Output, float64)
	GetOutput(key string) *Output
}

func GetUsableOuts(addr string) ([]Output, float64) {
	return QueryOutput.GetUsable(addr)
}

func GetUnspentOuts(addr string) ([]Output, float64) {
	return QueryOutput.GetUnspent(addr)
}

func GetSpentOuts(addr string) ([]Output, float64) {
	return QueryOutput.GetSpent(addr)
}

func GetMemoryUnspent(addr string) ([]Output, float64) {
	return QueryOutput.GetMemoryUnspent(addr)
}

func GetMemorySpent(addr string) ([]Output, float64) {
	return QueryOutput.GetMemorySpent(addr)
}

func GetLockedUnspent(addr string) ([]Output, float64) {
	return QueryOutput.GetLockUnspent(addr)
}

func GetBalance(addr string) ([]Output, float64) {
	return QueryOutput.GetBalance(addr)
}

func GetAllOuts(addr string) ([]Output, float64) {
	return QueryOutput.GetAllOuts(addr)
}

func GetOutput(key string) *Output {
	return Storage.GetOutput(key)
}
