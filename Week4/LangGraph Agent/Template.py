from langchain_openai import ChatOpenAI

llm = ChatOpenAI()

class State(TypedDict):
    message: Annotated[list, add_messages]

def chat(State: State):
    return {"message":[llm.invoke(State["messages"])]}

# 单节点
workflow = StateGraph(State)
workflow.add_node(chat)
workflow.set_entry_point("chat")
workflow.set_finish_point("chat")
graph = workflow.compile()
