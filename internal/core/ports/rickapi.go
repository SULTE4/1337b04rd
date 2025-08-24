package ports

type RickAPI interface {
	// Returns (name, imageURL, totalCharacters, err)
	RandomCharacter() (string, string, int, error)
}
