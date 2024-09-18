package notifier

type Notification struct {
	FromUserID uint
	Key        string
	Payload    interface{}
}
