// time.After() does NOT "become" a channel. It actually RETURNS a channel.

Excellent question. This confuses almost everyone when they first learn Go.

The trick is this:

> **`time.After()` does NOT "become" a channel. It actually RETURNS a channel.**

Let's look at its function signature.

```go
func After(d Duration) <-chan Time
```

Read it from right to left:

* `func After` → function named `After`
* `(d Duration)` → takes a duration
* `<-chan Time` → returns a **receive-only channel** of `time.Time`

So when you write

```go
time.After(3 * time.Second)
```

you're actually getting back a channel.

It's exactly like this:

```go
ch := time.After(3 * time.Second)
```

Now `ch` is a channel.

```go
fmt.Printf("%T\n", ch)
```

Output:

```text
<-chan time.Time
```

---

## What does `time.After()` do internally?

Conceptually, imagine Go implemented it like this:

```go
func MyAfter(d time.Duration) <-chan time.Time {

    ch := make(chan time.Time, 1)

    go func() {
        time.Sleep(d)
        ch <- time.Now()
    }()

    return ch
}
```

This isn't the real implementation, but it's very close to how it behaves.

Let's go through it.

### Step 1

```go
ch := make(chan time.Time, 1)
```

Create a channel.

---

### Step 2

```go
go func() {
```

Start another goroutine.

---

### Step 3

```go
time.Sleep(d)
```

Wait for the duration.

---

### Step 4

```go
ch <- time.Now()
```

After waiting, send the current time into the channel.

---

### Step 5

```go
return ch
```

Return the channel immediately.

So now your code has a channel that will receive one value in the future.

---

# Why do we write `<-time.After(...)`?

Suppose you do

```go
ch := time.After(3 * time.Second)
```

Remember,

```text
ch
|
v
(channel)
```

After 3 seconds,

```text
channel
   |
   |----> current time
```

To receive that value:

```go
<-ch
```

or directly

```go
<-time.After(3 * time.Second)
```

Since `time.After(...)` already returns a channel, you can immediately receive from it.

---

## Why is it used in `select`?

Look at this:

```go
select {

case result := <-apiChannel:
    fmt.Println(result)

case <-time.After(3 * time.Second):
    fmt.Println("Timeout")
}
```

Let's imagine the API never responds.

Initially:

```text
apiChannel
    empty

time.After channel
    empty
```

After 3 seconds:

```text
apiChannel
    empty

time.After channel
    12:30:15
```

Now this case becomes ready:

```go
case <-time.After(3 * time.Second):
```

The value (`time.Time`) is received and discarded because you don't need it. You only care that **the timer expired**.

---

## If you wanted the actual time

You could write:

```go
select {

case t := <-time.After(3 * time.Second):
    fmt.Println("Timer fired at:", t)

}
```

Output:

```text
Timer fired at: 2026-06-17 14:30:12 +0530 IST
```

---

## Mental model

Think of `time.After` as creating a "future delivery" channel.

```text
            time.After(3s)

                  |
                  |
          returns a channel
                  |
                  v

        +-------------------+
        |   time.Time chan  |
        +-------------------+
                  |
      (wait 3 seconds)
                  |
                  v
         sends current time
                  |
                  v
             <- receives
```

## The key takeaway

When you see:

```go
case <-time.After(3 * time.Second):
```

don't read it as:

> "`time.After` becomes a channel."

Instead read it as:

> "`time.After(3 * time.Second)` **returns a receive-only channel**. The `select` waits until that channel receives a value, which happens after 3 seconds."

This pattern is extremely common in Go for implementing **timeouts** without having to manually create goroutines and timer channels yourself.
