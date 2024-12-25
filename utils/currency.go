package utils

const (
	USD  = "USD"
	EUR  = "EUR"
	BDT  = "BDT"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, BDT:
		return true
	}
	return false
}
