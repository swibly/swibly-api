package dto

type APIKey struct {
	EnabledKeyManage   int `validate:"omitempty,mustbenumericalboolean" json:"enabled_key_manage"`
	EnabledAuth        int `validate:"omitempty,mustbenumericalboolean" json:"enabled_auth"`
	EnabledSearch      int `validate:"omitempty,mustbenumericalboolean" json:"enabled_search"`
	EnabledUserFetch   int `validate:"omitempty,mustbenumericalboolean" json:"enabled_user_fetch"`
	EnabledUserActions int `validate:"omitempty,mustbenumericalboolean" json:"enabled_user_actions"`
}
