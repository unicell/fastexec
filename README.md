# Fastexec

[![GoReportCard Widget]][GoReportCard] [![Travis Widget]][Travis] [![Coverage Status Widget]][Coverage Status]

[GoReportCard]: https://goreportcard.com/report/github.com/unicell/fastexec
[GoReportCard Widget]: https://goreportcard.com/badge/github.com/unicell/fastexec
[Travis]: https://travis-ci.org/unicell/fastexec
[Travis Widget]: https://travis-ci.org/unicell/fastexec.svg?branch=master
[Coverage Status]: https://coveralls.io/github/unicell/fastexec
[Coverage Status Widget]: https://coveralls.io/repos/github/unicell/fastexec/badge.svg

<hr>

### About

Fastexec is a tool to parallelly run shell command on data subset.

The idea is simple, fastexec takes data from standard input, divide it into
small chunks, pass it to workers running command line in parallel. And then it
will return the combined result from all workers.

### Examples

Run 1000 workers to curl URLs from the url_list.txt and output response code,
time cost for name resolution, connection and for entire request.

```bash
cat url_list.txt | fastexec -workers 10000 xargs curl -sL \
-w "%{http_code} %{time_namelookup} %{time_connect} %{time_total} %{url_effective}\\n" -o /dev/null
```

Parallelly run fping with 400 workers, and each takes 20 IPs from the
ip_list.txt. Fping will be triggered with `-t 1000` (1000 millisec timeout) and
all workers share the 1000 IP/s (1000 line of data / sec) rate limit.

```bash
cat ip_list.txt | fastexec -chunks 20 -workers 400 -ratelimit 1000 fping -t 1000
```

### Usage

```bash
  -alsologtostderr
        log to standard error as well as files
  -chunks int
        size of data chunk for one job (default 1)
  -log_backtrace_at value
        when logging hits line file:N, emit a stack trace (default :0)
  -log_dir string
        If non-empty, write log files in this directory
  -logtostderr
        log to standard error instead of files
  -ratelimit int
        data processing ratelimit / sec (default -1)
  -stderrthreshold value
        logs at or above this threshold go to stderr
  -v value
        log level for V logs
  -vmodule value
        comma-separated list of pattern=N settings for file-filtered logging
  -workers int
        num of workers (default 1)
```
