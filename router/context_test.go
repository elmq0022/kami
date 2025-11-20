package router_test

import (
	"context"
	"maps"
	"testing"

	"github.com/elmq0022/kami/router"
)

func TestParamsRoundTrip(t *testing.T) {
	want := map[string]string{"foo": "bar"}
	ctx := router.WithParams(context.Background(), want)

	got := router.GetParams(ctx)
	if !maps.Equal(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}

	empty := router.GetParams(context.Background())
	if len(empty) != 0 {
		t.Fatalf("expected empty map, got %v", empty)
	}
}
