// This package implements interceptors to limit the amount of requests that the server can handle.
//
// This approach, such as circuit breaker approach, aim on keep the servers most stable as possible,
// restarting it, after a given number of requests (to free memory/cache, etc).

package ratelimit
