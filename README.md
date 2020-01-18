# rdtp - Reliable Datagram Transfer Protocol

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/rdtp)](https://goreportcard.com/report/github.com/adrianosela/rdtp)
[![Documentation](https://godoc.org/github.com/adrianosela/rdtp?status.svg)](https://godoc.org/github.com/adrianosela/rdtp)
[![license](https://img.shields.io/github/license/adrianosela/rdtp.svg)](https://github.com/adrianosela/rdtp/blob/master/LICENSE)

Specification and implementation of a reliable transport layer protocol to be used over IP networks.

Goal: Eventually be able to perform HTTP communication over this homemade transport protocol

## Header Format

```
 0      7 8     15 16    23 24    31
+--------+--------+--------+--------+
|     Src. Port   |    Dst. Port    |
+--------+--------+--------+--------+
|      Length     |    Checksum     |
+--------+--------+--------+--------+
|             ( Data )              |
+               ....                +
```
