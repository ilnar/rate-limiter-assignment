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

### Multithreaded monolithic service

High precision, high performance, better scalability, complexity

### Distributed service under moderate load

Low precision, high performance, high scalability

Using Memcache shared between instances

High precision, low performance, high scalability

Using distributed counter  

### Large scale service with high load

Enforce rate limit per-region - not purely global service

Or use eventually consistent storage for counter - low precision 

## Alternatives

Regular quota refresh in the background -> high read/write  

## Future Work

TBD

## Scratch

Kernel level rate limiting for networking.
