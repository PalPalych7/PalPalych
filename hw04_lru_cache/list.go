package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	ListLen int       // длина списка
	FirstP  *ListItem // первый элемент списка
	LastP   *ListItem // последний элемент списка
}

func (l *list) Len() int {
	return l.ListLen
}

func (l *list) Front() *ListItem {
	return l.FirstP
}

func (l *list) Back() *ListItem {
	return l.LastP
}

func (l *list) PushFront(v interface{}) *ListItem {
	LI := new(ListItem)

	if l.FirstP == nil {
		l.LastP = LI
	} else {
		LI.Next = l.FirstP
		LI2 := l.FirstP
		LI2.Prev = LI
	}
	l.FirstP = LI
	l.ListLen++
	LI.Value = v
	return LI
}

func (l *list) PushBack(v interface{}) *ListItem {
	LI := new(ListItem)
	if l.LastP == nil {
		l.FirstP = LI
	} else {
		LI.Prev = l.LastP
		LI2 := l.LastP
		LI2.Next = LI
	}
	l.LastP = LI
	l.ListLen++
	LI.Value = v
	return LI
}

func (l *list) Remove(i *ListItem) {
	switch {
	case i.Prev == nil: // первый
		if i.Next != nil {
			i.Next.Prev = nil
		} else {
			l.LastP = i.Prev
		}
		l.FirstP = i.Next
	case i.Next == nil: // последний
		if i.Prev != nil {
			i.Prev.Next = nil
		}
		l.LastP = i.Prev
		//		}
	default:
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	l.ListLen--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev != nil {
		iVal := i.Value
		l.Remove(i)
		l.PushFront(iVal)
	}
}

func NewList() List {
	return new(list)
}
