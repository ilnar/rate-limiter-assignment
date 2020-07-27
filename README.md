# %YOUR_FAVOURITE_MYTHOLOGICAL_HERO%

Rate Limiter - GoLang Coding Assignment 

## Assignment

The task is to produce a rate-limiting module that stops a particular requester
from making too many http requests within a particular period.

The module should expose a method that keeps track of requests and limits it
such that a requester can only make 100 requests per hour. After the limit has
been reached, return a 429 with the text "Rate limit exceeded. Try again in #{n}
seconds".

## Solution

Solution implements the simplest single-threaded precise rate limiter (see more details below).
The focus here is on extensibility and quality of the code.  

### Build and run

```shell script
go build
./rate-limiter-assignment
```

Open http://localhost:8080 in browser.

### Test

The following URLs help to test different scenarios. The rate limits are set low and
can be reached by simply refreshing the page several times.

1. First user is limited at 10 requests per minute and has two API keys.
The rate limit applies on both keys http://localhost:8080/foo?api_key=user1key1
and http://localhost:8080/foo?api_key=user1key2 meaning that reaching limit on one key
will block requests on the other key.

2. Second user is limited at 10 requests per minute and has a single API key:
http://localhost:8080/foo?api_key=user2key1

3. Third user is limited at 1 request per second and has a single API key:
http://localhost:8080/foo?api_key=user3key1

4. Fourth user is unknown to the policy provider and this link returns error intentionally:
http://localhost:8080/foo?api_key=user4key1

## Requirements

The assignment should be completed in ~5 hours

Main focus on:

- support for arbitrary rate limiting / quota enforcement strategies
  
- implementation of the default strategy: 100 requests per hour

- returned errors should be informative and actionable

- code quality

- if the module includes integration with HTTP(S) server,
the latter will return machine-readable HTTP status code 429 together 
with error message when the rate limit is exceeded.

The task does not specify the characteristics of workloads (e.g. throughput)
and the architecture of system that will use this module.

## Assumptions

The rate limit is set per user. Users are authenticated by API keys.

## Scope

Given time constraints and abstract problem statement the simplest rate limiter is implemented
without any optimisations.

Extensibility is demonstrated through integration with other components.
These components are not implemented and fakes with hardcoded values are used instead.

The rate limiter does not perform garbage collection for users that stopped sending requests
or changed rate limiter.

## Design

### Monolithic service

High precision, high performance, low scalability

For a strategy that allows N requests over time interval T:

Keep a list of timestamps of requests, as soon as a new request comes in,
remove all the entries older than timestamp of the new request minus T,
calculate the size of the list and if it is less than N, add new timestamp
and allow request. Otherwise, discard request.

The list of timestamps will be sorted on single-threaded version of the
limiter, hence removing the old entries will be easy, using binary search.

The list itself is not a classical linked list, but a data structure
that support fast access by index and fast truncation (e.g. go slice).

This approach stores timestamps of requests and performs a series of reads
and writes in corresponding data structure, meaning it won't work well when
T and N are big and/or the data structure is accessed from several threads.

Nevertheless, this approach provides the highest precision and the simplest one.

## Alternatives
 
### Multithreaded monolithic service

Medium precision, high performance, medium scalability

Same as the previous one, but the list of timestamps is accessed from several threads.
To minimise synchronisation overhead, batching can be implemented, hence the precision will be
worse, meaning users might be able to go over rate limit or start to be blocked before
reaching it. Probably, can run on a very big machine.

### Distributed service with many users with moderate loads

Medium precision, high performance, medium scalability

Using distributed key-value storage shared between instances. Use API key or another bit of info
allowing to evenly distribute load between keys. At every key have a simple counter.
Every request decreases it. When the value is zero, rate limit is reached. 
Run regular counter refresh in the background encreasing counter. Use batching. 
 
This approach is likely to result in high read/writes per user.  

### Distributed service with few users with high loads

Low precision, low performance, medium scalability

Similar to previous approach, but the counter is sharded now and getting the current value of the
counter requires several calls and summation. This way, single customer won't be able to overload
single key.

This approach is likely to result in high reads, low writes per user. 

### Large scale service with a mix of users and patterns

Enforce rate limit per-region - not purely global service.

Or use eventually consistent storage for counter - low precision, because of even bigger batching. 
