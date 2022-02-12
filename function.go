package blog1

import (
	"context"
	"fmt"

	"cloud.google.com/go/functions/metadata"
	"github.com/suzuito/blog1-go/deployment/gcf"
)

func BlogUpdateArticle(ctx context.Context, ev gcf.GCSEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	return gcf.BlogUpdateArticle(ctx, meta, ev)
}
