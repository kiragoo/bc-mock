package main

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/emqx/emqx-bc-asynctasks/task"
	"github.com/emqx/emqx-operator/apis/apps/v1beta2"
)

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
		Env: map[string]string{
			"EMQX_CLUSTER__K8S__NAMESPACE": ns,
		},
		Certificate: generateCertificate(),
		// License:     []byte(LIC),
		LogType: "external",
		LogConfig: map[string]string{
			"EsUrl":       "http://139.224.110.115:8080/es9200",
			"EsUsername":  "esUsername",
			"EsPassword":  "esPassword",
			"EsIndexName": "EsIndexName",
		},
		MetricConfig: map[string]string{
			"BCMetricGateway": "http://www.example.com/",
		},
		ResourceLimits: task.ResourceLimits{
			Emqx: task.ResourceLimit{
				Requests: task.ResourceQuota{
					CPU:    "250m",
					Memory: "100Mi",
				},
				Limits: task.ResourceQuota{
					CPU:    "500m",
					Memory: "200Mi",
				},
			},
			Telegraf: task.ResourceLimit{
				Requests: task.ResourceQuota{
					CPU:    "250m",
					Memory: "100Mi",
				},
				Limits: task.ResourceQuota{
					CPU:    "1000m",
					Memory: "200Mi",
				},
			},
		},
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
		LogType: "external",
		LogConfig: map[string]string{
			"EsUrl":       "http://139.224.110.115:8080/es9200",
			"EsUsername":  "esUsername",
			"EsPassword":  "esPassword",
			"EsIndexName": "EsIndexName",
		},
		MetricConfig: map[string]string{
			"BCMetricGateway": "http://www.example.com/",
		},
		ResourceLimits: task.ResourceLimits{
			Emqx: task.ResourceLimit{
				Requests: task.ResourceQuota{
					CPU:    "250m",
					Memory: "100Mi",
				},
				Limits: task.ResourceQuota{
					CPU:    "500m",
					Memory: "200Mi",
				},
			},
			Telegraf: task.ResourceLimit{
				Requests: task.ResourceQuota{
					CPU:    "250m",
					Memory: "100Mi",
				},
				Limits: task.ResourceQuota{
					CPU:    "1000m",
					Memory: "200Mi",
				},
			},
		},
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

func updateLicense(id uint64, ns string) {
	msg := &task.Task{
		TaskID:    id,
		Type:      task.TaskUpdateLicense,
		Namespace: ns,
		License:   []byte(LIC),
		ClusterID: "cluster-" + strconv.FormatUint(id, 10),
	}
	task.Deploy(msg)
}

func updateCertificate(id uint64, ns string) {
	msg := &task.Task{
		TaskID:      id,
		Type:        task.TaskUpdateCertificate,
		Namespace:   ns,
		Certificate: generateCertificate(),
		ClusterID:   "cluster-" + strconv.FormatUint(id, 10),
	}
	task.Deploy(msg)
}

func updateImage(id uint64, ns string) {
	msg := &task.Task{
		TaskID:    id,
		Type:      task.TaskUpdateImage,
		Namespace: ns,
		Image:     "emqx/emqx-ee:4.4.1",
		ClusterID: "cluster-" + strconv.FormatUint(id, 10),
	}
	task.Deploy(msg)
}

func scaleUp(id uint64, ns string) {
	msg := &task.Task{
		TaskID:    id,
		Type:      task.TaskScaleUpDeployment,
		Namespace: ns,
		ClusterID: "cluster-" + strconv.FormatUint(id, 10),
		Replicas:  4,
	}
	task.Deploy(msg)
}

// helpers
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
