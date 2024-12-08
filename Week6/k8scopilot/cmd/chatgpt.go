/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ketitongxue/k8scopilot/utils"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/restmapper"
	"k8s.io/kubectl/pkg/scheme"
)

// chatgptCmd represents the chatgpt command
var chatgptCmd = &cobra.Command{
	Use:   "chatgpt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		startChat()
	},
}

func startChat() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("我是 K8s Copilot，请问有什么可以帮助你？")

	for {
		fmt.Print("> ")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" {
				fmt.Println("Bye!")
				break
			}
			if input == "" {
				continue
			}
			response := processInput(input)
			fmt.Println(response)
		}
	}
}

func processInput(input string) string {
	client, err := utils.NewOpenAIClient()
	if err != nil {
		return err.Error()
	}
	// response, err := client.SendMessage("你是一个 K8s Copilot，你要帮用户生成 YAML 文件，除了 YAML 内容以外不要输出任何内容，此外不要把 YAML 放在 ``` 代码块里", input)
	// if err != nil {
	// 	return err.Error()
	// }
	// return response
	response := functionCalling(input, client)
	return response
}

func functionCalling(input string, client *utils.OpenAI) string {
	//定义第一个函数，生成 K8s YAML，并部署资源
	f1 := openai.FunctionDefinition{
		Name:        "generateAndDeployK8sResource",
		Description: "生成 K8s YAML 文件，并部署资源",
		//Parameters 表示函数具体接收参数
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"user_input": {
					Type:        jsonschema.String,
					Description: "用户输入的内容，要求包含资源类型和镜像",
				},
			},
			Required: []string{"user_input"},
		},
	}

	t1 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f1,
	}

	// 定义查询 K8s 资源
	f2 := openai.FunctionDefinition{
		Name:        "queryK8sResource",
		Description: "查询 K8s 资源",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"namespace": {
					Type:        jsonschema.String,
					Description: "资源所在命名空间",
				},
				"resource_type": {
					Type:        jsonschema.String,
					Description: "K8s 标准资源类型，可以是：Pod、Deployment、Service、Ingress、Secret、ConfigMap、PersistentVolumeClaim、PersistentVolume",
				},
			},
			Required: []string{"namespace", "resource_type"},
		},
	}

	t2 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f2,
	}

	// 定义删除 K8s 资源
	f3 := openai.FunctionDefinition{
		Name:        "deleteK8sResource",
		Description: "删除 K8s 资源",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"namespace": {
					Type:        jsonschema.String,
					Description: "资源所在命名空间",
				},
				"resource_type": {
					Type:        jsonschema.String,
					Description: "K8s 标准资源类型，可以是：Pod、Deployment、Service、Ingress、Secret、ConfigMap、PersistentVolumeClaim、PersistentVolume",
				},
				"resource_name": {
					Type:        jsonschema.String,
					Description: "K8s 资源名称",
				},
			},
			Required: []string{"namespace", "resource_type", "resource_name"},
		},
	}

	t3 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f3,
	}

	dialog := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: input},
	}

	resp, err := client.Client.CreateChatCompletion(context.TODO(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: dialog,
			Tools:    []openai.Tool{t1, t2, t3},
		},
	)
	if err != nil {
		return err.Error()
	}

	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) != 1 {
		return fmt.Sprintf("未找到合适的工具调用，%v", len(msg.ToolCalls))
	}

	//组装对话历史
	dialog = append(dialog, msg)
	// return fmt.Sprintf("OpenAI 希望能请求函数 %s，参数：%s", msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	result, err := callFunction(client, msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	if err != nil {
		return fmt.Sprintf("Error calling function: %v\n", err)
	}
	return result
}

// 根据 OpenAI 返回的消息，调用对应的函数
func callFunction(client *utils.OpenAI, name, arguments string) (string, error) {
	if name == "generateAndDeployK8sResource" {
		params := struct {
			UserInput string `json:"user_input"`
		}{}
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return "", err
		}
		return generateAndDeployK8sResource(client, params.UserInput)
	}
	if name == "queryK8sResource" {
		params := struct {
			Namespace    string `json:"namespace"`
			ResourceType string `json:"resource_type"`
		}{}
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return "", err
		}
		return queryK8sResource(params.Namespace, params.ResourceType)
	}
	if name == "deleteK8sResource" {
		params := struct {
			Namespace    string `json:"namespace"`
			ResourceType string `json:"resource_type"`
			ResourceName string `json:"resource_name"`
		}{}
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return "", err
		}
		return deleteK8sResource(params.Namespace, params.ResourceType, params.ResourceName)
	}
	return "", fmt.Errorf("unknown function: %s", name)
}

func generateAndDeployK8sResource(client *utils.OpenAI, userInput string) (string, error) {
	yamlContent, err := client.SendMessage("你是一个 K8s Copilot，你要帮用户生成 YAML 文件，除了 YAML 内容以外不要输出任何内容，此外不要把 YAML 放在 ``` 代码块里", userInput)
	if err != nil {
		return "", err
	}
	// return yamlContent, nil
	// 调用 dynamic client 创建资源
	clientGo, err := utils.NewClientGo(kubeconfig)
	if err != nil {
		return "", err
	}
	resource, err := restmapper.GetAPIGroupResources(clientGo.DiscoveryClient)
	if err != nil {
		return "", err
	}
	fmt.Println(resource)
	// 把 YAML 转换成 Unstructured
	unstructuredObj := &unstructured.Unstructured{}
	_, _, err = scheme.Codecs.UniversalDeserializer().Decode([]byte(yamlContent), nil, unstructuredObj)
	if err != nil {
		return "", err
	}
	// 创建 mapper
	mapper := restmapper.NewDiscoveryRESTMapper(resource)
	// 从 unstructuredObj 中提取 GVK
	gvk := unstructuredObj.GroupVersionKind()
	// 用 GVK 转 GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return "", err
	}

	namespace := unstructuredObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	// 使用 Dynamic 创建资源
	_, err = clientGo.DynamicClient.Resource(mapping.Resource).Namespace(namespace).Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("YAML content:\n%s\n\nDeployment successful.", yamlContent), nil
}

func queryK8sResource(namespace, resourceType string) (string, error) {
	// 调用 dynamic client 查询资源
	clientGo, err := utils.NewClientGo(kubeconfig)
	if err != nil {
		return "", err
	}
	resourceType = strings.ToLower(resourceType)
	var gvr schema.GroupVersionResource
	switch resourceType {
	case "deployment":
		gvr = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	case "service":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	case "pod":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	default:
		return "", fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	// Query the resources using the dynamic client
	resourceList, err := clientGo.DynamicClient.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	result := ""
	for _, item := range resourceList.Items {
		result += fmt.Sprintf("Found %s: %s\n", resourceType, item.GetName())
	}

	return result, nil
}

func deleteK8sResource(namespace, resourceType string, resourceName string) (string, error) {
	// 调用 dynamic client 删除资源
	clientGo, err := utils.NewClientGo(kubeconfig)
	if err != nil {
		return "", err
	}
	resourceType = strings.ToLower(resourceType)

	var gvr schema.GroupVersionResource
	switch resourceType {
	case "deployment":
		gvr = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	case "service":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	case "pod":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	default:
		return "", fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	err = clientGo.DynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%sResource %s in %s deleted successed\n", resourceType, resourceName, namespace), nil
}

func init() {
	askCmd.AddCommand(chatgptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatgptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatgptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
