package notifier

type Notification struct {
	FromUserID uint        `json:"fromUserID"`
	Key        string      `json:"key"`
	Payload    interface{} `json:"payload"`
}
