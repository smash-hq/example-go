package main

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var K8sClient = &Client{}

type Client struct {
	Clientset *kubernetes.Clientset
}

func init() {
	conf := k8sRestConfig()
	K8sClient.Clientset = initClient(conf)
	printCurrentNamespaceInfo(K8sClient.Clientset)
	log.Info("Init k8s client success")
}

func printCurrentNamespaceInfo(clientset kubernetes.Interface) {
	// 获取集群版本
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		log.Errorf("Failed to get server version: %s", err.Error())
	} else {
		log.Info("Cluster Information:")
		log.Infof("  Kubernetes Version: %s", version.String())
	}
}

// initClient 返回初始化k8s-client
func initClient(conf rest.Config) *kubernetes.Clientset {
	// 创建一个Kubernetes的Clientset对象
	clientset, err := kubernetes.NewForConfig(&conf)
	if err != nil {
		log.Errorf("Create clientset failed, cause: %s", err.Error())
	}
	return clientset
}

func main() {
	log.Println("test get k8s info")
	select {}
}

func k8sRestConfig() rest.Config {
	var cfg rest.Config
	cfg, err := getRestConfByEnv()
	if err != nil {
		log.Fatalf("Error building kubeconfig, cause: %s", err.Error())
	}
	return cfg
}

// getRestConfByEnv 获取k8s restful client配置
func getRestConfByEnv() (rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return rest.Config{}, err
	}
	// 提高 QPS 和 Burst 限制
	config.QPS = 100
	config.Burst = 200
	return *config, nil
}
