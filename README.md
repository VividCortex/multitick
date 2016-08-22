multitick
=========

This package is a wrapper around a single `time.Ticker` in the Go language, which
provides a way for lots of listeners to subscribe to a single set of ticks.  If
a subscriber is busy and can't receive a tick, it will be discarded; the
implementation is non-blocking.  Additionally, ticks can be aligned to a
specific boundary easily, which eliminates annoying code from your programs.

We use this functionality at [VividCortex](https://vividcortex.com/) to help
coordinate worker goroutines, driving them by passing them ticker channels.
This makes the workers more testable, by avoiding nondeterministic use of
`time.Tick()` inside them.

Documentation
=============

Please read the generated [package
documentation](http://godoc.org/github.com/VividCortex/multitick).

Getting Started
===============

Install multitick in the usual way:

    go get github.com/VividCortex/multitick

And import it as usual:

    import (
        "github.com/VividCortex/multitick"
    )

To use multitick, create a ticker with the desired interval and offset.
For example, to create a once-per-second ticker that will send at 750
milliseconds past the clock tick,

    tick := multitick.NewTicker(time.Second, time.Millisecond*750)

If you want the ticker to start immediately instead of waiting, you can
pass a negative number as the second parameter.

Now subscribe to ticks. Here's an example of starting worker routines with
ticker channels:

    go someFunc(tick.Subscribe())
    go otherFunc(tick.Subscribe())

Contributing
============

We only accept pull requests for minor fixes or improvements. This includes:

* Small bug fixes
* Typos
* Documentation or comments

Please open issues to discuss new features. Pull requests for new features will be rejected,
so we recommend forking the repository and making changes in your fork for your use case.

[![Build Status](https://circleci.com/gh/VividCortex/multitick.png?circle-token=908b24495ba93c1070bddd2c0423f29056ef6007)](https://circleci.com/gh/VividCortex/multitick)

License
=======

Copyright (c) 2013 VividCortex, licensed under the MIT license.
Please see the LICENSE file for details.
