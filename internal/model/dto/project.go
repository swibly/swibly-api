package dto

type ProjectCreation struct {
	Name        string `validate:"required,min=3,max=32"    json:"name"`
	Description string `validate:"omitempty,min=3,max=5000" json:"description"`

	Content any `validate:"omitempty" json:"content"`
	Budget  int `validate:"omitempty" json:"budget"`

	OwnerID uint `json:"-"` // Set using JWT

	Public bool `json:"-"` // Set in an URL query param
}

type ProjectUpdate struct{}

type ProjectAllowList struct {
	Allow struct {
		View    bool `validate:"omitempty"`
		Edit    bool `validate:"omitempty"`
		Delete  bool `validate:"omitempty"`
		Publish bool `validate:"omitempty"`
		Share   bool `validate:"omitempty"`
		Manage  struct {
			Users    bool `validate:"omitempty"`
			Metadata bool `validate:"omitempty"`
		} `validate:"omitempty" json:"manage"`
	} `validate:"require" json:"allow"`
}
