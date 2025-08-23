package main

type Node struct {
	value string
	prev  *Node
	next  *Node
}

type DoublyLinkedList struct {
	head *Node
	tail *Node
	len  int
}

func newNode(value string) *Node {
	return &Node{value: value}
}

func NewDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{
		head: nil,
		tail: nil,
		len:  0,
	}
}

func (l *DoublyLinkedList) LPush(value string) {
	node := newNode(value)

	if l.len == 0 {
		l.head, l.tail = node, node
	} else {
		node.next = l.head
		l.head.prev = node
		l.head = node
	}
	l.len++
}

func (l *DoublyLinkedList) RPush(value string) {
	node := newNode(value)
	if l.len == 0 {
		l.head, l.tail = node, node
	} else {
		l.tail.next = node
		node.prev = l.tail
		l.tail = node
	}
	l.len++
}

func (l *DoublyLinkedList) LPop() (string, bool) {

	switch {
	case l.len == 1:
		result := l.head.value
		l.head, l.tail = nil, nil
		l.len--
		return result, true
	case l.len > 1:
		result := l.head.value
		l.head = l.head.next
		l.head.prev = nil
		l.len--
		return result, true
	default:
		return "", false
	}

}

func (l *DoublyLinkedList) RPop() (string, bool) {
	switch {
	case l.len == 1:
		result := l.head.value
		l.head, l.tail = nil, nil
		l.len--
		return result, true
	case l.len > 1:
		result := l.tail.value
		l.tail = l.tail.prev
		l.tail.next = nil
		return result, true
	default:
		return "", false
	}
}

func (l *DoublyLinkedList) LRange(start, end int) []string {

	if start < 0 {
		if Abs(start) > l.len {
			start = 0
		} else {
			start = l.len + start
		}
	}

	if end < 0 {
		if Abs(start) > l.len {
			start = 0
		} else {
			end = l.len + end
		}

	}

	if start > end || start >= l.len {
		return []string{}
	}

	if end > l.len {
		end = l.len
	}

	result := make([]string, 0)

	current := l.head
	currInd := 0

	for current != nil {
		if currInd >= start && currInd <= end {
			result = append(result, current.value)
		}
		if currInd == end {
			break
		}
		currInd++
		current = current.next
	}

	return result
}

func (l *DoublyLinkedList) Len() int {
	return l.len
}
