package main

import (
	"encoding/json"

	"github.com/emqx/emqx-bc-asynctasks/task"
)

func emqxClusterStatus(status []byte) error {
	s := task.EmqxClusterStatus{}

	if err := json.Unmarshal(status, &s); err != nil {
		panic(err.Error())
	}

	logger.Printf("---- emqx status: %+v\n", s)
	return nil
}

func pullKubeConfigs(configs []byte) error {
	c := task.KubeConfig{}

	// if err := json.Unmarshal(configs, &c); err != nil {
	// 	panic(err.Error())
	// }

	logger.Printf("---- kube configs: %+v\n", c)

	return nil
}

func deploymentStatus(status []byte) error {
	s := task.BCTaskStatus{}

	if err := json.Unmarshal(status, &s); err != nil {
		panic(err.Error())
	}

	logger.Printf("---- deployment status: %+v\n", s)
	return nil
}
