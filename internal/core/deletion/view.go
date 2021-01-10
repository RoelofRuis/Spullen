package deletion

import (
	"github.com/roelofruis/spullen/internal/core"
	"github.com/roelofruis/spullen/internal/core/object"
)

type Delete struct {
	Alert string

	Original *object.Form
	Form     *core.DeleteForm
}
