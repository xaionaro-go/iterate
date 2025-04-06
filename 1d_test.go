package iterate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegmentsDeduplicate(t *testing.T) {
	assert.Equal(t, []Segment[[]struct{}]{
		{S: 10, E: 20, Data: []struct{}{{}}},
		{S: 30, E: 40, Data: []struct{}{{}}},
	}, Segments[struct{}]{
		{S: 10, E: 20},
		{S: 30, E: 40},
	}.Deduplicate(nil))
}

func BenchmarkSegments(b *testing.B) {
	for srcCount := 1; srcCount <= 1024; srcCount *= 2 {
		b.Run(fmt.Sprintf("srcCount%d", srcCount), func(b *testing.B) {
			segments := make(Segments[struct{}], 0, srcCount)
			for i := range srcCount {
				segments = append(segments, Segment[struct{}]{
					S: i * 10,
					E: (i + 1) * 10,
				})
			}
			b.Run("Sort", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					segments.Sort()
				}
			})
			b.Run("Deduplicate", func(b *testing.B) {
				b.Run("nil-arg", func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						segments.Deduplicate(nil)
					}
				})
				var dedupSegments Segments[[]struct{}]
				b.Run("reuse-storage", func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						dedupSegments = segments.Deduplicate(dedupSegments)
					}
				})
			})
		})
	}
}
