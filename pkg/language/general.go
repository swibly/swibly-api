package language

type Language string

// NOTE: Some languages might not be available at launch of the API, like Russian or Portuguese (which are planned at the moment)

// Officially supported
const (
	PT Language = "pt" // Portuguese
	EN Language = "en" // English
	RU Language = "ru" // Russian
)

// Some utilities

var (
	Array       = []Language{PT, EN, RU}
	ArrayString = []string{string(PT), string(EN), string(RU)}
)
