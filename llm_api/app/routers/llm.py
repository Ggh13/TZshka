from fastapi import APIRouter, status, Depends
from fastapi.responses import JSONResponse

from langchain_core.exceptions import OutputParserException
from pydantic import ValidationError

from app.schemas.text_message import TextInput as TextInputSchema
from app.schemas.llm import LLMResponse as LLMResponseSchema
from app.schemas.error import Error as ErrorSchema
from app.schemas.mode import Mode as ModeSchema
from app.schemas.standard import Standard as StandardSchema

from app.llm.chains import CHAIN_TEXT_HANDLER as RULES_CHECKER


router = APIRouter(
    prefix='/llm',
    tags=['LLM']
)


@router.post('/text', response_model=LLMResponseSchema)
async def get_response_by_text(
    mode: ModeSchema,
    standard: StandardSchema,
    text: TextInputSchema = Depends(TextInputSchema.as_form),
):
    user_content = text.content
    try:
        llm_output = RULES_CHECKER.invoke({
            "technical_specification": user_content
        })

        return LLMResponseSchema(
            success=True,
            code=status.HTTP_200_OK,
            data=llm_output,
            error=None
        )

    except ValidationError as e:
        return JSONResponse(
            content=LLMResponseSchema(
                success=False,
                code=status.HTTP_422_UNPROCESSABLE_CONTENT,
                data=None,
                error=ErrorSchema(
                    type="validation_error",
                    message=str(e)
                )
            ).model_dump()
        )

    except OutputParserException:
        return JSONResponse(
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
