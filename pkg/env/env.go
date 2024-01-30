package env

import (
	"bufio"
	"bytes"
	"github.com/BurntSushi/toml"
	"os"
	"path"
	"strings"
)

const (
	defaultEnvFile = ".pal.toml"
)

type ClusterConfiguration struct {
	Default                  bool   `toml:"default"`
	ID                       string `toml:"id"`
	Name                     string `toml:"name"`
	GameServerLoadBalancerId string `toml:"gs_lb_id"`
}

type Configuration struct {
	Clusters []ClusterConfiguration `toml:"cluster"`
}

type EnvFile struct {
	path string
}

func NewEnvFile() *EnvFile {
	p, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return &EnvFile{
		path: path.Join(p, defaultEnvFile),
	}
}

func (e *EnvFile) Exists() bool {
	if _, err := os.Stat(e.path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (e *EnvFile) Create() error {
	_, err := os.Create(e.path)
	if err != nil {
		return err
	}
	return err
}

func (e *EnvFile) Read() (string, error) {
	return "", nil
}

func (e *EnvFile) ListAll() (Configuration, error) {
	var config Configuration

	// 读取 TOML 文件
	file, err := os.ReadFile(e.path)
	if err != nil {
		return config, err
	}

	// 解析 TOML 文件
	if _, err := toml.Decode(string(file), &config); err != nil {
		return config, err
	}

	return config, nil
}

func (e *EnvFile) ReadAll() (map[string]string, error) {
	file, err := os.Open(e.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	clusterMap := make(map[string]string)
	// 使用 bufio.NewReader 读取文件
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n') // 读取直到遇到换行符
		if err != nil {
			break
		}
		linePure := strings.Split(line, "\n")[0]
		nameId := strings.Split(linePure, ":")
		if len(nameId) != 2 {
			continue
		}
		clusterMap[nameId[0]] = nameId[1]
	}

	return clusterMap, nil
}

func (e *EnvFile) Write(stringToWrite string) error {
	file, err := os.Open(e.path)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(e.path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(stringToWrite + "\n") // 写入字符串并添加换行符
	if err != nil {
		return err
	}

	writer.Flush() // 确保所有数据都写入了文件

	return nil
}

func (e *EnvFile) AddNewCluster(newCluster *ClusterConfiguration) error {
	// 读取 TOML 文件
	file, err := os.ReadFile(e.path)
	if err != nil {
		return err
	}

	// 解析 TOML 文件
	var config Configuration
	if _, err := toml.Decode(string(file), &config); err != nil {
		return err
	}

	config.Clusters = append(config.Clusters, *newCluster)

	// 将新的配置写回文件
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	// 将缓冲区内容写入文件
	err = os.WriteFile(e.path, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (e *EnvFile) UpdateDefaultCluster(clusterId string) {
	// 读取 TOML 文件
	file, err := os.ReadFile(e.path)
	if err != nil {
		panic(err)
	}

	// 解析 TOML 文件
	var config Configuration
	if _, err := toml.Decode(string(file), &config); err != nil {
		panic(err)
	}

	for i, cluster := range config.Clusters {
		if cluster.Default {
			config.Clusters[i].Default = false
		}
		if cluster.ID == clusterId {
			config.Clusters[i].Default = true
		}
	}

	// 将新的配置写回文件
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err = encoder.Encode(config)
	if err != nil {
		panic(err)
	}

	// 将缓冲区内容写入文件
	err = os.WriteFile(e.path, buffer.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func (e *EnvFile) GetDefaultId() (string, string) {
	// 读取 TOML 文件
	file, err := os.ReadFile(e.path)
	if err != nil {
		panic(err)
	}

	// 解析 TOML 文件
	var config Configuration
	if _, err := toml.Decode(string(file), &config); err != nil {
		panic(err)
	}

	for _, cluster := range config.Clusters {
		if cluster.Default {
			return cluster.ID, cluster.GameServerLoadBalancerId
		}
	}
	return "", ""
}
