package common

import "github.com/prashantgupta24/mac-sleep-notifier/notifier"

type AwakeNotifierDarwn struct {
	n *notifier.Notifier
}

func NewAwakeNotifier() (*AwakeNotifierDarwn, error) {
	n := notifier.GetInstance()

	return &AwakeNotifierDarwn{n: n}, nil
}

func (a *AwakeNotifierDarwn) Start(out chan bool) error {
	suspendResumeNotifier := a.n.Start()

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

	return nil
}
