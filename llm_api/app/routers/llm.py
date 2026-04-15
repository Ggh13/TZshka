from fastapi import APIRouter, HTTPException, status, Depends, UploadFile, File
from langchain_core.exceptions import OutputParserException
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

    except OutputParserException:
        raise HTTPException(status_code=status.HTTP_422_UNPROCESSABLE_CONTENT,
                            detail="Invalid LLM response format")

    except TimeoutError:
        raise HTTPException(status_code=status.HTTP_504_GATEWAY_TIMEOUT,
                            detail="LLM timeout")

    except Exception:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                            detail="LLM service unavailable")

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

    user_content = (await file.read()).decode("utf-8")

    try:
        llm_output = LLM.invoke({
            "technical_specification": user_content
        })

    except OutputParserException:
        raise HTTPException(status_code=status.HTTP_422_UNPROCESSABLE_CONTENT,
                            detail="Invalid LLM response format")

    except TimeoutError:
        raise HTTPException(status_code=status.HTTP_504_GATEWAY_TIMEOUT,
                            detail="LLM timeout")

    except Exception:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                            detail="LLM service unavailable")

    return llm_output