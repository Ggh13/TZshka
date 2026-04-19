## Эндпоинты для взаимодействия с LLM_API

### 1. **Загрузка ТЗ в виде сплошного текста**

**POST** `llm/text`

**Request:**
```json
{
  "content": "some_text"
}
```
- `content` - текст технического задания.

**Response (200)**

```json
{
 "status": "valid | issues_found",
 "issues": [
   {
     "rule_id": "R1",
     "problem": "net tz",
     "explanation": "resultat hz"
   }
 ],
  "feedback": "very nice bro"
}
```
- `status` - было ли найдено хотя бы одно противоречие (valid - нет, issues_found - да);
- `issues` - список найденных противоречий;
- `rule_id` - идентификатор противоречия из реестра противоречий;
- `problem` - название найденного противоречия
- `explanation` - указание на противорчие в техническом задании пользователя и его объяснение;
- `feedback` - общая оценка технического задания.

### 2. **Загрузка ТЗ в виде файла .txt**

**POST** `llm/text_file`

**Request (multipart/form-data):**

```
file: tz.txt
```
- `file` - текст технического задания в файле формата .txt.

**Response (200)**

```json
{
 "status": "valid | issues_found",
 "issues": [
   {
     "rule_id": "R1",
     "problem": "net tz",
     "explanation": "resultat hz"
   }
 ],
  "feedback": "very nice bro"
}
```

- `status` - было ли найдено хотя бы одно противоречие (valid - нет, issues_found - да);
- `issues` - список найденных противоречий;
- `rule_id` - идентификатор противоречия из реестра противоречий;
- `problem` - название найденного противоречия
- `explanation` - указание на противорчие в техническом задании пользователя и его объяснение;
- `feedback` - общая оценка технического задания.