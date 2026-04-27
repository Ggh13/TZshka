from pydantic import BaseModel, Field, ConfigDict


class TextInput(BaseModel):
    content: str = Field(min_length=1, description="Текст технического задания")

    model_config = ConfigDict(from_attributes=True)