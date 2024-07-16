package app

import (
	"cmp"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MiniTagSetItem[S cmp.Ordered] interface {
	fyne.CanvasObject
	ShouldShow() bool
	SortKey() S
}

type MiniTagSet[S cmp.Ordered, T MiniTagSetItem[S]] struct {
	widget.BaseWidget

	items   []T
	minSize fyne.Size
}

func NewMiniTagSet[S cmp.Ordered, T MiniTagSetItem[S]]() *MiniTagSet[S, T] {
	ts := &MiniTagSet[S, T]{}
	ts.ExtendBaseWidget(ts)
	return ts
}

func (c *MiniTagSet[S, T]) Len() int {
	return len(c.items)
}

func (c *MiniTagSet[S, T]) Less(i, j int) bool {
	return c.items[i].SortKey() < c.items[j].SortKey()
}

func (c *MiniTagSet[S, T]) Swap(i, j int) {
	c.items[i], c.items[j] = c.items[j], c.items[i]
}

func (c *MiniTagSet[S, T]) Add(v T) {
	c.items = append(c.items, v)
	slices.SortStableFunc(c.items, func(a, b T) int {
		return cmp.Compare(a.SortKey(), b.SortKey())
	})
	v.Resize(v.MinSize())
}

func (c *MiniTagSet[S, T]) CreateRenderer() fyne.WidgetRenderer {
	return c
}

// The WidgetRenderer interface

func (c *MiniTagSet[S, T]) Destroy() {}

func max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func (c *MiniTagSet[S, T]) Layout(sz fyne.Size) {
	var x, y float32 = 0, 0
	var lineht float32 = 0
	for _, tsi := range c.items {
		if !tsi.ShouldShow() {
			tsi.Hide()
			continue
		}

		tsi.Show()
		tsi.Refresh()
		remain := sz.Width - x

		mtsz := tsi.MinSize()
		lineht = max(lineht, mtsz.Height+theme.Padding()/2)
		tsi.Resize(mtsz)
		if mtsz.Width > remain {
			x = 0
			y += lineht
		}

		// place
		tsi.Move(fyne.Position{X: x, Y: y})

		x += mtsz.Width + theme.Padding()
	}
}

func (c *MiniTagSet[S, T]) MinSize() fyne.Size {
	maxsz := c.minSize
	for _, mt := range c.items {
		mtsz := mt.MinSize()
		maxsz.Height = max(maxsz.Height, mtsz.Height)
		maxsz.Width = max(maxsz.Width, mtsz.Width)
	}

	return maxsz
}

func (c *MiniTagSet[S, T]) Objects() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 0, len(c.items))
	for _, mt := range c.items {
		objs = append(objs, mt)
	}

	return objs
}

func (c *MiniTagSet[S, T]) Refresh() {
	for _, o := range c.Objects() {
		o.Refresh()
	}
}
