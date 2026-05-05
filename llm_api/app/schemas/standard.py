from pydantic import BaseModel, field_validator
from typing import Literal
from fastapi import Form


class Standard(BaseModel):
    standard: Literal["ГОСТ-19", "ГОСТ-34"] | None = None


    @field_validator("standard")
    @classmethod
    def validate_value(cls, v):
        allowed = {"ГОСТ-19", "ГОСТ-34"}
        if v is not None and v not in allowed:
            raise ValueError("Standard must be 'ГОСТ-19', 'ГОСТ-34' or None")
        return v

    @classmethod
    def as_form(
            cls,
            standard: Literal["ГОСТ-19", "ГОСТ-34"] | None = Form(None),
    ):
        return cls(standard=standard)