from fastapi import FastAPI
from app.routers import llm

app = FastAPI(
    title="API для обращения к LLM"
)

app.include_router(llm.router)

@app.get('/')
async def root():
    return {'message': 'success'}