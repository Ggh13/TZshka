from pydantic import BaseModel, field_validator
from typing import Literal


class Mode(BaseModel):
    mode: Literal["Instant", "Thinking"] = "Instant"


    @field_validator("mode")
    @classmethod
    def validate_value(cls, v):
        allowed = {"Instant", "Thinking"}
        if v is not None and v not in allowed:
            raise ValueError("Mode value must be 'Instant' or 'Thinking'")
        return v