from fastapi import APIRouter, HTTPException, status, Depends, UploadFile, File
from fastapi.responses import JSONResponse

from langchain_core.exceptions import OutputParserException

from app.schemas.text_message import TextInput as TextInputSchema
from app.schemas.llm import LLMOutput as LLMOutputSchema, LLMResponse as LLMResponseSchema
from app.schemas.error import Error as ErrorSchema
from app.llm.chain import CHAIN_TEXT_HANDLER as LLM


router = APIRouter(
    prefix='/llm',
    tags=['LLM']
)


@router.post('/text', response_model=LLMResponseSchema)
async def get_response_by_text(
    payload: TextInputSchema = Depends(TextInputSchema.as_form)
):
    user_content = payload.content
    try:
        llm_output = LLM.invoke({
            "technical_specification": user_content
        })

        return LLMResponseSchema(
            success=True,
            code=status.HTTP_200_OK,
            data=llm_output,
            error=None
        )

    except OutputParserException:
        return JSONResponse(
            status_code=status.HTTP_422_UNPROCESSABLE_CONTENT,
            content=LLMResponseSchema(
                success=False,
                code=status.HTTP_422_UNPROCESSABLE_CONTENT,
                data=None,
                error=ErrorSchema(
                    type="invalid_llm_response",
                    message="Invalid LLM response format"
                )
            ).model_dump()
        )

    except TimeoutError:
        return JSONResponse(
            status_code=status.HTTP_504_GATEWAY_TIMEOUT,
            content=LLMResponseSchema(
                success=False,
                code=status.HTTP_504_GATEWAY_TIMEOUT,
                data=None,
                error=ErrorSchema(
                    type="timeout",
                    message="LLM timeout"
                )
            ).model_dump()
        )

    except Exception:
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=LLMResponseSchema(
                success=False,
                code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                data=None,
                error=ErrorSchema(
                    type="internal_error",
                    message="LLM service unavailable"
                )
            ).model_dump()
        )


@router.post('/text_file', response_model=LLMResponseSchema)
async def get_response_by_text_file(file: UploadFile = File(...)):
    if file.content_type != "text/plain":
        return JSONResponse(
            status_code=status.HTTP_400_BAD_REQUEST,
            content=LLMResponseSchema(
                success=False,
                code=status.HTTP_400_BAD_REQUEST,
                data=None,
                error=ErrorSchema(
                    type="invalid_file_type",
                    message="File must be .txt"
                )
            ).model_dump()
        )

    try:
        user_content = (await file.read()).decode("utf-8")
    except UnicodeDecodeError:
        return JSONResponse(
            status_code=status.HTTP_400_BAD_REQUEST,
            content=LLMResponseSchema(
                success=False,
                code=status.HTTP_400_BAD_REQUEST,
                data=None,
                error=ErrorSchema(
                    type="invalid_file_encoding",
                    message="File must be UTF-8 encoded text"
                )
            ).model_dump()
        )

    try:
        llm_output = LLM.invoke({
            "technical_specification": user_content
        })

        return LLMResponseSchema(
            success=True,
            code=status.HTTP_200_OK,
            data=llm_output,
            error=None
        )

    except OutputParserException:
        return JSONResponse(
            status_code=status.HTTP_422_UNPROCESSABLE_CONTENT,
            content=LLMResponseSchema(
                success=False,
                code=status.HTTP_422_UNPROCESSABLE_CONTENT,
                data=None,
                error=ErrorSchema(
                    type="invalid_llm_response",
                    message="Invalid LLM response format"
                )
            ).model_dump()
        )

    except TimeoutError:
        return JSONResponse(
            status_code=status.HTTP_504_GATEWAY_TIMEOUT,
            content=LLMResponseSchema(
                success=False,
                code=status.HTTP_504_GATEWAY_TIMEOUT,
                data=None,
                error=ErrorSchema(
                    type="timeout",
                    message="LLM timeout"
                )
            ).model_dump()
        )

    except Exception:
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=LLMResponseSchema(
                success=False,
                code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                data=None,
                error=ErrorSchema(
                    type="internal_error",
                    message="LLM service unavailable"
                )
            ).model_dump()
        )