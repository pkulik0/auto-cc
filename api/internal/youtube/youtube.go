package youtube

type Youtube interface{}

type youtube struct{}

func New() *youtube {
	return &youtube{}
}
