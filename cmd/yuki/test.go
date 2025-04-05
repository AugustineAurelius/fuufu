package yuki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/AugustineAurelius/fuufu/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type put struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
type get struct {
	Value json.RawMessage `json:"value"`
}

func registerTestCMD(manager *config.Manager) *cobra.Command {
	testCMD := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewWithManager(manager)

			cfg := manager.LoadYuki()
			url := "http://" + cfg.Host + ":" + strconv.Itoa(cfg.Port)

			log.Info("created url", zap.String("url", url))
			p := put{
				Key:   "123",
				Value: 123,
			}
			data, err := json.Marshal(&p)
			if err != nil {
				return err
			}
			resp, err := http.Post(url, "application/json", bytes.NewReader(data))
			if err != nil {
				return err
			}

			if resp.StatusCode != 201 {
				return fmt.Errorf("bad status code")
			}

			log.Info("post successfull")

			resp, err = http.Get(url + "/123")
			if err != nil {
				return err
			}

			var g get
			if err = json.NewDecoder(resp.Body).Decode(&g); err != nil {
				return err
			}

			log.Info("get successfull", zap.String("value", string(g.Value)))

			return nil
		},
	}

	return testCMD
}
