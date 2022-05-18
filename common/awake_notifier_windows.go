package common

type AwakeNotifierWindows struct{}

func (AwakeNotifierWindows) Start(chan bool) error {
	return nil
}

func NewAwakeNotifier() (AwakeNotifier, error) {
	return AwakeNotifierWindows{}, nil
}
