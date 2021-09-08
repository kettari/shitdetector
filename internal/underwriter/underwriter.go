package underwriter

import (
	"fmt"
	"github.com/kettari/shitdetector/internal/asset"
)

type (
	Underwriter interface {
		Score() *Score
	}
	Score struct {
		Stock *asset.Stock
		Criterias []*ScoreCriteria
	}
	ScoreCriteria struct {
		Name  string
		Value int64
	}
)

func (s Score) TotalScore() (total int64) {
	total = 0
	for _, v := range s.Criterias {
		total += v.Value
	}
	return total
}

func (s Score) MaxScore() int64 {
	return int64(len(s.Criterias) * 5)
}

func (s Score) Describe() (desc string) {
	desc = fmt.Sprintf(`<b>Скоринг по методу <a href="https://t.me/Finindie/767">Finindie</a> [%d/%d]</b>`, s.TotalScore(), s.MaxScore())
	position := 1
	for _, v := range s.Criterias {
		desc += fmt.Sprintf("\n%d) %s - %d", position, v.Name, v.Value)
		position++
	}
	return desc
}
