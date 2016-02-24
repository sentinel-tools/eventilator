# Failover Metrics Storage in Redis

The registrator will store various metrics on failovers into a Redis store.
This document will describe the keys used, what each metric represents, and
provide some examples on extracting information from the data stored.

# Configuration

The reconfigurator mode will use the `/etc/redis/reconfigurator.conf` TOML file
to learn where its Redis is and what, if any, authentication is required. The file needs to have the following lines:
```
RedisAddress=<address>
RedisPort=<port>
RedisAuth=<requirepass>
```

# Keys Used

## Global Pod List
First, the global pod list. Every pod which is added to the store, meaning it
has had a failover event, is stored in the key `pods-with-failovers`.
This is a global Redis Set which is stored the name of every pod which
has had a failover. There is no automated cleanup.

## Failover stats

There are several keys used to aggregate and represent failovers. There
are the aggregated ones which are global, and there are the pod-specific
ones.

### Pod-Specific Keys

Each pod with a failover will have the key `failovers:<podname>:log`.
This is a Sorted Set containing the newly promoted master IP and the
UNIX timestamp in seconds since the epoch as the score. From the Redis
CLI tool you can thus see all failovers, what IP was promoted, and when
it happened for a pod by issuing the command `zrange
failovers:<podname>:log 0 -1 withscores`

If you just want a log of when the pod has had failovers you will want
to work against the key `failovers:<podname>:timestamps`. This is a
Redis Set containing the UNIX timestamp of when the failover occured.

A final pod-specific key is one of counters. The key is
`failovers:success:<podname>:counters` and is a Redis Hash of integers.
Each member of the hash represents a unit of time. Currently these are:

	YYYY
	YYYY:M
	YYYY:M:D
	YYYY:M:D:H
	YYYY:M:D:H:M

With this layout you can easily pull and show failovers for the pod as
time-series data. Hopefully there are very few of these events but if
you have a pod with strange and frequent failovers this can be useful to
see if there is a temporal aspect to the occurrences.

Thus if you wanted to see how many times the pod 'pod`'had a failover on 01
March 2016 you could issue the following in the redis-cli: `hget
failovers:success:pod1:counters 2016:3:1`. If there were none you get
the standard redis `nil` response.


### Global Aggregates

For a larger view of your sentinel and redis setup you want to turn to
aggregated metrics. These are not pod-specific but represent varous
aggregations of failover events.

The key `failovers-aggregated` is the largest aggregation on failover metrics.
It is a Redis Hash where the members are the time units described above for
pod-specific keys. By looking at this data you can see how many total failovers
happened across your sentinel constellation by each time unit. For example to
see how many failovers happened in the constellation on the 24th of February in
2016: `hget failovers:aggregated 2016:2:24`. To see how many failovers happened
in all of February of 2016: `hget failovers:aggregated 2016:2`.

To see *which* pods failed over in a given time unit consult the
`failovers:<timeunit>` keys. These are Redis Sets contining every pod that
failed over in that window. For example to see what pods failed over on 24
February 2016: `smembers failovers:2016:2:24`. To narrow that down to the ones
which happened during the noon hour: `smembers failovers:2016:2:24:12`. 

With this being a Set you can leverage Redis' set operations such as
intersections, differenced, etc. to slice and dice this data. 

In addition to the `pods-with-failovers` index, there is also a time-window
index of the form `pods-with-failovers:<window>`. Of particular note here is
that the YYYY:M:D window has a rolling 60 day expire set. Thus data for
`pods-with-failovers:2016:1:1` will expire 60 days later. Strictly speaking
this data is the same as you can obtain with the previous Set description and
may be removed in the future. However, I suspect it is useful to more than just
me to have a set of data that can be auto-expired as well as data kept
long-term for record requirements. Indeed it may become a configurable item to
have each of the windows for these sets have custom expiration times.

The keys `failovers:aggregated-by-time:<WINDOW>` stores sorted sets of
failovers by the given window. Each entry will be a sorted set where the score
is the timestamp in seconds and the value is name of the pod which was flipped.
These keys are useful for seeing what pods failed over in a given time range.
As Sorted Sets you have the power of the Sorted Set operations such as zrange,
zrevrange, and the intersections, stores, and diff operations.

Where this group really comes into its own is when you need to do things such
as "ten most recent failovers". For example to see the last 10 failovers in
2016: `zrevrange  failovers:aggregated-by-time:2016:2 0 10 `. If you want to
see the timestamps for those failovers add ` withscores` to the command. This
can also be quite useful if building a web table and you want to paginate the
table.


# Storage Requirements

The data structures were selected with an eye toward keeping the memory
requirements small.  In testing running 1000 unique pods and simulating
thousands of failovers per day results in just a few MB of memory per day. Of
course that is a rather extreme case as I would expect that if you have a
thousand pods failing over dozens or hundreds of times per day the amount of
memory this Redis store would require is the least of your concerns. ;)

The key count is related to how many pods you have and how often they fail
over. The increase is sublinear in that the higher window aggregates require
less memory due to being counters or sets which don't expand. In the
aformentioned testing it came out to around a couple thousand per day when
running near-constant failovers over a thousand pods.

In normal real world conditions it should require less than 10MB of memory per
year. Of course, the longer your pod names, the more memory it taks to store
those names. ;)
