[![Go Reference](https://pkg.go.dev/badge/github.com/bitfield/qrand.svg)](https://pkg.go.dev/github.com/bitfield/qrand)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitfield/qrand)](https://goreportcard.com/report/github.com/bitfield/qrand)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go)
![Tests](https://github.com/bitfield/qrand/actions/workflows/test.yml/badge.svg)

# What is `qrand`?

`qrand` is a Go package that provides random numbers derived, ultimately, from a non-deterministic, quantum-mechanical process. 

```go
import "github.com/bitfield/qrand"
```

The random data is provided by the [ANU Quantum Numbers](https://quantumnumbers.anu.edu.au/) (AQN) API. You'll need an API key for this service, but it's free for limited use (see the website for details on how to pay for more data if you need it).

# Usage

Here are a couple of example programs that show how you might use `qrand`.

## Reading random bytes

A common use of `crypto/rand` in Go programs is to read a sequence of cryptographically secure random bytes (an initialization vector, for example). `qrand` can do the same thing, but deriving its data from the quantum randomness provider.

The [`numbers`](example/numbers/main.go) example shows how to do this:

```go
q := qrand.NewReader(apiKey)
buf := make([]byte, 10)
_, err := q.Read(buf)
```

As you can see, this is very similar to the corresponding [`crypto/rand` example](https://pkg.go.dev/crypto/rand#example-Read). The only difference here is that we need to create the reader first with `NewReader`, because the provider requires an API key.

## Generating random numbers

Go's `math/rand`, on the other hand, is commonly used to provide random numbers within a desired interval, using something like `rand.Intn`. This is useful in games, for example, or other programs that need “random-seeming” behaviour, but not strict cryptographic security.

The [`password`](example/password/main.go) example shows how to do this with `qrand`, by creating a *randomness source*:

```go
rnd := rand.New(qrand.NewSource(qrand.NewReader(apiKey)))
password := make([]byte, 32)
for i := range password {
    password[i] = chars[rnd.Intn(len(chars))]
}
```

# The CLI tool

There's a simple CLI tool to request and display a given number of random (hex) bytes. To install it:

```sh
go install github.com/bitfield/qrand/cmd/qrand@latest
```

To use it, pick the number of bytes you need (for example, 32), and run:

```sh
qrand 32
```
```
8e8c2771be5c2bb10d541a5bf6aa51203e0bce2d6d4fa267afd89a6e20df11f1
```

# Sources of randomness

> *Random numbers should not be generated with a method chosen at random.*
>
> —Donald Knuth, [“The Art of Computer Programming, Vol 2: Seminumerical Algorithms”](https://amzn.to/3Y8uMt3)

Most computer random number generators (RNGs) use a deterministic process, which means that given an initial seed value, the sequence of generated numbers is predictable.

For example, Go's standard `math/rand` library uses a fairly simple algorithm to generate a random-looking, but still deterministic sequence of numbers. For most applications this is absolutely fine when seeded with a suitable value, such as the current Unix time in nanoseconds (which is the default from Go 1.20 onwards). 

For any cryptographic purposes, though, `math/rand` is insecure, and we should use `crypto/rand` instead. `crypto/rand` will use the most secure randomness source provided by the operating system; for example, on Linux systems this might be the `/dev/urandom` or `/dev/random` devices. 

While this is still technically a pseudo-random source, it uses environmental 'noise' such as I/O activity, keystrokes, and so on, to generate numbers which are in practice (though not in principle) unpredictable.

For very high-security applications, though, we can use quantum-mechanical sources, such as the cosmic microwave background radiation:

* [Lee, J. S., & Cleaver, G. B. (2017). The cosmic microwave background radiation power spectrum as a random bit generator for symmetric-and asymmetric-key cryptography. Heliyon, 3(10).](https://arxiv.org/abs/1511.02511) 

The outcomes of quantum measurements, such as the spin of an electron or the polarization of a photon, are *in principle* unpredictable, to the best of our knowledge:

* [Bierhorst, P., Knill, E., Glancy, S., Zhang, Y., Mink, A., Jordan, S., ... & Shalm, L. K. (2018). Experimentally generated randomness certified by the impossibility of superluminal signals. Nature, 556(7700), 223-226.
](https://arxiv.org/abs/1803.06219)

Hardware RNGs are available that can use such measurements to generate random data at fairly high bitrates (many GiB/s). 

# The AQN service

Australia National University provides a public quantum randomness source, derived from a hardware RNG, via its [AQN](https://quantumnumbers.anu.edu.au/) service. This has a public API, which is the source of the data obtained with `qrand`.

The random data is generated by a device that uses a laser to measure the quantum fluctuations of the vacuum:

* [Symul, T., Assad, S. M., & Lam, P. K. (2011). Real time demonstration of high bitrate quantum random number generation with coherent laser light. Applied Physics Letters, 98(23)](https://arxiv.org/abs/1107.4438)
* [Haw, J. Y., Assad, S. M., Lance, A. M., Ng, N. H. Y., Sharma, V., Lam, P. K., & Symul, T. (2015). Maximization of extractable randomness in a quantum random-number generator. Physical Review Applied, 3(5), 054004](https://arxiv.org/abs/1411.4512).

# Why do I need quantum randomness?

You don't. The standard randomness source provided by your operating system, available via `crypto/rand`, is almost certainly good enough for any application requiring strong randomness, such as cryptography (otherwise, we're all in trouble).

However, it's fun to use a source of randomness which is entirely non-deterministic (so far as we know) and provided directly by the Universe itself. 

# Security note

`qrand` is primarily for fun, but just in case you're thinking of using it in programs, here's an important caveat. 

For games and other non-cryptographic applications, as I mentioned, any reasonably random-looking data is fine. For speed, you'd normally use `math/rand` for this.

But when you're doing cryptography (for example, hashing, signing, initializing ciphers, generating keys, and so on), `math/rand` is not good enough, and that's what `crypto/rand` is for. 

Since `crypto/rand` will use as secure a source of randomness as the local computer can provide, that's about as good as we can hope to get. This will probably derive from the operating system, as mentioned earlier.

Data from `qrand` is both more and less secure than `crypto/rand`'s. More secure, because the outcomes of quantum measurements are unpredictable in principle, as we saw earlier, whereas software-based RNGs are merely unpredictable in practice.

Less secure, because it's coming over the network. Although the connection to the API uses TLS, and is thus encrypted, it's still possible that the data could be intercepted or modified by a third party (via a [man-in-the-middle attack](https://en.wikipedia.org/wiki/Man-in-the-middle_attack), for example). To put it another way, the random data obtained by `qrand` is only as secure as your TLS connection to the AQN server.

And, of course, it all rather depends on whether you trust the AQN service itself. I offer no warranty of any kind that data obtained with `qrand` is cryptographically secure. Nor could I if I wanted to, because `qrand` itself is merely a client for this third-party service.

You have been warned.
