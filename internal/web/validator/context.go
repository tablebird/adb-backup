package validator

type Context interface {
	Set(key any, value any)

	SetTypeKey(value any)
}
