package external

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"1337b04rd/internal/core/ports"
)

type rickAPI struct {
	client *http.Client
	base   string
}

func NewRickAPI() ports.RickAPI {
	return &rickAPI{
		client: &http.Client{Timeout: 5 * time.Second},
		base:   "https://rickandmortyapi.com/api",
	}
}

type charactersResp struct {
	Info struct {
		Count int `json:"count"`
	} `json:"info"`
}

type character struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (r *rickAPI) RandomCharacter() (string, string, int, error) {
	// 1) get count
	req, _ := http.NewRequest(http.MethodGet, r.base+"/character", nil)
	resp, err := r.client.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()
	var cr charactersResp
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return "", "", 0, err
	}
	if cr.Info.Count == 0 {
		return "", "", 0, fmt.Errorf("no characters")
	}

	// 2) pick random id (1..count)
	rand.Seed(time.Now().UnixNano())
	id := 1 + rand.Intn(cr.Info.Count)

	// 3) fetch character
	req2, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/character/%d", r.base, id), nil)
	resp2, err := r.client.Do(req2)
	if err != nil {
		return "", "", cr.Info.Count, err
	}
	defer resp2.Body.Close()
	var c character
	if err := json.NewDecoder(resp2.Body).Decode(&c); err != nil {
		return "", "", cr.Info.Count, err
	}
	return c.Name, c.Image, cr.Info.Count, nil
}
