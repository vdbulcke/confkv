package controller

import "context"

func (c *Controller) Get(ctx context.Context, bucket, key string) ([]byte, error) {

	return c.store.Get(ctx, bucket, key)
}
