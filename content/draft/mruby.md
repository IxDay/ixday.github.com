Hey HN!

After wrestling with a recurring problem as a Golang developer, I finally found a solution worth sharing.
As a user of a compiled language I need from time to time a third party scripting
language to glue things together. In the golang ecosystem the default is bash + make, however
I grew more and more dissatified with those tools: bash/Mafiles scripting is hard,
it is not easily portable, and it lacks a lot of features. I decided to explore alternatives
and documented my journey in [this blog post](https://platipy.notion.site/The-quest-for-the-optimal-scripting-language-b013b3e35a5c4c6c8d5b4a7a31cb1508).

Long story short, I found a solution using an embeddedable version of the ruby runtime.
It ended up fulfilling all my requirements, which were:

1. Cross-Platform Compatibility: Ensuring the binary works consistently across different operating systems (Windows, Linux, MacOS).
2. CLI Primitives: The binary is built with support for optparse, glob patterns, and Windows path compatibility.
3. Built-in HTTP/HTTPS and JSON/YAML Support: Ability to make HTTP/HTTPS calls and parse JSON/YAML directly in scripts.
4. Small Self-Contained Binary: All-in-one executable without extra setup, with no additional dependencies.
5. Support for Task Execution: With the included mrake library, you can organize and manage complex task workflows, similar to Makefile dependencies.

Here is a link to the project: https://github.com/IxDay/mruby

Hopefully it solves some similar challenges for some of you as well!


Hey HN!

After wrestling with a recurring problem as a Golang developer + infrastructure engineer,
I finally found a solution worth sharing.
I repackaged mruby (a lightweight version of ruby for embedded systems)
to write cross platform complex scripts and build pipelines.

For instance, I have recurring tasks involving all those cloud native fancy tools:
vault, kubernetes, aws. Automating those usually requires calling https endpoints,
handle JSON/YAML, manipulate a few data structures and maybe pushing the result to
another endpoint. Those tasks can be launched from various places:
laptops (with any kind of distribution/OS you might think of), containers, remote VMs.
You can find two scripts illustrating some of my use cases in the examples directory of the project.


I am also frustrated by the state of the golang toolchain regarding the heavy reliance
on Bash and Make. I recurring hit the following problems:
- https://stackoverflow.com/questions/1950926/create-directories-using-make-file
- https://stackoverflow.com/questions/8889035/how-to-document-a-makefile
- https://stackoverflow.com/questions/25005637/makefile-rule-depend-on-directory-content-changes/25005796

This project solved it and for the same binary size as make and its dependencies.

None of the major scripting languages were fitting our needs. I ended up listing all the requirements which are:

1. Cross-Platform Compatibility: Ensuring the binary works consistently across different operating systems (Windows, Linux, MacOS).
2. CLI Primitives: The binary is built with support for optparse, glob patterns, and Windows path compatibility.
3. Built-in HTTP/HTTPS and JSON/YAML Support: Ability to make HTTP/HTTPS calls and parse JSON/YAML directly in scripts.
4. Small Self-Contained Binary: All-in-one executable without extra setup, with no additional dependencies.
5. Support for Task Execution: With the included mrake library, you can organize and manage complex task workflows, similar to Makefile dependencies.

I documented my exploration in the following blog post.

Here is a link to the project: https://github.com/IxDay/mruby

Hopefully it solves some similar challenges for some of you as well!





We ended up writing things in bash as it was the biggest common denominator.
We could have ended up using another



------------------



# Hey HN!

After wrestling with a recurring challenge as a Golang developer and infrastructure engineer,
I finally found a solution worth sharing.
I repackaged `mruby`—a lightweight version of Ruby designed for embedded systems—to create a tool for writing cross-platform, scripts and build pipelines.

As a mainly Golang developer I always had to struggle with Make + Bash for the build
and scripting of my projects (since they are the default of the ecosystem).
This left me frustrated as those are hardly portable and complex to maintain.

`ruby` + `rake` would be a good solution but it forces to install the ruby ecosystem
which can be hard to perform on a remote/container host, or conflict with the default installation on your system (I look at you OSX)
It is not self contained, i.e: a single binary, and it is quite heavyweight in the end (~24MB)

Here enters this project, MRuby allows to bundle Ruby code in one small (2MB with my dependencies) binary.
In my release I included:
- Optparse/Ansi colors to easily parse options and arguments and easily build CLIs scripts
- YAML/JSON support
- HTTP/HTTPS client/server capabilities
- A simplified version of Rake

Depending on your needs you can add or remove dependencies and repackage the binary.
I found this setup to be sufficient to cover 99% of my needs.

I explored some solutions you can check out this blog post to see my quest.

I now mostly use this as a replacement to any encounter where bash or Make would
have make sense.

I included a few scripts example in the repository to illustrate the kind
of
I also included a Golang project skeleton to showcase how I use it in this
situation.

Hopefully, it helps solve some of the challenges you’re facing too! Feedback and contributions are welcome.



I wanted something self contained (one binary) to ease installation on the various
OS my teammate may own. And potentially lightweight as I might want to install it within a container.




## Complex scripts

Nowadays, I mostly face tasks involving cloud-native tools, accessible through an API. Automating these tasks usually requires:

- Making HTTPS calls,
- Handling JSON or YAML,
- Manipulating data structures,
- Pushing results to other endpoints.

These scripts need to run across diverse environments: laptops (on any OS or distribution), containers, and remote VMs.
A lot of interpreted languages are available to handle this but they come at a
huge space cost and often lack some dependencies which needs to be installed on
the target machine (no native YAML support for Python, horrendous option parsing in nodejs by default).

Bash ended up to be our fallback scripting language but it needed some extra cli installed and
support even between Linux and Mac was unsatisfying to say the least.

Here enters this project, MRuby allows to bundle Ruby code in one small (2MB with my dependencies) binary.
In my release I included:
- Optparse/Ansi colors to easily parse options and arguments and easily build CLIs scripts
- YAML/JSON support
- HTTP/HTTPS client/server capabilities
- *Bonus:* a simplified version of Rake

To give you more insights, I’ve included two example scripts in the project’s repository that showcase some real life use cases.

## Build pipelines

In the Go ecosystem, Bash and Make are the go-to options for scripting and build automation. However, they come with their own set of frustrations. For example:
- [Creating directories in Makefiles](https://stackoverflow.com/questions/1950926/create-directories-using-make-file) feels unnecessarily complex.
- [Documenting Makefiles](https://stackoverflow.com/questions/8889035/how-to-document-a-makefile) is a chore.
- [Handling file dependencies](https://stackoverflow.com/questions/25005637/makefile-rule-depend-on-directory-content-changes/25005796) can quickly become convoluted.

Over time, I found Bash and Make difficult to maintain, especially across different platforms.
The `mruby` binary now allows me to handle this as well, packed into a binary comparable in size to Make and its dependencies, but with far greater flexibility and functionality.
You can find a sample Golang project using MRuby/MRake in the example directory as well

## More details
If this sounds interesting, I documented my journey and exploration in a [blog post](#).

The project is open source and available here: [GitHub Repository](https://github.com/IxDay/mruby).
You can check there what's been included and extend it if needed.

Hopefully, it helps solve some of the challenges you’re facing too! Feedback and contributions are welcome.




Here’s a smoother and polished version of your draft:

```markdown
# MRuby a fully-featured scripting platform for all operating systems in a single, small binary.

Hey HN!

After wrestling with a recurring challenge as a Golang developer and infrastructure engineer, I finally found a solution worth sharing.
I repackaged `mruby`—a lightweight version of Ruby designed for embedded systems—to create a tool for writing cross-platform, complex scripts and build pipelines.

## Complex Scripts

These days, many of my tasks involve working with cloud-native tools, often through APIs. Automating these tasks typically requires:
- Making HTTPS calls,
- Handling JSON or YAML,
- Manipulating data structures,
- Pushing results to other endpoints.

These scripts need to run across diverse environments: laptops (on any OS or distribution), containers, and remote VMs.
While there are many interpreted languages available for this purpose, they come with significant downsides:
- A large disk space footprint,
- Missing dependencies that require additional setup (e.g., no native YAML support in Python),
- Poor usability out-of-the-box (e.g., Node.js has subpar default CLI option parsing).

As a fallback, I often relied on Bash scripts (at least it is lightweight and available by default).
However, even Bash needed additional CLI tools, and compatibility between Linux and macOS was frustratingly inconsistent.

This project is my attempt at solving this problem.
I’ve included [two example scripts in the repository](https://github.com/IxDay/mruby/tree/main/examples) to demonstrate how it handles real-world use cases.

## Build Pipelines

In the Go ecosystem, Bash and Make are the go-to tools for scripting and build automation. While widely used, they come with their own set of challenges.
I talked about scripting in the previous section, so I will only address Make here:

- [Creating directories in Makefiles](https://stackoverflow.com/questions/1950926/create-directories-using-make-file) feels unnecessarily complex.
- [Documenting Makefiles](https://stackoverflow.com/questions/8889035/how-to-document-a-makefile) is tedious.
- [Handling file dependencies](https://stackoverflow.com/questions/25005637/makefile-rule-depend-on-directory-content-changes/25005796) can quickly become convoluted.

Rake is (I think) a much better tool for managing these tasks, so I decided to include a simplified version in the binary as well.
I now have a single, compact binary (about 2.2MB) that handles these tasks effortlessly.

To demonstrate, I’ve also included a [sample Golang project](https://github.com/IxDay/mruby/tree/main/examples/golang_project)
using `mruby` and `mrake` in the repository.

## More Details

I’ve documented my journey and exploration looking for a solution in a [blog post](https://platipy.notion.site/The-quest-for-the-optimal-scripting-language-b013b3e35a5c4c6c8d5b4a7a31cb1508).

The binary size is just slightly over 2MB, which is comparable to make + busybox + jq (about 2MB).
And way below fully fleged scripting languages: python is 42MB, ruby 21MB, nodejs 54MB.

The project is open source and available here: [GitHub Repository](https://github.com/IxDay/mruby).
Feel free to explore what’s included, and extend it to suit your needs.

I hope this tool helps address some of the challenges you’re facing as well! Feedback and contributions are welcome.
```



Hey HN!

After struggling with a recurring challenge as a Golang developer and infrastructure engineer, I’ve finally found a solution worth sharing.
I repackaged mruby—a lightweight Ruby runtime for embedded systems—into a tool for writing cross-platform scripts and build pipelines.

As a Golang developer, I often had to rely on Make + Bash for builds and scripting.
While they're the ecosystem defaults, they're far from ideal:
Bash does no support most data structures, it badly handles strings, and so on... (a good list here: https://www.quora.com/What-are-the-disadvantages-of-the-Bash-shell)
Make has no glob feature (subdirectory traversal need to be done through `find`),
it is often used as a task runner whereas it was not designed for this, etc... (I should make a list for this one as well)
On top of this, cross platform support can be really tricky (we all got a bug due to a GNU coreutils flag at one point).

Ruby + Rake seemed like a better fit, but the Ruby ecosystem isn’t lightweight.
Installing Ruby on remote hosts or containers can be tricky.
It can conflicts with system versions (macOS’s default Ruby, anyone?).
And it is not self-contained (multiple files to be copied instead of a single binary)
in addition to be a bit heavyweight: Ruby - Rake + dependencies add up to ~24MB.

Enter this project: a repackaged mruby binary (just ~2MB with my dependencies) that bundles additional libraries in a single file.
I added the following to fit my needs,
CLI tools: Optparse for argument parsing, ANSI colors for better output.
Data handling: Built-in YAML/JSON support.
Networking: HTTP/HTTPS client/server capabilities.
Task management: A simplified version of Rake.

You can customize it too—add or remove dependencies and repackage the binary as you see fit.

I now use this as a replacement for tasks where Bash or Make would have been my go-to.
The repository includes example scripts (using `kubectl` or `vault`) and
a Golang project skeleton to demonstrate how it works.

If you’re interested in my journey exploring alternatives, check out my blog post.

Feedback and contributions are welcome—I hope it solves some of your challenges too!


------

Show HN: MRuby a cross platform scripting environment in a single, small binary


Hey HN!

If you are a Golang/Infrastructure developer like me, you might also be struggling
with Bash and Make. Well, I think I found a solution worth sharing!
I repackaged mruby—a lightweight Ruby runtime for embedded systems—into a tool for writing cross-platform scripts and build pipelines.

While Make + Bash are the ecosystem default, they’re far from ideal:
Bash lacks support for most data structures, handles strings poorly, and has many other shortcomings (a good list here: https://www.quora.com/What-are-the-disadvantages-of-the-Bash-shell).
Make doesn’t include globbing for subdirectory traversal (you need to use find for that), is often misused as a task runner, and has its own limitations. On top of this, achieving cross-platform support is tricky (we’ve all run into bugs caused by GNU vs BSD coreutils flags).

Ruby + Rake seemed like a better fit, but
- The Ruby ecosystem isn’t lightweight: Ruby + Rake + dependencies add up to ~24MB.
- Installing Ruby on remote hosts or containers can be challenging.
- It may conflict with system versions (macOS’s default Ruby, for instance).
- It’s not self-contained (you need multiple files instead of a single binary).

This project offers a different approach: a repackaged mruby binary (just ~2.2MB with my dependencies) that bundles useful libraries into a single file.
I included the following to meet my needs:

- CLI tools: Optparse for argument parsing, ANSI colors for better output.
- Data handling: Built-in YAML/JSON support.
- Networking: HTTP/HTTPS client/server capabilities.
- Task management: A simplified version of Rake.

You can customize it (add or remove dependencies and repackage the binary) to fit your specific requirements.

I now use this as a replacement for tasks where Bash or Make would have been my first choice.
The repository includes example scripts (e.g., using kubectl or vault) and a Golang project skeleton to show how it all works.

If you’re interested in my journey exploring alternatives, check out my blog post: https://platipy.notion.site/The-quest-for-the-optimal-scripting-language-b013b3e35a5c4c6c8d5b4a7a31cb1508

Feedback and contributions are welcome—I hope it helps with some of your challenges too!
