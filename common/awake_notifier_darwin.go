package common

import "github.com/prashantgupta24/mac-sleep-notifier/notifier"

type AwakeNotifierDarwn struct {
	notifier *notifier.Notifer
}

func NewAwakeNotifier() (*AwakeNotifierDarwn, error) {
	n := notifier.GetInstance()

	return AwakeNotifierDarwn{notifier: n}, nil
}

func (a *AwakeNotifierDarwn) Start(out chan bool) error {
	suspendResumeNotifier := a.Start()

	go func() {
		for {
			select {
			case activity := <-suspendResumeNotifier:
				if activity.Type == notifier.Awake {
					out <- true
				}
			}
		}
	}()

	return out, nil
}
