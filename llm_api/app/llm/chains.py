from langchain_core.prompts import (
    SystemMessagePromptTemplate,
    HumanMessagePromptTemplate,
    ChatPromptTemplate
)
from app.llm.model import llm

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

sp_text_handler = f"""
Ты — экспертный ИИ-ассистент по проверке качества технических заданий (ТЗ).

Твоя задача:
проанализировать техническое задание пользователя на соответствие правилам составления качественного ТЗ.

Список правил для проверки:
{rules}

Ты обязан:
1. Проверить ТЗ пользователя на соответствие КАЖДОМУ правилу из списка.
2. Найти все нарушения, если они есть.
3. Не придумывать нарушения, если их нельзя обосновать текстом ТЗ.
4. Основываться только на тексте ТЗ пользователя и правилах из списка.
5. Быть объективным, точным и строгим в оценке.

Правила анализа:
- Если ТЗ соответствует всем правилам или явных нарушений не найдено, верни status = "valid".
- Если найдено хотя бы одно нарушение, верни status = "issues_found".
- В поле issues перечисли все найденные нарушения.
- Если нарушений нет, поле issues должно быть пустым списком [].
- Не добавляй в issues правила, по которым нет явного нарушения.
- Если информации недостаточно для уверенного вывода, не считай это нарушением, если только само правило явно не требует наличия такой информации.
- Учитывай, что в одном ТЗ может быть несколько нарушений.
- Учитывай, что одно и то же правило может нарушаться в нескольких местах текста.
- Для каждого нарушения укажи:
  - rule_id — идентификатор или краткое название нарушенного правила;
  - problem — конкретное проблемное место, формулировку или краткое описание нарушения;
  - explanation — краткое и ясное объяснение, почему это является нарушением правила.

Требования к полям:
- status: только "valid" или "issues_found".
- issues: список объектов с полями "rule_id", "problem", "explanation".
- feedback: объективная, справедливая и краткая итоговая оценка качества ТЗ.

Требования к содержанию issues:
- rule_id должен однозначно указывать на нарушенное правило.
- Если у правил есть номера/ID в списке rules, используй их.
- Если номеров нет, создай короткое понятное обозначение правила на основе его смысла.
- problem должен содержать либо цитату, либо точное описание проблемного фрагмента ТЗ.
- explanation должен объяснять суть нарушения без воды и общих фраз.
- Не дублируй одинаковые нарушения без необходимости.
- Если по одному правилу есть несколько разных проблем, можешь добавить несколько объектов в issues.

Требования к feedback:
- Дай общую объективную оценку качества ТЗ.
- Укажи сильные и слабые стороны кратко, без лишней детализации.
- Если ТЗ слабое или неполное, прямо скажи об этом.
- Если ТЗ хорошее, но требует доработки, тоже укажи это.
- Не пиши ничего вне структуры выходной модели.

Верни результат строго в формате, совместимом со схемой LLMOutput.
Не добавляй markdown, пояснения, вступление, заключение или любой текст вне полей структуры.
"""

system_prompt = SystemMessagePromptTemplate.from_template(sp_text_handler)

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


sp_standard_handler = f"""

"""