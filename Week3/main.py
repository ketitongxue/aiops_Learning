# pip install openai
from openai import OpenAI

# 定义 modify_config 函数，入参：service_name，key，value
def modify_config(service_name, key, value):
    print(f"service_name:{service_name},key:{key},value:{value}")

# 定义 restart_service 函数，入参：service_name
def restart_service(service_name):
    print(f"service_name:{service_name}")

# 定义 apply_manifest 函数，入参：resource_type，image
def apply_manifest(resource_type, image):
    print(f"resource_type:{resource_type},image:{image}")

    
def main(query):
    # Initialize the OpenAI client
    client = OpenAI()
    messages = [
        {
            "role": "system",
            "content": "你是一个 Kubernetes 集群管理助手，你可以帮助用户管理 Kubernetes 集群，你可以调用多个函数来帮助用户完成任务",
        },
        {
            "role": "user",
            "content": query,
        },
    ]
    tools = [
        {
            "type": "function",
            "function": {
                "name": "modify_config",
                "description": "修改配置",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "service_name": {"type": "string"},
                        "key": {"type": "string"},
                        "value": {"type": "string"},
                    },
                },
            },
        },
        {
            "type": "function",
            "function": {
                "name": "restart_service",
                "description": "重启服务",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "service_name": {"type": "string"},
                    },
                },
            },
        },
        {
            "type": "function",
            "function": {
                "name": "apply_manifest",
                "description": "部署应用",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "resource_type": {"type": "string"},
                        "image": {"type": "string"},
                    },
                },
            },
        },
    ]

    response = client.chat.completions.create(
        model="gpt-4o-mini",
        messages=messages,
        tools=tools,
        tool_choice="auto",
    )

    response_message = response.choices[0].message
    tool_calls = response_message.tool_calls

    print("\nChatGPT want to call function: ", tool_calls)

if __name__ == "__main__":
    query_List = ["帮我修改 gateway 的配置，vendor 修改为 alipay","帮我重启 gateway 服务","帮我部署一个 deployment，镜像是 nginx"]
    for query in query_List:
        print(f"query: {query}")
        main(query)
