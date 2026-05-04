from langchain_core.prompts import (
    SystemMessagePromptTemplate,
    HumanMessagePromptTemplate,
    ChatPromptTemplate
)
from app.llm.model import llm
from app.llm.prompts import (
    USER_PROMPT_TEMPLATE,
    get_standard_handler,
    get_text_handler
)

from pathlib import Path

rules_path = Path(__file__).parent.parent / 'static_files' / 'rules.txt'
gost_19_path = Path(__file__).parent.parent / 'static_files' / 'GOST-19.txt'
gost_34_path = Path(__file__).parent.parent / 'static_files' / 'GOST-34.txt'

with open(rules_path, 'r', encoding='utf-8') as file:
    rules = file.read()

with open(gost_19_path, 'r', encoding='utf-8') as file:
    gost_19 = file.read()

with open(gost_34_path, 'r', encoding='utf-8') as file:
    gost_34 = file.read()

system_prompt = SystemMessagePromptTemplate.from_template(get_text_handler(rules))

user_prompt = HumanMessagePromptTemplate.from_template(
    USER_PROMPT_TEMPLATE,
    input_variables=["technical_specification"]
)

chat_prompt = ChatPromptTemplate.from_messages([system_prompt, user_prompt])

CHAIN_TEXT_HANDLER = (
    {
        'technical_specification': lambda x: x['technical_specification']
    }
    | chat_prompt
    | llm
    | {
        "status": lambda x: x.status,
        "issues": lambda x: x.issues,
        "feedback": lambda x: x.feedback
    }
)

standard_chat_prompt = ChatPromptTemplate.from_messages([
    ("system", get_standard_handler()),
    user_prompt
])

CHAIN_STANDARD_HANDLER = (
    {
        "standard": lambda x: x["standard"],
        "technical_specification": lambda x: x["technical_specification"],
    }
    | standard_chat_prompt
    | llm
    | {
        "status": lambda x: x.status,
        "issues": lambda x: x.issues,
        "feedback": lambda x: x.feedback,
    }
)
