package external

import (
	"1337b04rd/internal/core/domain"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
)

type RickAndMortyApiRepo struct {
	logger *slog.Logger
}

func NewRandMApi(logger slog.Logger) *RickAndMortyApiRepo {
	return &RickAndMortyApiRepo{
		logger: &logger,
	}
}

func (a *RickAndMortyApiRepo) GetRandomCharacter(occupied []int) (domain.Character, error) {
	// get character count
	cnt, err := characterCount()
	if err != nil {
		a.logger.Error(err.Error())
		return domain.Character{}, err
	}
	// get random character id
	var id int
	for {
		id = rand.Intn(cnt) + 1
		if !contains(occupied, id) {
			break
		}
	}

	resp, err := http.Get(fmt.Sprintf("https://rickandmortyapi.com/api/character/%d", id))
	if err != nil {
		a.logger.Error(err.Error())
		return domain.Character{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Error(err.Error())
		return domain.Character{}, err
	}

	var character domain.Character
	err = json.Unmarshal(body, &character)
	if err != nil {
		a.logger.Warn(err.Error())
		return domain.Character{}, err
	}

	return character, nil
}

func characterCount() (int, error) {
	return 826, nil
}

// func main() {
// 	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
// 		AddSource: true,
// 	}))
// 	tmp := NewRandMApi(*logger)
// 	tmp.GetRandomCharacter([]int{})

// }

func contains(slice []int, element int) bool {
	for _, v := range slice {
		if v == element {
			return true // Element found
		}
	}
	return false // Element not found
}
