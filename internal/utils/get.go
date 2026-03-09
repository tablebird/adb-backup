package utils

type Getter interface {
	Get(key any) (value any, exists bool)
}

func GetType[T any](getter Getter, key any) (res T) {
	if val, ok := getter.Get(key); ok && val != nil {
		res, _ = val.(T)
	}
	return
}
