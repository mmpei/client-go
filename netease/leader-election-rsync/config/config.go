package config

type Configuration struct {
	// KubeConfig the kubeconfig, could not provide when running in k8s cluster
	KubeConfig string `yaml:"kube_config"`
	// Id
	Id string `yaml:"id"`
	// ConfigMapLockName
	ConfigMapLockName string `yaml:"config_map_lock_name"`
	// ConfigMapLockNamespace
	ConfigMapLockNamespace int `yaml:"config_map_lock_namespace"`
}
