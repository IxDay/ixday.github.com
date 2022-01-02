---
title:      "Bash base64 padding"
date:       2021-12-30
categories: ["Snippet"]
tags:       ["bash"]
url:        "post/bash_base64_padding"
---

__Disclaimer:__ I am not 100% percent sure of my formula. It works every time
I use it but I may have misunderstood some concepts. Take this article with a grain
of salt.

It happens that sometimes I have to process some base64 encoded strings which
are not padded correctly ([here][wiki_padding] is the wikipedia explanation of base64 padding).
When using the `base64 -d` command I end up getting
the following error: `base64: invalid input` and a non zero exit code.

__Note:__ This only happen with the GNU version of the `base64` utilitary. On
BSD it seems to support invalid padding.

Since I am relying on the exit code I can't ignore the error and need to fix the
padding before processing the string. Also, I would like this fix to be rather
short to insert it in a pipe flow without creating a monstruous line.
After a few experimentations I came up with a satisfying solution relying on `awk`.

In this example I will show a real life use case I encountered.
I had a [JWT][wiki_jwt] token stored in a file on my machine. This token come
with an expiration time and I would like to check if the token has expired before
requesting a new one. This will avoid going over the network if my token is still
valid.

```sh
	cat my_jwt_token.txt \
		| awk -F'.' '{l=length($2)+2; print substr($2"==",1,l-l%4)}' \
		| base64 -d | jq -r '.exp' \
		| xargs test "$(date +%s)" -lt
```

- The padding happens on the second line here is a quick explanation:

```txt
       ┌──── We are decoding a JWT token, base64 encoded fields are separated by a dot
       │                                    ┌─── We add two padding characters to the second field
       │                                    │     (the one we want to decode), which is the maximum possible padding
       │          ┌─ We calculate the new length with the additional padding
       │          │                         │
awk -F'.' '{l=length($2)+2; print substr($2"==",1,l-l%4)}'
                                                     │
  This is where we compute how much padding we need. ┘
  We calculate the modulo 4 of the string with extra padding and use it to
  truncate the length and remove the non required characters.
  The resulting string will contains the required padding characters to be base64 valid
```

- The third line decode the string and we use [jq][jq_website] to retrieve the proper field.
- Fourth line is a test against current time to determine if we are past the expiration date.

Even if I am not sure of my formula I ran a few tests using the [wiki article][wiki_padding]
to validate it works with most of the cases:

```sh
decode() {
    awk -vstr="${1}" 'BEGIN {l=length(str)+2; print substr(str"==",1,l-l%4)}' \
        | base64 -d \
        | { cat; echo; }
}

decode "bGlnaHQgd29yay4"
decode "bGlnaHQgd29yaw"
decode "bGlnaHQgd29y"
decode "bGlnaHQgd28"
decode "bGlnaHQgdw"
```

I hope this could help someone out there, I am keeping it here for my future self
to not forget about all those tricks.


[wiki_padding]: https://en.wikipedia.org/wiki/Base64#Output_padding
[wiki_jwt]: https://en.wikipedia.org/wiki/JSON_Web_Token
[jq_website]: https://stedolan.github.io/jq/
