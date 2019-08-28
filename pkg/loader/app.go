package loader

import (
  "context"
)

type App interface {
  Run(context.Context)
}
