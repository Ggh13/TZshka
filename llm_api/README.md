## Эндпоинты для взаимодействия с LLM_API

### 1. **Загрузка ТЗ в виде сплошного текста**

**Post** `llm/text`

**Request:**
```json
{
  "content": "some_text"
}
```

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
 ]
}
```