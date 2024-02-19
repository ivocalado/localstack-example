# localstack-example

```
curl -X POST -H "Content-Type: application/json" -d '{
  "id":"12345",
  "payload":"2022-01-01,John Doe,50000\n2022-01-02,John Doe,51000\n2022-01-03,Jane Doe,60000\n2022-01-04,Jane Doe,61000\n2022-01-05,Alice Smith,55000\n2022-01-06,Alice Smith,56000\n2022-01-07,Bob Johnson,65000\n2022-01-08,Bob Johnson,66000\n2022-01-09,Charlie Brown,70000\n2022-01-10,Charlie Brown,71000"
}' http://localhost:8080/v1/submit
```