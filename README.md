# Sentinel Event Relay

This utility runs under two modes: `eventilator` and `reconfigurator`. Both
modes relay or record certain Sentinel events. With the way Sentinel
works each pod which needs to be configured with the "script".`


# Eventilator

The ventilator mode is designed to be used as the sentinel
notification-script for the given pods. To add to a running sentinel pod
connect to each sentinel and execute `sentinel set <podname>
notification-script /path/to/eventilator`. Once configured any warning
event emitted by that sentinel will call eventilator.

# Reconfigurator

In reconfigurator mode the application will be called when a slave has
successfully been promoted to master. Registration in sentinel is accomplished
via `sentinel set <podname> client-reconfig-script /path/to/registrator`. In
this mode it will relay the failover event for the given pod.

# Installation

Due to Sentinel not supporting the passing of commandline options, or passing
environment variables, to the script each mode must be it's own command path.
The simplest route is to have the eventilator executable in place and symlink
`registrator` to it. They will need to be executable by the user running Redis
sentinel. Following traditional UNIX methodology the command detects the mode
by obtaining what name it was called by.


# Custom Eventilator Handlers

The simplest way to add custom handlers is to fork the repo and add custom
eventilator handlers in a separate file in `handlers/` with the build tag `//
+build custom` and use that tag to build. An example of this is visible in
`handlers/eventilator-custom.go`. NOTE: this process means you need to define
and register all event handlers you need to handle as the default ones will not
compile with the `custom` build tag passed to `go -build`.

These handlers will need to:
1) accept a parser.NotificationEvent
2) return an error (or nil)

To register you rcustom handler follow the pattern of using `HandlerMap.SetMandler(eventname,handlerfunc)` as is done in eventilator.go.
