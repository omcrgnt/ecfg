package ecfgtool

import "testing"

const benchPkg = "github.com/omcrgnt/ecfg/internal/testdata"

func benchCollectTemplate(b *testing.B, typeName string) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := CollectTemplateEntries(benchPkg, typeName, "BENCH"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCollectTemplate_root5(b *testing.B) {
	benchCollectTemplate(b, "BenchConfig5")
}

func BenchmarkCollectTemplate_root10(b *testing.B) {
	benchCollectTemplate(b, "BenchConfig10")
}

func BenchmarkCollectTemplate_root15(b *testing.B) {
	benchCollectTemplate(b, "BenchConfig15")
}
