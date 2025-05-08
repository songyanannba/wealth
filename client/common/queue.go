package common

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type synQueue struct {
	inLen    int32
	outLen   int32
	inLock   *sync.Mutex
	outLock  *sync.Mutex
	inQueue  *list.List
	outQueue *list.List
}

func NewSynQueue() *synQueue {
	return &synQueue{
		inLock:   new(sync.Mutex),
		outLock:  new(sync.Mutex),
		inQueue:  list.New(),
		outQueue: list.New(),
	}
}

// 添加队列
func (q *synQueue) Add(i interface{}) {
	q.inLock.Lock()
	defer q.inLock.Unlock()
	q.inQueue.PushBack(i)
	atomic.AddInt32(&q.inLen, 1)
}

// 交换 队列
func (q *synQueue) SwapQue() {
	q.inLock.Lock()
	q.outLock.Lock()

	q.inQueue, q.outQueue = q.outQueue, q.inQueue
	q.inLen, q.outLen = q.outLen, q.inLen

	q.outLock.Unlock()
	q.inLock.Unlock()
}

func (q *synQueue) outQueLen() int {
	return int(atomic.LoadInt32(&q.outLen))
}

func (q *synQueue) ConsumeQue() *list.List {
	q.outLock.Lock()
	defer q.outLock.Unlock()

	returnList := &list.List{}

	for q.outQueLen() > 0 {
		retQue := q.RemoveQue()
		if retQue != nil {
			returnList.PushBack(retQue)
		}
		atomic.AddInt32(&q.outLen, -1)
	}

	return returnList
}

func (q *synQueue) SwapAndConsumeQueue() *list.List {

	q.SwapQue()
	return q.ConsumeQue()
}

// 删除
func (q *synQueue) RemoveQue() interface{} {
	return q.outQueue.Remove(q.outQueue.Front())
}
