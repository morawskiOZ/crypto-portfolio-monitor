package price

import "strings"

// ConvertToUSDTPair converts coin symbol to it's USDT pair for valuation
func ConvertToUSDTPair(s string) (pair string, ok bool) {
	if strings.Contains(s, "USD") {
		return s, false
	}

	return strings.Join([]string{s, "USDT"}, ""), true
}
