from pydantic import BaseModel, Field, ConfigDict


class TextInput(BaseModel):
    content: str = Field(min_length=1, max_length=10000, description="Текст технического задания")

    model_config = ConfigDict(from_attributes=True)