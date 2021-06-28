+++
title = "The wild kubectl logs issue"
date = 2021-01-20
categories = ["Snippet"]
tags = ["kubernetes", "admin", "cli"]
url = "post/kubectl_logs"
+++

A quick post to present one of my finding during my Kubernetes journey. It may
help people since it took me some time to find this out.
I am currently using Kubernetes a lot for my job. I am part of the infrastructure team
and need to debug some setups. I am using `kubectl logs` extensively and I found a few
interesting options I'd like to share.

Most of the time you will have multiple containers handling requests and you want
to see what is happening in all of them. I did a bit of Googling and found out
that you can pass a label selector: `kubectl logs -l app.kubernetes.io/name=foobar`

This is handy but I also wanted to see from which pod the log line is coming from.
I checked the man page and found the `--prefix` which gives me this information.

For a moment everything was perfect, I used this command for some of my testing and it worked well.

A wild issue appears!
---------------------

Later on I had to debug an application receiving a lot of traffic. I used the label
selector and was not able to `grep` my entries. I extended the `--since` duration
but the amount of logs I received was not matching what I was asking for.
I wondered myself what I did wrong? Was there some environment variables messing things up?
One of my alias is conflicting and introduce an unexpected behavior?

__Is `kubectl logs` command broken when used with label filters!?__

Actually, ... it is totally expected.
I checked the man page and read it entirely to find out an interesting option and
its default value: `--tail`. Here is what the documentation says:

```txt
--tail=-1       Lines  of recent log file to display. Defaults to -1 with no selector, showing all log lines otherwise 10, if a selector is provided.
```

So, everything was working as expected it's just that this option is different when
a selector is present. __GREAT!__

I am now using the following alias to handle logs:
`alias klog='kubectl logs --tail=-1 --prefix --timestamps'`

And my life is wonderful since then :).

