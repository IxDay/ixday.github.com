+++
title = "Pytest command line"
date = 2016-03-16
categories = ["Snippet"]
tags = ["tests", "python", "cli"]
url = "post/pytest_command_line"
+++

I have recently dug into pytest documentation, and moreover into
the command line arguments and I finally found a better workflow
for running tests while I develop.

Here is the command I run when I just made some devs:
`py.test -xvvvs --tb=line --pdb`

  * **-x** will stop execution on first failue, useful when debugging tests
    in order of appearance (recommended)
  * **-vvv** will display current test path (reuasable in py.test), the
    path will avoid to rerun all the previous tests before going to the one
    you are currently working on. In addition, the verbose flags will display
    a full diff of the assert error. This will help troubleshoot from where
    the error is coming
  * **--tb=line** will only display the last line of the unhandled exception
    which make the test fail. I use this because in general the exception
    trace is not really relevant, as we already know what test is running
    and the name of the error.
  * **--pdb** this one is gold, it will enter pdb on fail, and if you have
    [pdbpp](https://pypi.python.org/pypi/pdbpp/) installed, this will be
    a real debugger :p

The second command I usually ran is: `py.test -v --tb=line --cov=src --cov-report=html`.
This one is my code coverage command. I run it on small changes or after the first
command described here, in order to check that my coverage has evolved.
I put a `-v` for displaying the path of the current test, and `--tb=line` in order
to not pollute the screen with stack traces. Those two options are here in the case
of an error during the coverage. It allows me to be faster on isolating tests to
debug.

## Bonus

Here I discussed about the coverage, and py.test commands. I use specific options
on CI (travis) in coordination with [coveralls](https://coveralls.io/).

My `.travis.yml` file which run py.test with coverage.

```yaml

language: python
install:
  - pip install coveralls
script:
  - py.test -v -tb=line -cov-report= --cov=src
after_success:
  coveralls
```

It will just put the report into coveralls and displays the
same informations as the other command explained above.
