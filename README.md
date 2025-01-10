# **request_rate_limiter**

Limit the number of times a a resource can be used within a given timeframe.

An example of a use case for a Limiter is when making requests to a server that has a limit of say 100 requests a second. The Limiter can be used to make sure only upto 100 requests are made per second.

## Installation
Download a copy of the limiter.go file and use in your program.

> Note: In the future I might extract the limiter to its own repo where it will be obtainable with a go get ...

## Usage
Once the limiter.go is imported to you program you can make use of the Limiter.

Suppose we have a struct for which we want to impose a limit for its use of a resourse we can do it with a limiter.

Create a stuct with a pointer to a Limiter and declare the struct and new *Limiter using the NewLimiter function which takes two parameters an int for the limit and a time.Duration for the rate interval.
```go
type ItemWithLimiter struct {
    Limiter *Limiter
    // other variables...
}

func main() {
    limetedItem := &ItemWithLimiter {
        Limiter := NewLimiter(100, time.Second)
        // other variables...
    }
}
```
The above coded will create a limtedItem with a limiter that handles 100 uses per second.

To make use of the limiter we can create receiver on the limeted item and call the Limiter.Wait function which takes a pointer to a time.Duration variable, like so.

```go
func (i *ItemWithLimiter) UseResource() error {
    err := i.Limiter.Wait(nil)
    if err != nil {
        // request timed out
        return err
    }
    // access the resource that is time limited
}
```
The Wait function takes a pointer to a time.Duration parameter which is used when we want the request to use the resource to timeout, by sending nil the request won't time out. Wait blocks until it received permission to return or errors when it times out. In the above example the Limiter.Wait won't time out.

> Note: in its current form the limiter DOES NOT GURANTEE any order to which requests will receive permission to access the resource. This can lead to STARVATION and/or race conditions with unintended consiquences.
>
> In a future version I plan to implement a way where the usage is allocated based on the order the limiter received the request. This might still have race conditions when multiple threads can request access so this should not be used in cases where ther thread rely on eachother for information from the limited resource.

## Contributing
Not accepting ATM
