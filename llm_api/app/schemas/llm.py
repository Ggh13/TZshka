from pydantic import BaseModel, Field
from app.schemas.error import Error as ErrorSchema


class LLMOutput(BaseModel):
    status: str = Field(..., description="valid | issues_found")
    issues: list[dict] = Field(..., description="""
    Формат ответа в случае status == issues_found:
    {
      "rule_id": "...",
      "problem": "...",
      "explanation": "..."
    }
    """)
    feedback: str = Field(..., description="Объективная и справедливая оценка технического задания")


class LLMResponse(BaseModel):
    success: bool = Field(..., description="Успех/не успех запроса")
    code: int = Field(..., description="HTTP status code")
    data: LLMOutput | None = Field(
        default=None,
        description="Полезные данные при успешном ответе"
    )
    error: ErrorSchema | None = Field(
        default=None,
        description="Информация об ошибке"
    )
