package main

const (
	PickV2 = "pick.v2"
)

func HasFeature(customer string, feature string) bool {
	switch feature {
	case PickV2:
		if customer == "gita" {
			return true
		}
		return false
	default:
		panic("unknown feature")
	}
}
