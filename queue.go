package go_ds

type Queue struct {
	dataCh chan interface{}
}

func (q *Queue) Push(data interface{}) {
	q.dataCh <- data
}

func (q *Queue) Pop() interface{} {
	return <-q.dataCh
}

func NewQueue() *Queue {
	return &Queue{
		dataCh: make(chan interface{}),
	}
}
