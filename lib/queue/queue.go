package queue

import "sync/atomic"

type (
	//Queue 队列
	Queue struct {
		top    *node
		rear   *node
		length int64
	}
	//双向链表节点
	node struct {
		pre   *node
		next  *node
		value interface{}
	}
)

// Create a new queue
func New() *Queue {
	return &Queue{nil, nil, 0}
}

//获取队列长度
func (q *Queue) Len() int64 {
	return q.length
}

//返回队列顶端元素
func (q *Queue) Peek() interface{} {
	if q.top == nil {
		return nil
	}
	return q.top.value
}

//入队操作
func (q *Queue) Push(v interface{}) {
	n := &node{nil, nil, v}
	if q.length == 0 {
		q.top = n
		q.rear = q.top
	} else {
		n.pre = q.rear
		q.rear.next = n
		q.rear = n
	}
	atomic.AddInt64(&q.length, 1)
}

//出队操作
func (q *Queue) Pop() interface{} {
	if q.length == 0 {
		return nil
	}
	n := q.top
	if q.top != nil {
		if q.top.next == nil {
			q.top = nil
		} else {
			q.top = q.top.next
			//q.top.pre.next = nil
			q.top.pre = nil
		}
	}
	atomic.AddInt64(&q.length, -1)
	if n != nil {
		return n.value
	} else {
		return nil
	}
}
