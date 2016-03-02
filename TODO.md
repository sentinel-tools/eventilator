# Things I Want to Add


# Alerting/Monitoring Systems
I'd like to have stock handlers which send alerts to various alerting systems
such as  Nagios, Sensu, ZenDesk, PagerDuty, New Relic, etc.. Each would require a specific
section in a `eventilator.conf` file and not take those actions if missing.

For each alerting destination I am thinking of a common/required set of
parameters such as:

	enabled=true
	retry-on-fail=true

# Notification Metrics

Storing the Failover metrics is useful but it would also be useful to record
warnings such as `sdown` events and no-good slave events. 

# Prometheus Integration

Prometheus is a neat and growing metrics storage system. An option to store
metrics in there would likely be welcome as well.

# Make Default Eventilator use Syslog

The default currently just logs to stderr. This isn't really useful but logging
to a syslog server could be a decent default action.

# Standard Log to Syslog

While stderr is useful, all the output from Sentinel executed scripts is
slurped in and output by sentinel. There should be an option to log it to
syslog; just in case you don't log Setinel to syslog but do want these events
logged there.


# More Environment Variables

While Redis itself can't set environment variables it is possible to set them
in the shell that launches Sentinel and they will be propogated to
reconfigurator/eventilator. For Sentinels launched in Docker where you can set these environemnt variables per-container this could be quite useful.

# Consul Config Storage

I am quite fond of using a configuration backing store such as Consul instead
of a config file. I'd like to add that capability here as well.

# Consul Events

For "not massive" setups, ie. where the event traffic is small, it could be
useful to relay the sentinel events to a Consul server as user events. I'd like
to add that in.

# Multiple Handlers

I think I'd like to be able to call multiple handlers for a given event. It
seems somewhat straightforward in that I could make the HandlerMapper map to a
list of EventHandlers, but then it becomes a question of ordering and error
handling. I am somewhat of the opinion that stuff like that should actually be
handled by something that already implements that mechanism such as Consul or
Sensu though. So it may not show up soon, if at all.
