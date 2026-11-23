package users

import (
	"context"
	"log/slog"
)

func (r *repo) CountUser(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	count, err := r.col.EstimatedDocumentCount(ctx, nil)
	if err != nil {

		slog.ErrorContext(ctx, "failed to count users", slog.String("error", err.Error()))
		return 0, err
	}

	return count, nil
}
