import os
from dotenv import load_dotenv
from langchain_openrouter import ChatOpenRouter

from app.schemas.llm import LLMOutput as LLMOutputSchema

load_dotenv()

api_key = os.getenv("OPENROUTER_API_KEY")

llm = None
if api_key:
    llm = ChatOpenRouter(
        model="qwen/qwen3-next-80b-a3b-instruct",
        temperature=0,
        api_key=api_key
    )
    llm = llm.with_structured_output(LLMOutputSchema)