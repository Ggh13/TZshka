from pydantic import BaseModel, Field, ConfigDict
from fastapi import Form


class TextInput(BaseModel):
    content: str = Field(min_length=1, description="Текст технического задания")

    @classmethod
    def as_form(
            cls,
            content: str = Form(...)
    ):
        return cls(content=content)

    model_config = ConfigDict(from_attributes=True)