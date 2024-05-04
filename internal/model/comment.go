package model

type Comment struct {
	ID       uint   `gorm:"primaryKey"`
	OwnerID  uint   `gorm:"index"`
	Message  string `gorm:"type:text"` // Type `text`, as `text` is bigger than the default for `string`
	Likes    uint
	Dislikes uint
}
