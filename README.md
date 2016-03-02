# Sentinel Event Relay

[![Build
Status](https://travis-ci.org/sentinel-tools/eventilator.svg?branch=master)](https://travis-ci.org/sentinel-tools/eventilator)

This utility runs under two modes: `eventilator` and `reconfigurator`. Both
modes relay or record certain Sentinel events. With the way Sentinel
works each pod which needs to be configured with the "script".`


# Eventilator

The eventilator mode is designed to be used as the sentinel
notification-script for the given pods. To add to a running sentinel pod
connect to each sentinel and execute `sentinel set <podname>
notification-script /path/to/eventilator`. Once configured any warning
event emitted by that sentinel will call eventilator.

# Reconfigurator

In reconfigurator mode the application will be called when a slave has
successfully been promoted to master. Registration in sentinel is accomplished
via `sentinel set <podname> client-reconfig-script /path/to/registrator`. In
this mode it will relay the failover event for the given pod. To see the
details on how the failovver metrics are stored see the
[redis.md](handlers/redis.md) in the `handlers` directory.

# Which One to Use

The main criteria is whether you care about things only happening on the
elected leader sentinel, or just want/need to capture all warning level events.
If, for example, you are wanting to update a database or DNS when a failover
happens you will want to use registrator. This is because only the leader
executes registrator when configured properly. With that in place you don't
have to worry about getting three events for the same failover.

If, however, you are wanting to capture all the events and have a mechanism to
dedupe certain events such as `+switch-master` or are making idempotent calls
then eventilator is much more amenable in that it handles all warning level
events. Another option for handling failvoer events is to not look at
`+switch-master` but catch the `+promoted-slave` event. You get the same
information but *only* the leader isues this event. Thus if you catch this
event you can still have eventilator handle a failvoer without worrying about
having it executed on each sentinel.

# Installation

Due to Sentinel not supporting the passing of commandline options, or passing
environment variables, to the script each mode must be it's own command path.
The simplest route is to have the eventilator executable in place and symlink
`registrator` to it. They will need to be executable by the user running Redis
sentinel. Following traditional UNIX methodology the command detects the mode
by obtaining what name it was called by.

# Configuration

With the addition of a default handler for registrator which stores failover
metrics in a Redis instance there will be config files for each mode. These are
expected to be stored in `/etc/redis/eventilator.conf` and
`/etc/redis/registrator.conf`. The command checks how it is called and loads
the appropriate config file. If one is not found it uses a default value of a
localhost Redis instance on the default port with authentication.  As new
handlers such as monitoring hooks are added into eventilator it will look for
the configuration in it's config file.

## Eventilator: Slack Handler

Slack integration has been added. For details on how to configure it see: [slack.md](handlers/slack.md)

![Redacted Screenshot](eventilator-slack-screenshot.png)


# Custom Eventilator Handlers

The simplest way to add custom handlers is to fork the repo and add custom
eventilator handlers in a separate file in `handlers/` with the build tag `//
+build custom` and use that tag to build. An example of this is visible in
`handlers/eventilator-custom.go`. 

NOTE: this process means you need to define and register all event handlers you
need to handle as the default ones will not compile with the `custom` build tag
passed to `go -build`.

These handlers will need to:

1. accept a parser.NotificationEvent
1. return an error (or nil)

To register your custom handler follow the pattern of using `HandlerMap.SetMandler(eventname,handlerfunc)` as is done in eventilator.go.


## Neat Ideas You Could Implement

Just a collection of some things that might be really cool to do with this.

Imagine you run your Redis instances on a Docker Swarm of hosts. If so you
could write a handler for `sdown` on a slave to boot a new instance and enslave
it to the master. 

Alternatively you could do the same thing but instead spin up a new cloud
server such as an AWS EC2 or Rackspace Cloud Server VMs configured w/Redis that
you then enslave to the pod.

Use the event system to reconfigure an HAProxy somewhere.

Track the frequency of a given node being subjectively down and once it hits a
given threshold migrate it to a new instance to try to improve reliability of
the node.
