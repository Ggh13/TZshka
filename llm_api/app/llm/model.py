import os
from dotenv import load_dotenv
from langchain_openrouter import ChatOpenRouter

from app.schemas.llm import LLMOutput as LLMOutputSchema

load_dotenv()

llm = ChatOpenRouter(
    model="meta-llama/llama-3.3-70b-instruct",
    temperature=0,
    api_key=os.getenv("OPENROUTER_API_KEY")
)

llm = llm.with_structured_output(LLMOutputSchema)