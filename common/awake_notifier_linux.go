package common

type AwakeNotifierEmpty struct{}

func (AwakeNotifierEmpty) Start(chan bool) error {
	return nil
}

func NewAwakeNotifier() (AwakeNotifier, error) {
	return AwakeNotifierEmpty{}, nil
}
