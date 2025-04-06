package iterate

import (
	"iter"
)

type Area interface {
	MinY() int
	MaxY() int
	MinX(y int) int
	MaxX(y int) int
}
type Point2D[T any] struct {
	X, Y int
	Data T
}

type Rectangle[T any] struct {
	Min  Point2D[struct{}]
	Max  Point2D[struct{}]
	Data T
}

var _ Area = (*Rectangle[struct{}])(nil)

func Rect[T any](x0, y0, x1, y1 int) Rectangle[T] {
	return Rectangle[T]{
		Min: Point2D[struct{}]{
			X: x0,
			Y: y0,
		},
		Max: Point2D[struct{}]{
			X: x1,
			Y: y1,
		},
	}
}

func (r Rectangle[T]) MinY() int {
	return r.Min.Y
}
func (r Rectangle[T]) MaxY() int {
	return r.Max.Y
}
func (r Rectangle[T]) MinX(int) int {
	return r.Min.X
}
func (r Rectangle[T]) MaxX(int) int {
	return r.Max.X
}

type TwoDReusableBuffers struct {
	YSegments      Segments[Area]
	YSegmentsDedup Segments[[]Area]
	XSegments      Segments[Area]
	XSegmentsDedup Segments[[]Area]
}

func TwoD(
	reuseBuffers *TwoDReusableBuffers,
	areas ...Area,
) iter.Seq[Point2D[[]Area]] {
	return func(yield func(Point2D[[]Area]) bool) {
		for y, xSegments := range TwoDForEachY(reuseBuffers, areas...) {
			for _, xSegment := range xSegments {
				for x := xSegment.S; x < xSegment.E; x++ {
					if !yield(Point2D[[]Area]{
						X:    x,
						Y:    y,
						Data: xSegment.Data,
					}) {
						return
					}
				}
			}
		}
	}
}

func TwoDPointsCount(
	reuseBuffers *TwoDReusableBuffers,
	areas ...Area,
) int {
	totalCount := 0
	for _, xSegments := range TwoDForEachY(reuseBuffers, areas...) {
		for _, xSegment := range xSegments {
			totalCount += xSegment.E - xSegment.S
		}
	}
	return totalCount
}

func TwoDForEachY(
	reuseBuffers *TwoDReusableBuffers,
	areas ...Area,
) iter.Seq2[int, Segments[[]Area]] {
	return func(yield func(int, Segments[[]Area]) bool) {
		if reuseBuffers == nil {
			reuseBuffers = &TwoDReusableBuffers{}
		}

		reuseBuffers.YSegments = reuseBuffers.YSegments[:0]
		for _, area := range areas {
			reuseBuffers.YSegments = append(reuseBuffers.YSegments, Segment[Area]{
				S:    area.MinY(),
				E:    area.MaxY(),
				Data: area,
			})
		}

		reuseBuffers.YSegments.Sort()
		reuseBuffers.YSegmentsDedup = reuseBuffers.YSegments.Deduplicate(reuseBuffers.YSegmentsDedup)

		for _, ySegment := range reuseBuffers.YSegmentsDedup {
			for y := ySegment.S; y < ySegment.E; y++ {
				reuseBuffers.XSegments = reuseBuffers.XSegments[:0]
				for _, area := range ySegment.Data {
					if y < area.MinY() || y >= area.MaxY() {
						continue
					}
					reuseBuffers.XSegments = append(reuseBuffers.XSegments, Segment[Area]{
						S:    area.MinX(y),
						E:    area.MaxX(y),
						Data: area,
					})
				}

				reuseBuffers.XSegments.Sort()
				reuseBuffers.XSegmentsDedup = reuseBuffers.XSegments.Deduplicate(reuseBuffers.XSegmentsDedup)
				if !yield(y, reuseBuffers.XSegmentsDedup) {
					return
				}
			}
		}
	}
}
