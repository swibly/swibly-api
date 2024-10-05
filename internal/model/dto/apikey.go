package dto

type UpdateAPIKey struct {
	Owner string `validate:"omitempty" json:"owner"`

	EnabledKeyManage   int `validate:"omitempty,mustbenumericalboolean" json:"enabled_key_manage"`
	EnabledAuth        int `validate:"omitempty,mustbenumericalboolean" json:"enabled_auth"`
	EnabledSearch      int `validate:"omitempty,mustbenumericalboolean" json:"enabled_search"`
	EnabledUserFetch   int `validate:"omitempty,mustbenumericalboolean" json:"enabled_user_fetch"`
	EnabledUserActions int `validate:"omitempty,mustbenumericalboolean" json:"enabled_user_actions"`

	TimesUsed uint `validate:"omitempty" json:"times_used"`
	MaxUsage  uint `validate:"omitempty" json:"max_usage"`
}

type ReadAPIKey struct {
	Key   string `json:"key"`
	Owner string `json:"owner"`

	EnabledKeyManage   int `json:"enabled_key_manage"`
	EnabledAuth        int `json:"enabled_auth"`
	EnabledSearch      int `json:"enabled_search"`
	EnabledUserFetch   int `json:"enabled_user_fetch"`
	EnabledUserActions int `json:"enabled_user_actions"`
	EnabledProjects    int `json:"enabled_projects"`

	TimesUsed uint `json:"times_used"`
	MaxUsage  uint `json:"max_usage"`
}
