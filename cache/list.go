package cache

type Node struct {
	Key   string
	Value ExecutionOutput
	Prev  *Node
	Next  *Node
}

type DoublyLinkedList struct {
	Head *Node
	Tail *Node
}

func (l *DoublyLinkedList) AddtoFront(node *Node) {
	node.Prev = nil
	node.Next = l.Head

	if l.Head != nil {
		l.Head.Prev = node
	} else {
		l.Tail = node
	}

	l.Head = node
}

func (l *DoublyLinkedList) MovetoFront(node *Node) {
	l.Remove(node)
	l.AddtoFront(node)
}

func (l *DoublyLinkedList) Remove(node *Node) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		l.Head = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		l.Tail = node.Prev
	}

	node.Next = nil
	node.Prev = nil
}

func (l *DoublyLinkedList) RemoveTail() *Node {
	if l.Tail == nil {
		return nil
	}

	oldTail := l.Tail
	l.Remove(l.Tail)
	return oldTail
}
