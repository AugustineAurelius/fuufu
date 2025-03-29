package server

import (
	"context"

	"github.com/AugustineAurelius/fuufu/api/todo"
	todo_repository "github.com/AugustineAurelius/fuufu/internal/repository/todo"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TodoHandler struct {
	todo.StrictServerInterface

	Repo      *todo_repository.CommandRepository
	Telemetry trace.Tracer
}

// Get Todo list
// (GET /api/v1/todo)
func (h *TodoHandler) GetAllTodos(ctx context.Context, _ todo.GetAllTodosRequestObject) (todo.GetAllTodosResponseObject, error) {
	ctx, span := h.Telemetry.Start(ctx, "TodoHandler.GetAllTodos")
	defer span.End()

	todos, err := h.Repo.GetMany(ctx)
	if err != nil {
		return todo.GetAllTodos500JSONResponse{
			Error: err.Error(),
		}, nil
	}

	result := make([]todo.Task, 0, len(todos))
	for _, task := range todos {
		result = append(result, todo.Task{
			CreatedBy:   todo.WeAll(task.CreatedBy),
			Description: task.Description,
			DoBefore:    *task.DoBefore,
			Doer:        todo.WeAll(task.Doer),
			Done:        task.Done,
			Id:          task.ID,
			Name:        task.Name,
			Range:       task.RepeatAfter,
			Repeatable:  task.Repeatable,
		})
	}

	return todo.GetAllTodos200JSONResponse{Tasks: result}, nil
}

// Creates a new task
// (POST /api/v1/todo)
func (h *TodoHandler) CreateNewTask(ctx context.Context, request todo.CreateNewTaskRequestObject) (todo.CreateNewTaskResponseObject, error) {
	ctx, span := h.Telemetry.Start(ctx, "TodoHandler.CreateNewTask", trace.WithAttributes(
		attribute.String("created_by", string(request.Body.CreatedBy)),
		attribute.String("doer", string(request.Body.Doer)),
		attribute.String("name", request.Body.Name),
		attribute.Bool("repeateble", request.Body.Repeatable),
	))
	defer span.End()

	id := uuid.New()
	span.SetAttributes(attribute.Stringer("todo_id", id))

	err := h.Repo.Create(ctx, &todo_repository.Task{
		ID:          id,
		CreatedBy:   string(request.Body.CreatedBy),
		Description: request.Body.Description,
		DoBefore:    &request.Body.DoBefore,
		Doer:        string(request.Body.Doer),
		Name:        request.Body.Name,
		RepeatAfter: request.Body.Range,
		Repeatable:  request.Body.Repeatable,
	})
	if err != nil {
		return todo.CreateNewTask500JSONResponse{
			Error: err.Error(),
		}, nil
	}

	return todo.CreateNewTask201JSONResponse{
		TaskId: id,
	}, nil
}

// Get task by todo_id
// (GET /api/v1/todo/{todo_id})
func (h *TodoHandler) GetTaskByID(ctx context.Context, request todo.GetTaskByIDRequestObject) (todo.GetTaskByIDResponseObject, error) {
	ctx, span := h.Telemetry.Start(ctx, "TodoHandler.GetTodoByID", trace.WithAttributes(
		attribute.Stringer("todo_id", request.TodoId),
	))
	defer span.End()

	todoOne, err := h.Repo.Get(ctx, todo_repository.WithID(request.TodoId))
	if err != nil {
		return todo.GetTaskByID500JSONResponse{
			Error: err.Error(),
		}, nil
	}

	return todo.GetTaskByID200JSONResponse{
		Id:          todoOne.ID,
		CreatedBy:   todo.WeAll(todoOne.CreatedBy),
		Description: todoOne.Description,
		DoBefore:    *todoOne.DoBefore,
		Doer:        todo.WeAll(todoOne.Doer),
		Name:        todoOne.Name,
		Range:       todoOne.RepeatAfter,
		Repeatable:  todoOne.Repeatable,
	}, nil
}

// delete task by todo_id
// (DELETE /api/v1/todo/{todo_id})
func (h *TodoHandler) DeleteTaskByID(ctx context.Context, request todo.DeleteTaskByIDRequestObject) (todo.DeleteTaskByIDResponseObject, error) {
	ctx, span := h.Telemetry.Start(ctx, "TodoHandler.DeleteByID", trace.WithAttributes(
		attribute.Stringer("todo_id", request.TodoId),
	))
	defer span.End()

	_, err := h.Repo.Get(ctx, todo_repository.WithID(request.TodoId))
	if err != nil {
		return todo.DeleteTaskByID404Response{}, nil
	}

	err = h.Repo.Delete(ctx, todo_repository.WithID(request.TodoId))
	if err != nil {
		return todo.DeleteTaskByID500JSONResponse{
			Error: err.Error(),
		}, nil
	}

	return todo.DeleteTaskByID200Response{}, nil
}
