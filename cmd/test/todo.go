package test

import (
	"fmt"
	"time"

	"github.com/AugustineAurelius/fuufu/api/todo"
	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/AugustineAurelius/fuufu/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func createTodoCMD(manager *config.Manager) *cobra.Command {

	todoCMD := &cobra.Command{
		Use: "todo",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.NewWithManager(manager)

			createResponse, err := todoClient.CreateNewTaskWithResponse(cmd.Context(), todo.CreateNewTaskJSONRequestBody{
				CreatedBy: todo.Kirill,
				DoBefore:  time.Now().UTC().Add(time.Hour),
				Doer:      todo.Dasha,
				Name:      "Uberi kakashki",
			}, HeaderCreator)
			checkError(err)
			checkResponse(createResponse.HTTPResponse)

			log.Info("task created")

			getAllResponse, err := todoClient.GetAllTodosWithResponse(cmd.Context(), HeaderCreator)
			checkError(err)
			checkResponse(getAllResponse.HTTPResponse)
			log.Info("get all tasks", zap.Any("by id", string(getAllResponse.Body)))

			var index int
			for i, task := range getAllResponse.JSON200.Tasks {
				if task.Id == createResponse.JSON201.TaskId {
					index = i
					break
				}
			}
			if createResponse.JSON201.TaskId != getAllResponse.JSON200.Tasks[index].Id {
				checkError(fmt.Errorf("ids not equal"))
			}

			log.Info("tasks ID is equal")

			getByIDResponse, err := todoClient.GetTaskByIDWithResponse(cmd.Context(), createResponse.JSON201.TaskId, HeaderCreator)
			checkError(err)
			checkResponse(getByIDResponse.HTTPResponse)
			log.Info("find by id", zap.Any("by id", string(getByIDResponse.Body)))

			deleteResponse, err := todoClient.DeleteTaskByIDWithResponse(cmd.Context(), createResponse.JSON201.TaskId, HeaderCreator)
			checkError(err)
			checkResponse(deleteResponse.HTTPResponse)

			log.Info("success test")
		},
	}

	return todoCMD
}
