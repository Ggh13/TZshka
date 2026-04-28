from pydantic import BaseModel, Field

class Error(BaseModel):
    type: str = Field(..., description="Тип ошибки")
    message: str = Field(..., description="Описание ошибки")
