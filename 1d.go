package iterate

import (
	"github.com/go-ng/sort"
)

type Segment[T any] struct {
	S    int // start
	E    int // end
	Data T
}

type Segments[T any] []Segment[T]

func (s Segments[T]) Len() int {
	return len(s)
}

func (s Segments[T]) Less(i, j int) bool {
	return s[i].S < s[j].S
}

func (s Segments[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Segments[T]) Sort() {
	sort.Sort(s)
}

// Note: This function requires segments to be already sorted (see method `Sort`).
func (s Segments[T]) Deduplicate(storage []Segment[[]T]) []Segment[[]T] {
	storage = storage[:0]
	if len(s) == 0 {
		return storage
	}
	var data []T
	if cap(storage) > 0 {
		data = storage[:1][0].Data
		data = data[:0]
	}
	data = append(data, s[0].Data)
	curItem := Segment[[]T]{
		S:    s[0].S,
		E:    s[0].E,
		Data: data,
	}
	for _, segment := range s[1:] {
		if segment.S <= curItem.E+1 {
			if segment.E > curItem.E {
				curItem.E = segment.E
			}
			curItem.Data = append(curItem.Data, segment.Data)
			continue
		}
		storage = append(storage, curItem)
		var data []T
		if cap(storage) > len(storage) {
			data = storage[:len(storage)+1][len(storage)].Data
			data = data[:0]
		}
		data = append(data, segment.Data)
		curItem = Segment[[]T]{
			S:    segment.S,
			E:    segment.E,
			Data: data,
		}
	}
	storage = append(storage, curItem)
	return storage
}
