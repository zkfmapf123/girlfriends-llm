# go-llm 

## API Examples

```sh

## vacation register
curl -X POST -H"Content-type: application/json" \
    -d'{"favorite_season": "winter", "hobbies": ["computer coding","exercise"], "budget":1000}' \
    http://localhost:8080/vacation/create

## introduce vacation plan
curl -X GET -H"Content-type: application/json" \
    http://localhost:8080/vacation/...



```

## Reference

- <a href="https://platform.openai.com/docs/overview"> OpenAPI </a>