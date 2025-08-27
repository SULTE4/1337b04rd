package ports

import "1337b04rd/internal/core/domain"

type ExternalApi interface {
	GetRandomCharacter([]int) (domain.Character, error)
}
