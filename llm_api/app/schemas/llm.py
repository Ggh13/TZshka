from pydantic import BaseModel, Field
from app.schemas.error import Error as ErrorSchema



class Issue(BaseModel):
    rule_id: str = Field(..., description="Идентификатор или краткое название нарушенного правила")
    problem: str = Field(..., description="Проблемный фрагмент ТЗ или точное описание нарушения")
    explanation: str = Field(..., description="Почему это является нарушением правила")

class LLMOutput(BaseModel):
    status: str = Field(..., description="valid | issues_found")
    issues: list[Issue] = Field(default_factory=list, description="Список найденных нарушений")
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
