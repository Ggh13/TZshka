from fastapi import APIRouter, HTTPException, status
from app.schemas.text_message import TextInput as TextInputSchema
from app.schemas.llm_output import LLMOutput as LLMOutputSchema
from app.llm.model import llm


router = APIRouter(
    prefix='/llm',
    tags=['LLM']
)


@router.get('/text', response_model=LLMOutputSchema)
async def get_response_by_text(
    payload: TextInputSchema
):
    user_content = payload.content
    try:
        llm_output = llm.invoke(user_content)
    except:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="LLM is not available")

    return llm_output