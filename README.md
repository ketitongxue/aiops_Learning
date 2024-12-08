极客时间 AIOps 训练营作业
Week1
使用 Terraform 创建腾讯云虚拟机并安装 Docker

Week3
1. 实现 Function Calling

定义 modify_config 函数，入参：service_name，key，value
定义 restart_service 函数，入参：service_name
定义 apply_manifest 函数，入参：resource_type，image

2. 实践 Function Calling，观察以下输入是否能正确选择对应的函数

帮我修改 gateway 的配置，vendor 修改为 alipay
帮我重启 gateway 服务
帮我部署一个 deployment，镜像是 nginx


Week6
完善 chatgpt.go，实现 deleteResource 方法，能其能以对话的方式删除 K8s 资源。

```bash
./k8scopilot ask chatgpt
我是 K8s Copilot，请问有什么可以帮助你？
> 帮我部署一个 deploy，镜像是nginx
YAML content:
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80

Deployment successful.

// 进入终端查询
$ kubectl get deploy
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   0/3     3            0           8s


> 查询default ns下的deploy
Found deployment: nginx-deployment

> 查询 default ns 下的 pod
Found pod: nginx-deployment-cd55c47f5-chrpd
Found pod: nginx-deployment-cd55c47f5-fxkkc
Found pod: nginx-deployment-cd55c47f5-nr622

> 删除 default ns 下的 名字为 nginx-deployment 的deploy
deploymentResource nginx-deployment in default deleted successed

// 进入终端查询
$ kubectl get deploy
No resources found in default namespace.
```