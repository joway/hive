package hive

type Task func()
type Callback func()
type Job struct {
	task     Task
	callback Callback
}

type Worker struct {
	jobs chan Job
}

func NewWorker() *Worker {
	w := &Worker{
		jobs: make(chan Job),
	}
	go w.Run()
	return w
}

func (w *Worker) Submit(task Task, callback Callback) {
	w.jobs <- Job{
		task:     task,
		callback: callback,
	}
}

func (w *Worker) Run() {
	for job := range w.jobs {
		if job.task != nil {
			job.task()
		}
		if job.callback != nil {
			job.callback()
		}
	}
}

func (w *Worker) Close() {
	close(w.jobs)
}
