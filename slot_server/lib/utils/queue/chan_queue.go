package queue

import (
	"errors"
	"sync"
)

// QueueDataKeyMap map
type QueueDataKeyMap struct {
	Mu     *sync.RWMutex
	KeyMap map[string]int64
}

func NewQueueDataKeyMap() *QueueDataKeyMap {
	return &QueueDataKeyMap{
		Mu:     &sync.RWMutex{},
		KeyMap: make(map[string]int64, 0),
	}
}

func (qm *QueueDataKeyMap) TryAdd(key string, val int64) error {
	qm.Mu.Lock()
	defer qm.Mu.Unlock()

	getVal, _ := qm.KeyMap[key]
	if getVal > 0 {
		return errors.New("key already exists")
	}

	qm.KeyMap[key] = val
	return nil
}

func (qm *QueueDataKeyMap) Add(key string, val int64) {
	qm.Mu.Lock()
	defer qm.Mu.Unlock()
	qm.KeyMap[key] = val
}

func (qm *QueueDataKeyMap) Del(key string) {
	qm.Mu.Lock()
	defer qm.Mu.Unlock()
	delete(qm.KeyMap, key)
}

func (qm *QueueDataKeyMap) Get(key string) int64 {
	qm.Mu.RLock()
	defer qm.Mu.RUnlock()
	v, ok := qm.KeyMap[key]
	if !ok {
		return -1
	}
	return v
}

type QueueData struct {
	Key  string //用户ID + func || 其他
	Data []byte
}

// 简单队列结构
type Queue struct {
	ch chan QueueData // 用一个 channel 作为队列
}

// 创建队列
func NewQueue(capacity int) *Queue {
	return &Queue{ch: make(chan QueueData, capacity)}
}

// 往队列里加数据
func (q *Queue) Enqueue(item QueueData) {
	q.ch <- item // 放不下的时候，会阻塞在这里
}

// 从队列里取数据
func (q *Queue) Dequeue() QueueData {
	return <-q.ch // 没数据的时候，会阻塞在这里
}

// 获取当前队列长度
func (q *Queue) Length() int {
	return len(q.ch)
}

// 获取队列容量
func (q *Queue) Capacity() int {
	return cap(q.ch)
}
