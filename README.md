# `iterate`

This package is focused on providing generic iterators. Currently, we have only one iterator:
* `TwoD` for iterating efficiently over every point on a big 2D plane.

For example let's say you want to draw two overlapping rectangles, but you need to do it in a performance-optimal way. A naive solution would be to just draw then as a normal human being:
```go
draw.Draw(canvas, image.Rect(100, 100, 300, 300), &image.Uniform{colorFG}, image.Point{}, draw.Src)
draw.Draw(canvas, image.Rect(200, 200, 400, 400), &image.Uniform{colorFG}, image.Point{}, draw.Src)
```
But here some computation power is lost due to walking twice in area (20,20)-(30,30) and other stuff may be done suboptimal. To avoid these issues you may use a "smart" 2D iterator from this package:
```go
for i := 0; i < b.N; i++ {
	for p, _ := range TwoD(
		nil,
		Rect[struct{}](100, 100, 300, 300),
		Rect[struct{}](200, 200, 400, 400),
	) {
		canvasB.SetGray(p.X, p.Y, colorFG)
	}
}
```
In result, you'll get a performance gain:
```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/iterate
cpu: AMD Ryzen 9 5900X 12-Core Processor
                        │    sec/op     │    sec/op      vs base                │
TwoD/sizeFactor1/-24       161.0n ± ∞ ¹    808.8n ± ∞ ¹  +402.36% (p=0.008 n=5)
TwoD/sizeFactor2/-24       330.5n ± ∞ ¹   1022.0n ± ∞ ¹  +209.23% (p=0.008 n=5)
TwoD/sizeFactor4/-24       939.4n ± ∞ ¹   1617.0n ± ∞ ¹   +72.13% (p=0.008 n=5)
TwoD/sizeFactor8/-24       3.401µ ± ∞ ¹    3.148µ ± ∞ ¹    -7.44% (p=0.008 n=5)
TwoD/sizeFactor16/-24     13.167µ ± ∞ ¹    7.590µ ± ∞ ¹   -42.36% (p=0.008 n=5)
TwoD/sizeFactor32/-24      54.62µ ± ∞ ¹    21.87µ ± ∞ ¹   -59.96% (p=0.008 n=5)
TwoD/sizeFactor64/-24     215.02µ ± ∞ ¹    76.75µ ± ∞ ¹   -64.31% (p=0.008 n=5)
TwoD/sizeFactor128/-24     842.1µ ± ∞ ¹    280.3µ ± ∞ ¹   -66.71% (p=0.008 n=5)
TwoD/sizeFactor256/-24     3.340m ± ∞ ¹    1.069m ± ∞ ¹   -68.00% (p=0.008 n=5)
TwoD/sizeFactor512/-24    13.462m ± ∞ ¹    4.214m ± ∞ ¹   -68.69% (p=0.008 n=5)
TwoD/sizeFactor1024/-24    54.25m ± ∞ ¹    16.72m ± ∞ ¹   -69.17% (p=0.008 n=5)
```

The implementation of the iteration also could be improved, but yet it is already significantly better than the original variant: 3 times faster if the picture is big enough.
