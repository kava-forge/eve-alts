package bindings

import (
	"fyne.io/fyne/v2/data/binding"
)

type DataProxy[T any] interface {
	AddListener(l binding.DataListener)
	RemoveListener(l binding.DataListener)
	Get() (val T, err error)
	Set(val T) error
}

var (
	_ DataProxy[int] = (*DataStructProxy[int])(nil)
	_ DataProxy[int] = (*DataStructMapProxy[int])(nil)
)

type DataStruct[T any] struct {
	binding.Untyped
}

func NewDataStruct[T any]() *DataStruct[T] {
	return &DataStruct[T]{
		Untyped: binding.NewUntyped(),
	}
}

func (d *DataStruct[T]) Get() (val T, err error) {
	raw, err := d.Untyped.Get()
	if err != nil {
		return val, err
	}

	return raw.(T), nil
}

func (d *DataStruct[T]) Set(val T) error {
	return d.Untyped.Set(val)
}

type DataList[T any] struct {
	binding.UntypedList
}

func NewDataList[T any]() *DataList[T] {
	raw := binding.NewUntypedList()
	return &DataList[T]{UntypedList: raw}
}

func (cl *DataList[T]) Append(value T) error  { return cl.UntypedList.Append(value) }
func (cl *DataList[T]) Prepend(value T) error { return cl.UntypedList.Prepend(value) }

func (cl *DataList[T]) SetValue(index int, value T) error {
	return cl.UntypedList.SetValue(index, value)
}

func (cl *DataList[T]) Set(list []T) error {
	l := make([]interface{}, 0, len(list))
	for _, it := range list {
		l = append(l, it)
	}
	return cl.UntypedList.Set(l)
}

func (cl *DataList[T]) Get() ([]T, error) {
	list, err := cl.UntypedList.Get()
	if err != nil {
		return nil, err
	}

	l := make([]T, 0, len(list))
	for _, it := range list {
		l = append(l, it.(T))
	}

	return l, nil
}

func (cl *DataList[T]) GetValue(index int) (t T, err error) {
	val, err := cl.UntypedList.GetValue(index)
	if err != nil {
		return t, err
	}

	return val.(T), nil
}

func (cl *DataList[T]) Child(index int) *DataStructProxy[T] {
	return NewDataStructProxy(cl, index)
}

type DataStructProxy[T any] struct {
	parent *DataList[T]
	idx    int
}

func NewDataStructProxy[T any](parent *DataList[T], idx int) *DataStructProxy[T] {
	return &DataStructProxy[T]{
		parent: parent,
		idx:    idx,
	}
}

func (d *DataStructProxy[T]) AddListener(l binding.DataListener) {
	item, err := d.parent.GetItem(d.idx)
	if err != nil {
		return
	}

	item.AddListener(l)
}

func (d *DataStructProxy[T]) RemoveListener(l binding.DataListener) {
	item, err := d.parent.GetItem(d.idx)
	if err != nil {
		return
	}

	item.RemoveListener(l)
}

func (d *DataStructProxy[T]) Get() (val T, err error) {
	raw, err := d.parent.GetValue(d.idx)
	if err != nil {
		return val, err
	}

	return raw, nil
}

func (d *DataStructProxy[T]) Set(val T) error {
	return d.parent.SetValue(d.idx, val)
}

type DataMap[V any] struct {
	binding.UntypedMap
}

func NewDataMap[V any]() *DataMap[V] {
	raw := binding.NewUntypedMap()
	return &DataMap[V]{UntypedMap: raw}
}

func (cl *DataMap[V]) HasKey(k string) bool {
	_, err := cl.UntypedMap.GetValue(k)
	return err == nil
}

func (cl *DataMap[V]) Delete(k string) {
	cl.UntypedMap.Delete(k)
}

func (cl *DataMap[V]) SetValue(k string, value V) error {
	return cl.UntypedMap.SetValue(k, value)
}

func (cl *DataMap[V]) Set(data map[string]V) error {
	l := make(map[string]interface{}, len(data))
	for k, v := range data {
		l[k] = v
	}
	return cl.UntypedMap.Set(l)
}

func (cl *DataMap[V]) Get() (map[string]V, error) {
	raw, err := cl.UntypedMap.Get()
	if err != nil {
		return nil, err
	}

	l := make(map[string]V, len(raw))
	for k, it := range raw {
		l[k] = it.(V)
	}

	return l, nil
}

func (cl *DataMap[V]) GetValue(k string) (t V, err error) {
	val, err := cl.UntypedMap.GetValue(k)
	if err != nil {
		return t, err
	}

	return val.(V), nil
}

func (cl *DataMap[V]) Child(k string) *DataStructMapProxy[V] {
	return NewDataStructMapProxy(cl, k)
}

type DataStructMapProxy[T any] struct {
	parent *DataMap[T]
	idx    string
}

func NewDataStructMapProxy[T any](parent *DataMap[T], idx string) *DataStructMapProxy[T] {
	return &DataStructMapProxy[T]{
		parent: parent,
		idx:    idx,
	}
}

func (d *DataStructMapProxy[T]) AddListener(l binding.DataListener) {
	item, err := d.parent.GetItem(d.idx)
	if err != nil {
		return
	}

	item.AddListener(l)
}

func (d *DataStructMapProxy[T]) RemoveListener(l binding.DataListener) {
	item, err := d.parent.GetItem(d.idx)
	if err != nil {
		return
	}

	item.RemoveListener(l)
}

func (d *DataStructMapProxy[T]) Get() (val T, err error) {
	raw, err := d.parent.GetValue(d.idx)
	if err != nil {
		return val, err
	}

	return raw, nil
}

func (d *DataStructMapProxy[T]) Set(val T) error {
	return d.parent.SetValue(d.idx, val)
}
