from fastapi import FastAPI

app = FastAPI(
    title="API для обращения к LLM"
)


@app.get('/')
async def root():
    return {'message': 'success'}