# Slack Integration

Eventilator can send events to your Slack channel!


# Configuration


``` go
type SlackConfig struct {
    Token         string
    Enabled       bool 
    Channel       string
    AuthorName    string
    AuthorSubname string
    TriggerOn   []string
}
```


Defaults:
	Token: "use-your-slack-token"
	Enabled: false
	Channel: "sentinel-events"
	AuthorName: "eventilator"
	AuthorSubname: <sentinel hostname>
	TriggerOn: ["+sdown","-sdown","+odown","-odown","+promoted-slave"]

In order to use slack integration you need to add it to the `eventilator.conf` config as follows:

```
[slack]
Enabled = true
Token = "your-token-here"
```

Subsitituing your actual slack OAuth token of course. You can also specify the list of events you want it to trigger on by doing the following:

```
[slack]
Enabled = true
Token = "your-token-here"
TriggerOn = ["+odown","-odown","+promoted-slave"]
```

In which case it will only call ot to the slack channel for objective down state changes as well as whenever a slave is promoted.

