package notifier

type Notification struct {
	FromUserID uint64
	Key        string
	Payload    interface{}
}
