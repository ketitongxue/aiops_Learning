作业
1. 实现 Function Calling

定义 modify_config 函数，入参：service_name，key，value
定义 restart_service 函数，入参：service_name
定义 apply_manifest 函数，入参：resource_type，image


2. 实践 Function Calling，观察以下输入是否能正确选择对应的函数

帮我修改 gateway 的配置，vendor 修改为 alipay
帮我重启 gateway 服务
帮我部署一个 deployment，镜像是 nginx

```python
python main.py

query: 帮我修改 gateway 的配置，vendor 修改为 alipay

ChatGPT want to call function:  [ChatCompletionMessageToolCall(id='call_SM0e5HpodWyabghHEyfoyvVz', function=Function(arguments='{"service_name":"gateway","key":"vendor","value":"alipay"}', name='modify_config'), type='function')]
query: 帮我重启 gateway 服务

ChatGPT want to call function:  [ChatCompletionMessageToolCall(id='call_wNKFw36SMZ2L4gblaGS5VOpl', function=Function(arguments='{"service_name":"gateway"}', name='restart_service'), type='function')]
query: 帮我部署一个 deployment，镜像是 nginx

ChatGPT want to call function:  [ChatCompletionMessageToolCall(id='call_oBzq1nNUG3BrXVAjWhLyc4zm', function=Function(arguments='{"resource_type":"deployment","image":"nginx"}', name='apply_manifest'), type='function')]
```