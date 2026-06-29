package config_test

import (
	"fmt"
	"testing"

	"github.com/omcrgnt/ecfg/config"
	"github.com/omcrgnt/ecfg/internal/testdata"
)

const benchPrefix = "BENCH"

func setBenchEnv(b *testing.B, rootBlocks int) {
	b.Helper()
	for i := 1; i <= rootBlocks; i++ {
		block := fmt.Sprintf("B%02d", i)
		for _, leaf := range []string{"F1", "F2", "F3", "F4", "F5"} {
			key := fmt.Sprintf("%s_%s_%s", benchPrefix, block, leaf)
			b.Setenv(key, "v")
		}
	}
}

func benchParse[T any](b *testing.B) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := config.Parse[T](config.WithPrefix(benchPrefix)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParse_root5(b *testing.B) {
	setBenchEnv(b, 5)
	benchParse[testdata.BenchConfig5](b)
}

func BenchmarkParse_root10(b *testing.B) {
	setBenchEnv(b, 10)
	benchParse[testdata.BenchConfig10](b)
}

func BenchmarkParse_root15(b *testing.B) {
	setBenchEnv(b, 15)
	benchParse[testdata.BenchConfig15](b)
}
