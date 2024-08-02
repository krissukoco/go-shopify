package shopify

import (
	"fmt"
	"strings"
)

type Shop struct {
	// base is the shortest identifier of the shop
	// e.g.: example-store-1
	// which can be converted later to FullID like 'example-store-1.myshopify.com'
	// or BaseURL like 'https://example-store-1.myshopify.com'
	base string
}

func ParseShop(s string) Shop {
	// Strip 'https://' prefix and '.myshopify.com' suffix
	s = strings.Replace(s, "https://", "", -1)
	s = strings.Replace(s, "http://", "", -1)
	s = strings.Replace(s, ".myshopify.com", "", -1)
	return Shop{s}
}

func (s Shop) FullID() string {
	return fmt.Sprintf("%s.myshopify.com", s.base)
}

func (s Shop) BaseUrl() string {
	return fmt.Sprintf("https://%s", s.FullID())
}

func (s Shop) BaseID() string {
	return s.base
}
