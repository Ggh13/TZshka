from fastapi import APIRouter, status
from fastapi.responses import JSONResponse

from langchain_core.exceptions import OutputParserException
from pydantic import ValidationError

from app.schemas.llm import LLMResponse as LLMResponseSchema, LLMCombinedOutput as LLMCombinedOutputSchema
from app.schemas.error import Error as ErrorSchema

from app.llm.chains import CHAIN_TEXT_HANDLER as RULES_CHECKER
from app.llm.chains import CHAIN_STANDARD_HANDLER as STANDARD_CHECKER


router = APIRouter(
    prefix='/llm',
    tags=['LLM']
)


@router.post('/text', response_model=LLMResponseSchema)
async def get_response_by_text(
    mode: str = "Instant",
    standard: str = None,
    content: str = "",
):
    user_content = content
    try:
        if RULES_CHECKER is None:
            return JSONResponse(
                content=LLMResponseSchema(
                    success=False,
                    code=status.HTTP_503_SERVICE_UNAVAILABLE,
                    data=None,
                    error=ErrorSchema(
                        type="service_unavailable",
                        message="LLM not configured"
                    )
                ).model_dump()
            )

        rules_output = RULES_CHECKER.invoke({
            "technical_specification": user_content
        })

        standard_output = None

        if standard is not None:
            standard_output = STANDARD_CHECKER.invoke({
                "standard": standard,
                "technical_specification": user_content
            })

        return LLMResponseSchema(
            success=True,
            code=status.HTTP_200_OK,
            data=LLMCombinedOutputSchema(
                rules_checker=rules_output,
                standard_checker=standard_output
            ),
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

    except Exception as e:
        import traceback
        return JSONResponse(
            content=LLMResponseSchema(
                success=False,
                code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                data=None,
                error=ErrorSchema(
                    type="internal_error",
                    message=f"LLM service unavailable: {str(e)}"
                )
            ).model_dump()
        )