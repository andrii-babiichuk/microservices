package repository

import "context"

type Visits interface {
	Inc(ctx context.Context) error
	Get(ctx context.Context) (int, error)
}
