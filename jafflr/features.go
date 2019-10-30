package main

const (
	PickV2 = "pick.v2"
)

func HasFeature(customer string, feature string) bool {
	switch feature {
	case PickV2:
		return true
	default:
		panic("unknown feature")
	}
}
