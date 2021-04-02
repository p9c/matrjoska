# ![https://raw.githubusercontent.com/l0k18/sporeOS/75d621579257aae2849a565ec633392970ded597/pkg/logo/logo-icon.svg](https://raw.githubusercontent.com/l0k18/sporeOS/75d621579257aae2849a565ec633392970ded597/pkg/logo/logo-icon.svg) sporə

###### Operating System as a Service

Building on the principles of distributed systems, and the legacy of Unix and
Plan 9 operating systems, spore aims to become a testbed for fracturing the
concept of operating system completely.

[![github](https://img.shields.io/badge/github-page-yellow.svg)](https://l0k18.github.io/sporeOS)
[![GoDoc](https://img.shields.io/badge/godoc-documentation-blue.svg)](https://godoc.org/github.com/l0k18/sporeOS)
[![GoDoc](https://img.shields.io/badge/chat-telegram-white.svg)](https://t.me/sporeOS)

# [Specification](https://github.com/l0k18/sporeOS/wiki/specification)

# About

###### Von Neumann Doesn't Live Here Anymore

The idea of modeling computer systems as serial, synchronised, centralised
systems is becoming very very irrelevant in modern times. More and more
important now is the connections between nodes that can execute programs, and
dealing with data from real world sources, which tend to be random and
unpredictable, both in timing and frequency.

So, operating systems designs, conventional kernel architectures, are based on
assumptions devised in Von Neumann's model of serial processing engines, but
with the speed of current processors, there is now the need for systems on the
chip to keep message streams synchronised at the micro and nanosecond level.

If distributed systems principles are starting to apply at this level, then the
whole system needs to move on to thinking about computing as ad-hoc clusters as
this is invariably what they are. By this, code can be written and used in more
freeform ways with regard to the paths and channels and sites of processing,
dynamically, making better use of all available resources connected to a user's
hardware.

## Spore shell

The starting point for any operating system is the user interface, the shell,
and the ability to gain access to binary executables that implement applications
the user wishes to run.

This is the centre of the system in the development phase, and is integrated
into a processing schedule server when the binary interfacing is done. What the
shell does is provide arbitrary duplex channels with protocol translators
potentially making RPC and IPC transparent to the applications.

To start, we need to abstract away the platform differences. In this aspect, we
implement alternative ways of using the standard I/O, as a series of pipes, with
the core feature of logging in one direction and shutdown request in the other.
This abstraction is transparently implemented for the platform, and the details
are not necessary to use it.

Next, the shell does not execute binaries, it loads the source, generates the
binaries, and as they use the low level IPC, they can be started and stopped at
will.

Further, the controlling process can broker named pipes between processes, and
eventually including ones with network transport phases in their connection.

Final first phase of development is to implement a pipe transport for conveying
Gio events to child processes and return ops to render the display accordingly,
including embedding from child processes inside child interfaces.
