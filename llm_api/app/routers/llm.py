from fastapi import APIRouter, HTTPException, status, Depends, UploadFile, File
from app.schemas.text_message import TextInput as TextInputSchema
from app.schemas.llm_output import LLMOutput as LLMOutputSchema
from app.llm.chain import CHAIN_TEXT_HANDLER as LLM


router = APIRouter(
    prefix='/llm',
    tags=['LLM']
)


@router.post('/text', response_model=LLMOutputSchema)
async def get_response_by_text(
    payload: TextInputSchema = Depends(TextInputSchema.as_form)
):
    user_content = payload.content
    try:
        llm_output = LLM.invoke({
            "technical_specification": user_content
        })
    except:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="LLM is not available")

    return llm_output


@router.post('/text_file')
async def get_response_by_text_file(
        file: UploadFile = File(...)
):
    if file.content_type != "text/plain":
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="File must be .txt"
        )

    content = (await file.read()).decode("utf-8")

    try:
        llm_output = LLM.invoke({
            "technical_specification": content
        })
    except Exception:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="LLM is not available"
        )

    return llm_output