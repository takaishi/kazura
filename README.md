# Kazura

kazura is a deployment tool for Amazon EventBridge.

- Create Rule
- Create Target (Support Lambda only now)

## eventbridge.json
[.gitignore](.gitignore)
The eventbridge.json is a definition for EventBridge Rule and Target. 

```json
{
  "rule": {
    "name": "example-rule",
    "eventBusName": "default",
    "eventPattern": {
      "source": [
        "aws.s3"
      ],
      "detail-type": [
        "Object Created"
      ],
      "detail": {
        "bucket": {
          "name": [
            "example-bucket"
          ]
        }
      }
    }
  },
  "lambdaTarget": {
    "name": "example-function"
  }
}
```



```json
{
  "rule": {
    "name": "example-rule",
    "scheduleExpression": "cron(0 12 * * ? *)",
    "eventBusName": "default"
  },
  "lambdaTarget": {
    "name": "example-function"
  }
}
```

