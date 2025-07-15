package cache

import (
	"testing"
)

func TestLRUCache_SetAndGet(t *testing.T) {
	cache := NewLRUCache(2)

	out1 := ExecutionOutput{Stdout: "One", Stderr: "", Err: ""}
	out2 := ExecutionOutput{Stdout: "Two", Stderr: "", Err: ""}
	out3 := ExecutionOutput{Stdout: "Three", Stderr: "", Err: ""}

	cache.Set("a", out1)
	cache.Set("b", out2)

	got, err := cache.Get("a")
	if err != nil {
		t.Errorf("expected no error for key 'a', got %v", err)
	}
	if got != out1 {
		t.Errorf("expected %v, got %v", out1, got)
	}

	cache.Set("c", out3)

	_, err = cache.Get("b")
	if err == nil {
		t.Error("expected error when getting evicted key 'b', but got none")
	}

	if _, err := cache.Get("a"); err != nil {
		t.Errorf("expected key 'a' to exist, got error: %v", err)
	}
	if _, err := cache.Get("c"); err != nil {
		t.Errorf("expected key 'c' to exist, got error: %v", err)
	}
}

func TestLRUCache_OverwriteValue(t *testing.T) {
	cache := NewLRUCache(2)
	out1 := ExecutionOutput{Stdout: "Old"}
	out2 := ExecutionOutput{Stdout: "New"}

	cache.Set("x", out1)
	cache.Set("x", out2)

	got, err := cache.Get("x")
	if err != nil {
		t.Errorf("expected to find key 'x', got error: %v", err)
	}
	if got != out2 {
		t.Errorf("expected updated value %v, got %v", out2, got)
	}
}
func TestNewLRUCache_PanicsOnZeroCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when creating zero-capacity cache, but got none")
		}
	}()

	NewLRUCache(0)
}

func TestDoublyLinkedList_BasicOperations(t *testing.T) {
	l := &DoublyLinkedList{}

	n1 := &Node{Key: "1"}
	n2 := &Node{Key: "2"}
	n3 := &Node{Key: "3"}

	l.AddtoFront(n1)
	if l.Head != n1 || l.Tail != n1 {
		t.Error("expected head and tail to be node1 after first insert")
	}

	l.AddtoFront(n2)
	if l.Head != n2 || l.Tail != n1 {
		t.Error("expected node2 as new head")
	}

	l.AddtoFront(n3)
	if l.Head != n3 {
		t.Error("expected node3 as new head")
	}

	// Move n1 to front
	l.MovetoFront(n1)
	if l.Head != n1 {
		t.Error("expected node1 to be moved to front")
	}

	// Remove tail (node2)
	removed := l.RemoveTail()
	if removed != n2 {
		t.Errorf("expected tail to be node2, got %v", removed.Key)
	}
}
