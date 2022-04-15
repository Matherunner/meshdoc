package meshdoc

import (
	"github.com/Matherunner/meshdoc/context"
)

type Postprocessor interface {
	Process(ctx context.Context, r ParsedReader) (ParsedReader, error)
}
