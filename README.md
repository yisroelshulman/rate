# **rate**

An implementation of a buffered and unbuffered limiter to control the frequency of operations over time. The limiters will grant or deny permission at the rate provided when creating them.

## Installation

```sh
go get github.com/yisroleshulman/rate
```

## Features

- concurrent usage
- non-concurrent usage

## Usage

The rate package contains 2 types of limiters:
- BufferedLimiter
- UnbufferedLimiter

While the general functionality is the same they each have some unique behaviors which will be useful depending on what the limiter is needed for.

Once the rate package was installed (verify in the go.mod file) it is ready to be imported for use.
```go
// imports
import (
    "time"
    "github.com/yisroleshulman/rate"
)


// declare and initialize the limiters
var limiter *rate.UnbufferedLimiter
var bufLimiter *rate.BufferedLimiter
limiter = NewUnbufferedLimiter(2, time.Second)
bufLimiter = NewBufferedLimiter(1, 5, time.Second)

```

The above code will create 2 limiters a *UnbufferedLimiter with a rate of 2 permissions per second, and a *BufferedLimiter with a rate of 1 permission per second and a buffer with a max capacity of 5.

Using the limiters
```go
func main() {
//... prior code

 err := operation(limiter, timeout)
 if err != nil {
    // do something
 }
 err = op(bufLimiter, timeout)
 if err != nil {
    // do something
 }

// more code ...
}

func operation(limiter *rate.UnbufferedLimiter, timeout *time.Duration) error{
    err := limiter.Wait(timeout)
    if err != nil {
        // request timed out
        return err
    }
    doOperation()
    return nil
}

func op(limiter *rate.BufferedLimiter, timeout *time.Duration) error {
    err := limiter.Wait(timeout)
    if err != nil {
        switch err.(type) {
            case *LimiterWaitTimedOut:
                // do something if needed
                return err
            case *LimiterBufferFull:
                // do something if needed
                return err
            default:
                // something is wrong
                return someErr
        }
    }
    doOperation()
    return err
}
```

It is important to note that Wait(timeout) is a blocking function and only returns when permission is granted or the request error either by timeout or the buffer is full.

Only the UnbufferedLimiter supports non-blocking request.

```go
func (limiter *UnbufferedLimiter, timeout *time.Duration) {
    timeRemaining, err := limiter.TryWait(timeout)
    if err != nil { // the rate has been reached
        // do something ex.
        return
    }
    doOperation()
}
```

The TryWait(timeout) function does not block and returns an error if the limiter has reached the limit and the time remaining until the next permission can be granted.

> Note: Although both limiters can be used for both concurrent and non-concurrent uses it is recommended to use the UnbufferedLimiter for non-concurrent use cases to reduce resource usage. and the BufferdLimiter for concurrent uses to prevent request starvation.
>
> Since wait is a blocking function calling wait concurrently on a UnbufferedLimiter may lead to request starvation as there is no way to guarantee the order in which permission is granted. And the BufferedLimiter grants permission in the order the requests were put in the buffer. The BufferedLimiter is SUBJECT TO RACE CONDITIONS due to the time delta from the approval and returning from Wait.

## Tests

The tests take < 30 seconds

## Contributing
I consider this project feature complete.
Reported issued will be looked at for bug fixes.
