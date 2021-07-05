---
title:      "Terraform retrieve sensible data"
date:       2021-07-03
categories: ["Snippet"]
tags:       ["cli", "terraform"]
url:        "post/terraform_sensitive"
---

Another post on Terraform. Here to share a little snippet to retrieve some data from
your Terraform state.
When you are dealing with some Terraform resources or writing modules you may
encounter the `sensitive` keyword/value ([here][sensitive_doc] is a bit of doc).
This is handy to avoid leaking some data, but from time to time you may want to
extract one of those.

## Code

For this article let's imagine I am creating an `aws_iam_user`
for another team to access some specific resources. I need to share the
`AWS_SECRET_ACCESS_KEY` and `AWS_ACCESS_KEY_ID` values but it is actually classified
as sensistive.
Here is the terraform code I am using:

```terraform
resource "aws_iam_user" "foo" {
  name = "foo"
}

resource "aws_iam_access_key" "foo" {
  user = aws_iam_user.foo.name
}
```

And here is the command to retrieve the info but it will be obfuscated:

```sh
terraform state show aws_iam_access_key.foo

...
# aws_iam_access_key.foo:
resource "aws_iam_access_key" "foo" {
    id                   = "AKIAEXAMPLE"
    secret               = (sensitive value)
    ses_smtp_password    = (sensitive value)
    ses_smtp_password_v4 = (sensitive value)
    status               = "Active"
    user                 = "foo"
}
```

## Snippet

The trick is to use the `-json` flag when running `terraform show`.
[Here][show_doc] is the doc explaining that sensitive data can be displayed when
passing the proper flag (see the blue note block).

The last thing to do is to use the `jq` command to extract what we are looking for:

```sh
terraform show -json | \
	jq '.values.root_module.resources[] | select(.address == "aws_iam_access_key.foo") | .values'
```

[sensitive_doc]: https://learn.hashicorp.com/tutorials/terraform/sensitive-variables
[show_doc]: https://www.terraform.io/docs/cli/commands/show.html
