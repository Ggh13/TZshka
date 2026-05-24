import json

from app.llm.chains import (
    CHAIN_CONSISTENCY_AGGREGATOR,
    CHAIN_STANDARD_HANDLER,
    CHAIN_STANDARD_HANDLER_THINKING,
    CHAIN_TEXT_HANDLER,
    CHAIN_TEXT_HANDLER_THINKING,
)
from app.schemas.llm import LLMCombinedOutput as LLMCombinedOutputSchema


SELF_CONSISTENCY_RUNS = 3


def check_text_instant(
        technical_specification: str,
        standard: str | None
) -> LLMCombinedOutputSchema:
    rules_output = CHAIN_TEXT_HANDLER.invoke({
        "technical_specification": technical_specification
    })

    standard_output = None

    if standard is not None:
        standard_output = CHAIN_STANDARD_HANDLER.invoke({
            "standard": standard,
            "technical_specification": technical_specification
        })

    return LLMCombinedOutputSchema(
        rules_checker=rules_output,
        standard_checker=standard_output
    )


def check_text_thinking(
        technical_specification: str,
        standard: str | None
) -> LLMCombinedOutputSchema:
    rules_output = _run_self_consistency(
        chain=CHAIN_TEXT_HANDLER_THINKING,
        technical_specification=technical_specification,
        payload={
            "technical_specification": technical_specification
        }
    )

    standard_output = None

    if standard is not None:
        standard_output = _run_self_consistency(
            chain=CHAIN_STANDARD_HANDLER_THINKING,
            technical_specification=technical_specification,
            payload={
                "standard": standard,
                "technical_specification": technical_specification
            }
        )

    return LLMCombinedOutputSchema(
        rules_checker=rules_output,
        standard_checker=standard_output
    )


def check_text(
        mode: str,
        technical_specification: str,
        standard: str | None
) -> LLMCombinedOutputSchema:
    if mode == "Thinking":
        return check_text_thinking(
            technical_specification=technical_specification,
            standard=standard
        )

    return check_text_instant(
        technical_specification=technical_specification,
        standard=standard
    )


def _run_self_consistency(
        chain,
        technical_specification: str,
        payload: dict
) -> dict:
    checker_outputs = [
        chain.invoke(payload)
        for _ in range(SELF_CONSISTENCY_RUNS)
    ]

    return CHAIN_CONSISTENCY_AGGREGATOR.invoke({
        "technical_specification": technical_specification,
        "checker_outputs": json.dumps(
            checker_outputs,
            ensure_ascii=False,
            default=str
        )
    })
