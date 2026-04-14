from pydantic import BaseModel, Field


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