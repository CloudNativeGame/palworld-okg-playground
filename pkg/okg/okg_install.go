package okg

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"strings"
	"time"
)

type DeployRequest struct {
	RepoURL      string                 // 仓库地址
	ChartName    string                 // Chart名称
	ChartVersion string                 // Chart版本
	Namespace    string                 // 命名空间
	ReleaseName  string                 // 在kubernetes中的程序名
	Values       map[string]interface{} // values.yaml 配置文件
	config       *rest.Config
}

func InstallOpenKruiseGame(config *rest.Config) error {
	entry := &repo.Entry{
		Name: "openkruise",
		URL:  "https://openkruise.github.io/charts/",
	}
	if err := add(entry); err != nil {
		return err
	}

	if err := update(); err != nil {
		return err
	}

	kruiseRequest := &DeployRequest{
		RepoURL:      "https://openkruise.github.io/charts/",
		ChartName:    "kruise",
		ChartVersion: "1.5.1",
		Namespace:    "default",
		ReleaseName:  "kruise",
		config:       config,
	}

	if err := installChart(kruiseRequest); err != nil {
		return err
	}

	kruiseGameRequest := &DeployRequest{
		RepoURL:      "https://openkruise.github.io/charts/",
		ChartName:    "kruise-game",
		ChartVersion: "0.7.0",
		Namespace:    "default",
		ReleaseName:  "kruise-game",
		config:       config,
	}

	if err := installChart(kruiseGameRequest); err != nil {
		return err
	}

	return nil
}

func add(entry *repo.Entry) error {
	settings := cli.New()

	repoFile := settings.RepositoryConfig

	// 加载仓库配置文件
	repositories, err := repo.LoadFile(repoFile)
	// 如果文件不存在
	if err != nil {
		// 创建一个新的仓库配置对象
		repositories = repo.NewFile()
	}

	// 检查要添加的仓库是否已存在
	if repositories.Has(entry.Name) {
		log.Printf("仓库 %s 已存在", entry.Name)
	}

	// 添加仓库信息到仓库配置
	repositories.Add(entry)

	// 保存更新后的仓库配置到文件
	if err = repositories.WriteFile(repoFile, 0644); err != nil {
		return fmt.Errorf("无法保存仓库配置文件：%s", err)
	}

	log.Printf("成功添加仓库地址：%s。", entry.Name)
	return nil
}

func update() error {
	settings := cli.New()
	// 加载仓库配置文件
	repositories, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return fmt.Errorf("无法加载仓库配置文件：%s", err)
	}

	// 遍历每个仓库
	for _, repoEntry := range repositories.Repositories {
		if repoEntry.Name != "openkruise" {
			continue
		}
		// 添加要检索的仓库
		chartRepository, err := repo.NewChartRepository(repoEntry, getter.All(settings))
		if err != nil {
			return fmt.Errorf("无法添加仓库：%s\n", err)
		}

		// 更新仓库索引信息
		if _, err := chartRepository.DownloadIndexFile(); err != nil {
			return fmt.Errorf("无法下载仓库索引：%s\n", err)
		}

		log.Printf("...Successfully got an update from the %s chart repository", repoEntry.Name)
		break
	}

	log.Println("Update Complete. ⎈Happy Helming!⎈")
	return nil
}

func installChart(deployRequest *DeployRequest) error {
	settings := cli.New()

	actionConfig := new(action.Configuration)
	restGetter := genericclioptions.NewConfigFlags(false)
	restGetter.WrapConfigFn = func(*rest.Config) *rest.Config {
		return deployRequest.config
	}
	if err := actionConfig.Init(restGetter, deployRequest.Namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return fmt.Errorf("初始化 action 失败\n%s", err)
	}

	install := action.NewInstall(actionConfig)
	install.RepoURL = deployRequest.RepoURL
	install.Version = deployRequest.ChartVersion
	install.Timeout = 300 * time.Second
	install.CreateNamespace = true
	install.Wait = true
	// kubernetes 中的配置
	install.Namespace = deployRequest.Namespace
	install.ReleaseName = deployRequest.ReleaseName

	chartRequested, err := install.ChartPathOptions.LocateChart(deployRequest.ChartName, settings)
	if err != nil {
		return fmt.Errorf("下载失败\n%s", err)
	}

	chart, err := loader.Load(chartRequested)
	if err != nil {
		return fmt.Errorf("加载失败\n%s", err)
	}

	_, err = install.Run(chart, nil)
	if err != nil && !strings.Contains(err.Error(), "still in use") {
		return fmt.Errorf("执行失败\n%s", err)
	}

	return nil
}
