package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	pushFirstItem(i *ListItem)
}
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}
type list struct {
	First  *ListItem
	Last   *ListItem
	Length int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.First
}

func (l *list) Back() *ListItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v}

	if l.First == nil && l.Last == nil {
		l.pushFirstItem(newListItem)
		return newListItem
	}

	if l.First != nil {
		newListItem.Next = l.First
		l.First.Prev = newListItem
		l.First = newListItem
		l.Length++
	}

	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v, Next: nil, Prev: nil}

	if l.First == nil && l.Last == nil {
		l.pushFirstItem(newListItem)
		return newListItem
	}

	if l.Last != nil {
		newListItem.Prev = l.Last
		l.Last.Next = newListItem
		l.Last = newListItem
		l.Length++
	}

	return newListItem
}

func (l *list) Remove(i *ListItem) {
	// Удалить последний элемент
	if i.Next == nil {
		i.Prev.Next = nil
		l.Last = i.Prev
	}
	// Удалить первый элемент
	if i.Prev == nil {
		i.Next.Prev = nil
		l.First = i.Next
	}
	// Удалить "средний" элемент
	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.Length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}

func (l *list) pushFirstItem(i *ListItem) {
	l.First = i
	l.Last = i
	l.Length++
}
