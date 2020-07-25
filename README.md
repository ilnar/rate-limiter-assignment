# %MYTHOLOGICAL_NAME% - Rate Limiter

GoLang Coding Assignment

## Assignment

The task is to produce a rate-limiting module that stops a particular requester
from making too many http requests within a particular period.

The module should expose a method that keeps track of requests and limits it
such that a requester can only make 100 requests per hour. After the limit has
been reached, return a 429 with the text "Rate limit exceeded. Try again in #{n}
seconds".

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

The rate limit is set per user. Corresponding ID is provided.

## Scope

Given time constraints and abstract problem statement the module will implement
basic components and one or two rate limiting strategies.

Integrations with web server is likely to be put out of scope.

For distributed implementation, fakes will be used.

User is identified via API key or token, the authentication system is out of scope. 

## Design Options

### Monolithic service

High precision, high performance, low scalability

For a strategy that allows N requests over time interval T:

Keep a list of timestamps of requests, as soon as a new request comes in,
remove all the old entries older than timestamp of the new request minus T,
calculate the size of the array and if it is less than N, add new timestamp
and allow request. Otherwise, discard request.

The list of timestamps will be sorted on single-threaded version of the
limiter, hence removing the old entries will be easy using binary search.

This approach stores timestamps of requests and performs a series of reads
and writes in corresponding data structure, meaning it won't work well when
T and N are big and/or the data structure will be accessed from many threads.

Nevertheless, this approach provides the highest precision.

### Multithreaded monolithic service

Medium precision, high performance, medium scalability

### Distributed service with many users with moderate load

Medium precision, high performance, medium scalability

Using hashtable (key-value storage) shared between instances

Regular quota refresh in the background -> high read/write per user 

### Distributed service with few users with high load

Low precision, low performance, medium scalability

Using distributed counter -> high read, low write per user 

### Large scale service with a mix of users and patterns

Enforce rate limit per-region - not purely global service

Or use eventually consistent storage for counter - low precision 

## Future Work

TBD

## Scratch

Kernel level rate limiting for networking.
