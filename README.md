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