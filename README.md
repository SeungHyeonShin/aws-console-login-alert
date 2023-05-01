# About
Lambda function that sends an alarm to Slack when Webconsole login event occurs on AWS

<img width="435" alt="image" src="https://user-images.githubusercontent.com/70582181/235418642-52695fbd-4f97-471b-8a55-33de2fe5465d.png">


# Requirements
1. Slack
- Create Slack App
  - `chat:write` Permission required
  - Copy OAuth Token
- Create Channel
  - Copy `Channel id`

2. Lambda Function
- Name Handler from 'hello' to 'main'
- Set Lambda environment variable
  - `SLACK_AUTH_TOKEN`
  - `CHANNEL_ID`

3. Eventbridge
- Event Pattern
```
{
  "source": ["aws.signin"],
  "detail-type": ["AWS Console Sign In via CloudTrail"],
  "detail": {
    "eventName": ["ConsoleLogin"]
  }
}
```
- Designate as Lambda Function as Eventbridge's destination

# Build
```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
zip function.zip main
```
