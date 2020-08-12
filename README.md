gsuite-collector
=======

gsuite-collector: Open-Source GSuite Log Collector

[gsuite-collector](https://github.com/rfizzle/gsuite-collector) is an open-source collector designed to pull activity and audit logs for Google's GSuite service. It provides the ability to export results to a number of different destinations, such as Google Cloud Storage, Amazon S3, Stackdriver, file, and HTTP endpoint.

### Install

Installation of gsuite-collector is dead-simple - just download and extract the zip containing the [release for your system](https://github.com/rfizzle/gsuite-collector/releases/), and run the binary. gsuite-collector has binary releases for Windows, Mac, and Linux platforms.

### Building From Source
**If you are building from source, please note that gsuite-collector requires Go v1.13 or above, due to its use of Go Modules!**

To build gsuite-collector from source, simply run `go get github.com/rfizzle/gsuite-collector` and `cd` into the project source directory. Then, run `go build`. After this, you should have a binary called `gsuite-collector` in the current directory.

### Docker
You can also get gsuite-collector via the official Docker container [here](https://hub.docker.com/r/rfizzle/gsuite-collector/).
The collector was built with Kubernetes in mind.

### Documentation

Documentation can be found via the [docs](https://github.com/rfizzle/gsuite-collector/tree/master/docs). Find something missing? Let us know by filing an issue!

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let us know! Please don't hesitate to [file an issue](https://github.com/rfizzle/gsuite-collector/issues/new) and we'll get right on it.

### License
```
MIT License

Copyright (c) 2020 Coleton Pierson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```