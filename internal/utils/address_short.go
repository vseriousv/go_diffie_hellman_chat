package utils

import "fmt"

func AddressShort(address string) string {
	return fmt.Sprintf("%s...%s", address[:7], address[len(address)-7:])
}
