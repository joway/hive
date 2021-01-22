# Hive
![GitHub release](https://img.shields.io/github/tag/joway/hive.svg?label=release)
[![Go Report Card](https://goreportcard.com/badge/github.com/joway/hive)](https://goreportcard.com/report/github.com/joway/hive)
[![codecov](https://codecov.io/gh/joway/hive/branch/master/graph/badge.svg?token=Y1YO11FZKU)](https://codecov.io/gh/joway/hive)
[![CircleCI](https://circleci.com/gh/joway/hive.svg?style=shield)](https://circleci.com/gh/joway/hive)

Hive is a high-efficiency **Goroutine Pool** based on [pond](https://github.com/joway/pond).

## Get Started

```go
h, err := hive.New(
    hive.WithSize(10),
    hive.WithNonblocking(false),
)
if err != nil {
    log.Fatal(err)
}
err = h.Submit(context.Background(), func () {
    //do something
    time.Sleep(time.Millisecond * 100)
})
if err != nil {
    log.Fatal(err)
}
```

## Configuration

| Option            | Default        | Description  |
| ------------------|:--------------:| :------------|
| Size              | 10             |Size of the goroutine pool.|
| MinIdle           | 0              |The minimum size of the idle goroutines.|
| MaxIdle           | 10             |The maximal size of the idle goroutines.|
| MinIdleTime       | 5m             |The minimum time that an idle goroutine should be reserved.|
| Nonblocking       | false          |If true, return ErrOverload when all workers is busy. Otherwise, block Submit method.|
