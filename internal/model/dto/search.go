package dto

type search struct {
	Name *string `json:"name"`

	OrderAscending bool `json:"ascending"`

	OrderAlphabetic bool `json:"order_alphabetic"`

	FollowedUsersOnly bool `json:"followed"`

	OrderCreationDate bool `json:"order_created"`
	OrderModifiedDate bool `json:"order_modified"`
}

type SearchUser struct {
	search

	VerifiedOnly bool `json:"verified_only"`

	MostFollowers bool `json:"most_followers"`
}

type SearchProject struct {
	search

	MostFavorites bool `json:"most_favorites"`
	MostClones    bool `json:"most_clones"`

	MinArea int `json:"min_area"`
	MaxArea int `json:"max_area"`

	MinBudget int `json:"min_budget"`
	MaxBudget int `json:"max_budget"`
}
