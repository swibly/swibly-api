package notification

type NotificationType string

const (
	Information NotificationType = "information"
	Warning     NotificationType = "warning"
	Danger      NotificationType = "danger"
)

var (
	Array       = []NotificationType{Information, Warning, Danger}
	ArrayString = []string{string(Information), string(Warning), string(Danger)}
)
