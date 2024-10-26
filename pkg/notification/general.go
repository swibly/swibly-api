package notification

type NotificationType string

const (
	Info    NotificationType = "info"
	Warning NotificationType = "warning"
	Error   NotificationType = "error"
	Danger  NotificationType = "danger"
)

var (
	Array       = []NotificationType{Info, Warning, Error, Danger}
	ArrayString = []string{string(Info), string(Warning), string(Error), string(Danger)}
)
