package game

type Answer struct {
	AfterReleaseYear  *int `json:"afterReleaseYear"`
	BeforeReleaseYear *int `json:"beforeReleaseYear"`
	GuessIndex        *int `json:"guessIndex"`
}
