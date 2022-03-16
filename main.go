package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/emqx/emqx-bc-asynctasks/role"
	"github.com/emqx/emqx-bc-asynctasks/task"
	"github.com/emqx/emqx-bc-asynctasks/taskmq"
	"github.com/emqx/emqx-operator/apis/apps/v1beta2"
	"gopkg.in/yaml.v2"
)

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

func getKubeConfig() []byte {
	config, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err != nil {
		panic(err.Error())
	}

	return config
}

func sendKubeConfig(id uint64) {
	msg := &task.KubeConfig{
		Configs: map[string][]byte{"cluster-" + strconv.FormatUint(id, 10): getKubeConfig()},
	}
	task.SendKubeConfig(msg)
}

func createDeployment(id uint64, ns string) {
	msg := &task.Task{
		TaskID:           id,
		Type:             task.TaskCreateDeployment,
		Namespace:        ns,
		ClusterID:        "cluster-" + strconv.FormatUint(id, 10),
		Image:            "emqx/emqx-ee:4.4.0",
		Replicas:         3,
		StorageClassName: "standard",
		StorageClassSize: "20Mi",
		// Labels: map[string]string{
		// 	"cluster": "emqx",
		// },
		Env: map[string]string{
			"EMQX_CLUSTER__K8S__NAMESPACE": ns,
		},
		Certificate: generateCertificate(),
		License:     []byte(LIC),
	}
	task.Deploy(msg)
}

func stopDeployment(id uint64, ns string) {
	msg := &task.Task{
		TaskID:    id,
		Type:      task.TaskStopDeployment,
		Namespace: ns,
		ClusterID: "cluster-" + strconv.FormatUint(id, 10),
	}
	task.Deploy(msg)
}

func deleteDeployment(id uint64, ns string) {
	msg := &task.Task{
		TaskID:    id,
		Type:      task.TaskDeleteDeployment,
		Namespace: ns,
		Replicas:  3,
		ClusterID: "cluster-" + strconv.FormatUint(id, 10),
	}
	task.Deploy(msg)
}

func startDeployment(id uint64, ns string) {
	msg := &task.Task{
		TaskID:           id,
		Type:             task.TaskStartDeployment,
		Namespace:        ns,
		ClusterID:        "cluster-" + strconv.FormatUint(id, 10),
		Image:            "emqx/emqx-ee:4.4.0",
		Replicas:         3,
		StorageClassName: "standard",
		StorageClassSize: "20Mi",
		Env: map[string]string{
			"EMQX_CLUSTER__K8S__NAMESPACE": ns,
		},
		// Status: task.DeploymentNotStarted,
	}
	task.Deploy(msg)
}

type Cluster struct {
}

func (c *Cluster) SubEmqxClusterStatus(status []byte) error {
	s := task.EmqxClusterStatus{}

	if err := json.Unmarshal(status, &s); err != nil {
		panic(err.Error())
	}

	log.Printf("---- emqx status: %+v\n", s)
	return nil
}

func main() {
	config := getConfig("./conf.yaml")
	taskmq.Server.CreateMachineryServer("test_worker", config)
	role.RegisterDefaultTasks()
	task.RegisterEmqxActions(&Cluster{})

	sendKubeConfig(1)
	time.Sleep(time.Second * 2)
	createDeployment(1, "emqx-1")

	// deleteDeployment(1, "emqx-1")
	// stopDeployment(1, "emqx-1")

	taskmq.Server.LaunchWorker(task.ROUTINGKEY_BC_TASKS)
}

func generateCertificate() v1beta2.Certificate {
	certConf := v1beta2.CertificateConf{
		StringData: v1beta2.CertificateStringData{
			CaCert:  CA_CERT,
			TLSCert: TLS_CERT,
			TLSKey:  TLS_KEY,
		},
	}
	return v1beta2.Certificate{
		WSS:   certConf,
		MQTTS: certConf,
	}
}
