package utils

import (
	"github.com/vistormu/go-dsa/hashmap"
	"github.com/vistormu/go-dsa/queue"
	"github.com/vistormu/go-dsa/set"
	"github.com/vistormu/go-dsa/stack"
)

type QueueArray[T any] = queue.QueueArray[T]
type StackArray[T any] = stack.StackArray[T]
type HashSet[T comparable] = set.HashSet[T]
type UniqueStack[T comparable] = stack.UniqueStack[T]
type BiHashmap[K, V comparable] = hashmap.BiHashmap[K, V]

func NewQueueArray[T any]() *QueueArray[T] {
	return queue.NewQueueArray[T]()
}

func NewStackArray[T any]() *StackArray[T] {
	return stack.NewStackArray[T]()
}

func NewHashSet[T comparable]() *HashSet[T] {
	return set.NewHashSet[T]()
}

func NewUniqueStack[T comparable]() *UniqueStack[T] {
	return stack.NewUniqueStack[T]()
}

func NewBiHashmap[K, V comparable]() *BiHashmap[K, V] {
	return hashmap.NewBiHashmap[K, V]()
}
