package touch

type MouseType int32

const (
	MouseTypeUnknown MouseType = iota
	MouseTypeDown
	MouseTypeMove
	MouseTypeUp
)
