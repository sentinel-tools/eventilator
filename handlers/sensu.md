# SensuJIT Integration

Eventilator can send events to your SensuJIT channel. Its default configuration
is also the recommended one, with the exception on enabling it course. What it
will do is send a JSON string representing either an `+odown`, `-odown`, or an
`-failover-abort-no-good-slave` event.

In the `+odown` event it will set the status to 1 to indicate this is a
warning. If an `-down` comes in for the same pod it will send the same message
with a status=1 to clear it from the Sensu display. If, however, the failover
fails do to an `-failover-abort-no-good-slave` event it will instead send an
alert indicating the pod was unable to failover de to lacking a good slave and
this will be status=2 representing a critical condition. There is no auto-clear
code for this alert.


# Configuration

``` go
type SensuJIT struct {
    Enabled     bool 
    IP       	string
	Port 		int
    TriggerOn   []string
}
```


Defaults:
	Enabled: false
	TriggerOn: [+odown","-odown","-failover-abort-no-good-slave"]

In order to use this integration you need to add it to the `eventilator.conf` config as follows:

```
[SensuJIT]
Enabled = true
Token = "your-token-here"
IP = "127.0.0.1"
Port = 3030
```

Subsitituing your actual bound IP and port for the sensu-client and port.


#Sample Messages

## `+odown`
```json
{
    "host": "the.sentinel.hostname",
    "name": "Redis Master for fa77dac3-4ee0-4040-a2cd-80aeb1be4120 is Down",
    "occurrences": 1,
    "output": "Pod web_cache_1 has a non-responsive Master",
    "source": "the.sentinel.hostname",
    "status": 1,
    "zd_description": "Pod web_cache_1 has a non-responsive Master",
    "zd_set_tags": "monitoring redis_master_odown",
    "zd_subject": "Monitoring alert for web_cache_1 Master Down"
}
```


## `-down`
```json
{
    "host": "the.sentinel.hostname",
    "name": "Redis Master for web_cache_1 is Down",
    "occurrences": 1,
    "output": "Pod web_cache_1 has a non-responsive Master",
    "source": "the.sentinel.hostname",
    "status": 0,
    "zd_description": "Pod web_cache_1 has a non-responsive Master",
    "zd_set_tags": "monitoring redis_master_odown",
    "zd_subject": "Monitoring alert for web_cache_1 Master Down"
}
```

## `-failover-abort-no-good-slave`
```json
{
    "host": "the.sentinel.hostname",
    "name": "Failover Aborted: web_cache_1 Has No Promotable Slave!",
    "occurrences": 1,
    "output": "Pod web_cache_1 has a non-responsive Master and no promotable slave",
    "source": "the.sentinel.hostname",
    "status": 2,
    "zd_description": "Pod web_cache_1 Is in Aborted Failover State",
    "zd_set_tags": "monitoring redis_pod_down",
    "zd_subject": "Monitoring alert for web_cache_1 Failover Failure"
}
```
