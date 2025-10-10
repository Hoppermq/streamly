package clickhouse

import (
	"log/slog"

	"github.com/hoppermq/streamly/pkg/domain"
)

type Translator struct {
	logger *slog.Logger
}

func (t *Translator) Execute(query *domain.QueryAstRequest) (string, error) {
	return "", nil
}
