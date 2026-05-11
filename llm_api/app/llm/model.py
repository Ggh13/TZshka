import os
from dotenv import load_dotenv
from langchain_openrouter import ChatOpenRouter

from app.schemas.llm import LLMOutput as LLMOutputSchema

load_dotenv()

base_llm = ChatOpenRouter(
    # model="meta-llama/llama-3.3-70b-instruct",
    model="qwen/qwen3-next-80b-a3b-instruct",
    temperature=0,
    max_tokens=11865,
    api_key=os.getenv("OPENROUTER_API_KEY")
)

thinking_base_llm = ChatOpenRouter(
    model="qwen/qwen3-next-80b-a3b-instruct",
    temperature=0.5,
    max_tokens=11865,
    api_key=os.getenv("OPENROUTER_API_KEY")
)

llm = base_llm.with_structured_output(LLMOutputSchema)
thinking_llm = thinking_base_llm.with_structured_output(LLMOutputSchema)
