package iterate

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTwoD(t *testing.T) {
	buf := TwoDReusableBuffers{}
	assert.Equal(t, 300,
		TwoDPointsCount(&buf,
			Rect[struct{}](10, 10, 20, 40), // (10,10)-(20,40): 300
			Rect[struct{}](10, 30, 20, 40), // (10,30)-(20,40): 100 (but overlaps)
		),
	)
	assert.Equal(t, 200,
		TwoDPointsCount(&buf,
			Rect[struct{}](10, 10, 20, 20), // (10,10)-(20,20): 100
			Rect[struct{}](10, 30, 20, 40), // (10,30)-(20,40): 100
		),
	)
	assert.Equal(t, 521528,
		TwoDPointsCount(&buf,
			Rect[struct{}](0, 1557, 556, 1870), // 174028
			Rect[struct{}](0, 779, 556, 1091),  // 173472
			Rect[struct{}](0, 0, 556, 313),     // 174028
		),
	)
}

func BenchmarkTwoD(b *testing.B) {
	colorFG := color.Gray{Y: 255}
	for k := 1; k <= 1024; k *= 2 {
		canvasA := image.NewGray(image.Rect(0, 0, 10*k, 10*k))
		canvasB := image.NewGray(image.Rect(0, 0, 10*k, 10*k))
		b.Run(fmt.Sprintf("sizeFactor%d", k), func(b *testing.B) {
			canvasAReady := false
			b.Run("image.Draw", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					draw.Draw(canvasA, image.Rect(1*k, 1*k, 3*k, 3*k), &image.Uniform{colorFG}, image.Point{}, draw.Src)
					draw.Draw(canvasA, image.Rect(2*k, 2*k, 4*k, 4*k), &image.Uniform{colorFG}, image.Point{}, draw.Src)
				}
				canvasAReady = true
			})

			canvasBReady := false
			b.Run("TwoD", func(b *testing.B) {
				b.Run("simple", func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						for p := range TwoD(
							nil,
							Rect[struct{}](1*k, 1*k, 3*k, 3*k),
							Rect[struct{}](2*k, 2*k, 4*k, 4*k),
						) {
							canvasB.SetGray(p.X, p.Y, colorFG)
						}
					}
					canvasBReady = true
				})
				b.Run("reuse-buffers", func(b *testing.B) {
					var buf TwoDReusableBuffers
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						for p := range TwoD(
							&buf,
							Rect[struct{}](1*k, 1*k, 3*k, 3*k),
							Rect[struct{}](2*k, 2*k, 4*k, 4*k),
						) {
							canvasB.SetGray(p.X, p.Y, colorFG)
						}
					}
					canvasBReady = true
				})
			})

			if canvasAReady && canvasBReady {
				require.Equal(b, canvasA.Pix, canvasB.Pix)
			}

			b.Run("TwoDPointsCount", func(b *testing.B) {
				b.Run("simple", func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						TwoDPointsCount(
							nil,
							Rect[struct{}](1*k, 1*k, 3*k, 3*k),
							Rect[struct{}](2*k, 2*k, 4*k, 4*k),
						)
					}
					canvasBReady = true
				})
				b.Run("reuse-buffers", func(b *testing.B) {
					var buf TwoDReusableBuffers
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						TwoDPointsCount(
							&buf,
							Rect[struct{}](1*k, 1*k, 3*k, 3*k),
							Rect[struct{}](2*k, 2*k, 4*k, 4*k),
						)
					}
					canvasBReady = true
				})
			})
		})
	}
}
