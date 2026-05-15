package walker_test

import (
	"testing"

	"github.com/omcrgnt/ecfg/pkg/walker"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestWalk_protoFieldIsLeaf(t *testing.T) {
	type Config struct {
		TS   *timestamppb.Timestamp
		Name string
	}

	cfg := Config{Name: "app", TS: &timestamppb.Timestamp{}}
	p, err := walker.NewReflectProvider(&cfg)
	require.NoError(t, err)

	var visited []string
	w := walker.New()
	err = w.Walk(p, func(f walker.Field) error {
		_, sf, err := f.Value()
		if err != nil {
			return err
		}
		if sf.Name != "" {
			visited = append(visited, sf.Name)
		}
		return nil
	})
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"TS", "Name"}, visited)
}
