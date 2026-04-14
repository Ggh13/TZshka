from langchain_core.prompts import (
    SystemMessagePromptTemplate,
    HumanMessagePromptTemplate,
    ChatPromptTemplate
)
from app.llm.model import llm

from pathlib import Path

rules_path = Path(__file__).parent.parent / 'static_files' / 'rules.txt'

with open(rules_path, 'r', encoding='utf-8') as file:
    rules = file.read()

sp = """
Ты полезный ассистент. Твоя роль заключается в том, чтобы анализировать 
технические задания пользователей и находить в них несоответствия правилам 
составления технических заданий. 
Все правила составления хорошего технического задания описаны в списке: 
""" + rules + """
Ты должен проанализировать техническое задание и проверить его на соответствие 
КАЖДОМУ правилу из списка. 
В случае возникновения противоречия какому либо правилу составления хорошего 
технического задания ты должен вывести сообщение с правилом, которое нарушено, 
и частью (частями) текста, в котором есть противоречие. Учти, что таких нарушений 
может быть насколько в одном техническом задании пользователя. Также ты должен 
дать объективную оценку качества технического задания.
"""

system_prompt = SystemMessagePromptTemplate.from_template(sp)

user_prompt = HumanMessagePromptTemplate.from_template(
    """
    Техническое задание для анализа:

    {technical_specification}

    """,
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