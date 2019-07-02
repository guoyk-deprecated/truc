package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	optHome       string
	optDeployment string
	optEnv        string
	optCluster    string
	optNamespace  string
	optRegistry   string
)

func checkFile(name string, exec bool) (err error) {
	log.Println("检查文件 " + name)
	var info os.FileInfo
	if info, err = os.Stat(name); err != nil {
		return
	}
	if exec {
		mode := info.Mode()
		if mode&0111 != 0111 {
			return os.Chmod(name, mode|0111)
		}
	}
	return
}

func run(rel bool, saveOut bool, name string, args ...string) (out string, err error) {
	log.Println("执行命令 " + name + " " + strings.Join(args, " "))
	if rel {
		var wd string
		if wd, err = os.Getwd(); err != nil {
			return
		}
		name = filepath.Join(wd, name)
	}
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	if saveOut {
		buf := &bytes.Buffer{}
		cmd.Stdout = buf
		err = cmd.Run()
		out = strings.TrimSpace(buf.String())
	} else {
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	}
	return
}

func exit(err *error) {
	if *err != nil {
		log.Fatalln(*err)
	}
}

func checkStr(str *string, desc string) {
	*str = strings.TrimSpace(*str)
	if *str == "" {
		log.Fatal(desc)
	}
}

func main() {
	var err error
	defer exit(&err)

	log.SetPrefix("deployer: ")
	log.SetFlags(0)

	flag.StringVar(&optRegistry, "registry", "", "镜像仓库，指定镜像要推往的仓库地址")
	flag.StringVar(&optCluster, "cluster", "", "集群，决定使用哪个 kubectl 配置文件，一般和 Rancher 控制台左上角集群名一致")
	flag.StringVar(&optNamespace, "namespace", "", "命名空间，工作负载在集群中的命名空间")
	flag.StringVar(&optDeployment, "deployment", "", "工作负载，要求和 Kubernetes 上的工作负载名称完全一致")
	flag.StringVar(&optEnv, "env", "", "环境，决定 docker-build.XXX.sh， Dockerfile.XXX 文件的选择，和构建后的镜像标签")
	flag.Parse()

	// extract deployment/env from JOB_NAME
	if optCluster == "" || optNamespace == "" || optDeployment == "" || optEnv == "" {
		split := strings.Split(os.Getenv("JOB_NAME"), ".")
		if len(split) == 2 {
			optDeployment = strings.ReplaceAll(split[0], "_", "-")
			optEnv = split[1]
			log.Println("从 Jenkins $JOB_NAME 获取到 工作负载 " + optDeployment + ", 环境 " + optEnv)
		}
		if len(split) == 4 {
			optCluster = strings.ReplaceAll(split[0], "_", "-")
			optNamespace = strings.ReplaceAll(split[1], "_", "-")
			optDeployment = strings.ReplaceAll(split[2], "_", "-")
			optEnv = strings.ReplaceAll(split[3], "_", "-")
			log.Println("从 Jenkins $JOB_NAME 获取到 集群 " + optCluster + "，命名空间 " + optNamespace + "，工作负载 " + optDeployment + ", 环境 " + optEnv)
		}
	}

	// extract $HOME
	optHome = os.Getenv("HOME")

	// check options
	checkStr(&optRegistry, "错误：镜像仓库未指定，使用 --registry 指定镜像仓库")
	checkStr(&optCluster, "错误：集群未指定，使用 --cluster 指定集群")
	checkStr(&optNamespace, "错误：命名空间未指定，使用 --namespace 指定命名空间")
	checkStr(&optDeployment, "错误: 工作负载未指定，使用 --deployment，或者 $JOB_NAME 指定工作负载")
	checkStr(&optEnv, "错误：环境未指定，使用 --env 或者 $JOB_NAME 指定环境")
	checkStr(&optHome, "错误：$HOME 环境变量未指定，无法获取 $HOME/.kube/config-XXX 文件")

	// calculate image name, docker-build, Dockerfile
	imageName := fmt.Sprintf("%s/%s:%s", optRegistry, optNamespace+"-"+optDeployment, optEnv)
	log.Println("镜像名称 " + imageName)
	log.Println("部署目标 " + optCluster + " > " + optNamespace + " > " + optDeployment)
	log.Println("------------")

	// check docker-build.XXX.sh
	buildFile := fmt.Sprintf("docker-build.%s.sh", optEnv)
	if err = checkFile(buildFile, true); err != nil {
		return
	}

	// check Dockerfile.XXX
	dockerFile := fmt.Sprintf("Dockerfile.%s", optEnv)
	if err = checkFile(dockerFile, false); err != nil {
		return
	}

	// check docker executable
	if _, err = run(false, false, "docker", "--version"); err != nil {
		return
	}

	// check kubectl executable
	if _, err = run(false, false, "kubectl", "version", "--client"); err != nil {
		return
	}

	// execute docker-build.XXX.sh
	if _, err = run(true, false, buildFile); err != nil {
		return
	}

	// execute docker build
	if _, err = run(false, false, "docker", "build", "-t", imageName, "-f", dockerFile, "."); err != nil {
		return
	}

	// execute docker push
	if _, err = run(false, false, "docker", "push", imageName); err != nil {
		return
	}

	// execute docker inspect
	var canonicalName string
	if canonicalName, err = run(false, true, "docker", "inspect", "--format", "{{index .RepoDigests 0}}", imageName); err != nil {
		return
	}
	log.Println("完整镜像 " + canonicalName)

	// check kubectl config
	kubeconfig := filepath.Join(optHome, ".kube", "config-"+optCluster)
	if err = checkFile(kubeconfig, false); err != nil {
		return
	}

	// execute kubectl
	if _, err = run(false, false,
		"kubectl", "--kubeconfig", kubeconfig, "--namespace", optNamespace, "set", "image", "deployment/"+optDeployment, optDeployment+"="+canonicalName); err != nil {
		return
	}
}
