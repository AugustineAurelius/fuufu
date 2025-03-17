package analytic

import (
	"context"
	"fmt"
	"time"

	todo_repository "github.com/AugustineAurelius/fuufu/internal/repository/todo"
	"go.opentelemetry.io/otel/metric"
)

type Analytic struct {
	Meter    metric.Meter
	TodoRepo *todo_repository.QueryRepository
}

func (a *Analytic) Run(ctx context.Context) error {

	allTodoTasks, err := a.Meter.Int64Gauge("current_todo_tasks")
	if err != nil {
		return err
	}
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for range t.C {
		tasks, err := a.getAllTodoCount(ctx)
		if err != nil {
			fmt.Println(err)
			continue
		}
		allTodoTasks.Record(ctx, tasks)

	}
	return nil
}

func (a Analytic) getAllTodoCount(ctx context.Context) (int64, error) {
	tasks, err := a.TodoRepo.GetMany(ctx, todo_repository.WithDone(false))
	if err != nil {
		return 0, err
	}

	return int64(len(tasks)), nil
}
