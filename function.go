package blog1

import (
	"context"

	"github.com/suzuito/blog1-go/deployment/gcf"
)

func BlogUpdateArticle(ctx context.Context, ev gcf.GCSEvent) error {
	return gcf.BlogUpdateArticle(ctx, ev)
}
