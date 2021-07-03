package longjob

type Queue struct {
	Name string
	WorkerCount uint
}

func NewQueue(queueName string, workerCount uint) Queue {
	return Queue{
		Name: queueName,
		WorkerCount: workerCount,
	}
}
