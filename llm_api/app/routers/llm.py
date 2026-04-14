from fastapi import APIRouter, HTTPException, status
from app.schemas.text_message import TextInput as TextInputSchema
from app.schemas.llm_output import LLMOutput as LLMOutputSchema
from app.llm.chain import CHAIN_TEXT_HANDLER as LLM


router = APIRouter(
    prefix='/llm',
    tags=['LLM']
)


@router.post('/text', response_model=LLMOutputSchema)
async def get_response_by_text(
    payload: TextInputSchema
):
    user_content = payload.content
    try:
        llm_output = LLM.invoke(user_content)
    except:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="LLM is not available")

    return llm_output