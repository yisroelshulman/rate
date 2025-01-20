// Package rate implements rate limiters to control the frequency of operations over time. The
// limiter will grant or deny permissions enforcing the provided rate.
//
// The package contains two types of limiters the buffered limiter and the unbuffered limiter. The
// buffered limiter has an internal buffer that ensures that concurrent requests for permissions are
// granted in the order that they were received by the limiter. The unbuffered limiter also has a
// non-blocking option.
package rate
