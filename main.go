package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/emqx/emqx-bc-asynctasks/role"
	"github.com/emqx/emqx-bc-asynctasks/task"
	"github.com/emqx/emqx-bc-asynctasks/taskmq"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

var logger *log.Logger

func init() {
	writer1, err := os.OpenFile("logger.txt", os.O_WRONLY|os.O_CREATE, 0755)
	writer2 := os.Stdout
	if err != nil {
		log.Fatalf("create file logger.txt failed: %v", err)
	}

	logger = log.New(io.MultiWriter(writer1, writer2), "", log.Lshortfile|log.LstdFlags)
}

type Config struct {
	TaskMqConfig *taskmq.BCAsyncTaskConfig `yaml:"taskmq"`
}

func getConfig(file string) *taskmq.BCAsyncTaskConfig {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		panic(err)
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		panic(err)
	}
	return cfg.TaskMqConfig
}

var tasks = map[string]interface{}{
	task.TASK_EMQX_CLUSTER_UPDATE: emqxClusterStatus,
	task.RECALL_KUBECONFIG:        pullKubeConfigs,
	task.TASK_BC_UPDATE_TASK:      deploymentStatus,
}

var taskList = map[string]interface{}{
	SEND_KUBE_CONFIG:   sendKubeConfig,
	CREATE_DEPLOYMENT:  createDeployment,
	STOP_DEPLOYMENT:    stopDeployment,
	START_DEPLOYMENT:   startDeployment,
	DELETE_DEPLOYMENT:  deleteDeployment,
	UPDATE_LICENSE:     updateLicense,
	UPDATE_CERTIFICATE: updateCertificate,
	UPDATE_IMAGE:       updateImage,
	SCALE_UP:           scaleUp,
}

func main() {
	config := getConfig("./conf.yaml")
	taskmq.Server.CreateMachineryServer("test_worker", config)
	role.RegisterDefaultTasks()
	taskmq.Server.RegisterTasks(tasks)

	{
		prompt := promptui.Select{
			Label: "Select Tasks",
			Items: []string{SEND_KUBE_CONFIG, CREATE_DEPLOYMENT, STOP_DEPLOYMENT, START_DEPLOYMENT, DELETE_DEPLOYMENT,
				UPDATE_LICENSE, UPDATE_CERTIFICATE, UPDATE_IMAGE, SCALE_UP, EXIT},
		}

		_, result, err := prompt.Run()

		if err != nil {
			logger.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case SEND_KUBE_CONFIG:
			taskList[SEND_KUBE_CONFIG].(func(uint64))(TASK_ID)
		case CREATE_DEPLOYMENT:
			taskList[CREATE_DEPLOYMENT].(func(uint64, string))(TASK_ID, NAMESPACE)
		case STOP_DEPLOYMENT:
			taskList[STOP_DEPLOYMENT].(func(uint64, string))(TASK_ID, NAMESPACE)
		case START_DEPLOYMENT:
			taskList[START_DEPLOYMENT].(func(uint64, string))(TASK_ID, NAMESPACE)
		case DELETE_DEPLOYMENT:
			taskList[DELETE_DEPLOYMENT].(func(uint64, string))(TASK_ID, NAMESPACE)
		case UPDATE_LICENSE:
			taskList[UPDATE_LICENSE].(func(uint64, string))(TASK_ID, NAMESPACE)
		case UPDATE_CERTIFICATE:
			taskList[UPDATE_CERTIFICATE].(func(uint64, string))(TASK_ID, NAMESPACE)
		case UPDATE_IMAGE:
			taskList[UPDATE_IMAGE].(func(uint64, string))(TASK_ID, NAMESPACE)
		case SCALE_UP:
			taskList[SCALE_UP].(func(uint64, string))(TASK_ID, NAMESPACE)
		case EXIT:
			return
		}
		taskmq.Server.LaunchWorker(task.ROUTINGKEY_BC_TASKS)
	}
}
