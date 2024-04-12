---
date: 2023-12-21
filename: yourbase.html
type: draft
---

# The Sunday Pivot

One sleepy Sunday morning in January 2021,
I decided to write a Python library.
This was no ordinary Python library,
since its purpose was the same as that of the entire company of YourBase,
the startup I worked at at the time.

# Part 1

YourBase,
which I had joined eight months prior,
had a mission to accelerate software testing speed for monolithic applications.
When I joined they already had some really clever techniques for doing this.

The bulk of the magic happened in some low-level tracing code written in C and Rust,
which operated at the kernel level and so required that we control the build machine.
The magic code would track which specific test was currently executing,
and inspect the function calls made during the life of the test.
If a function was called during a test,
that test was said to depend on that function.
And if no function depended on by a test had changed since last time the test was run,
that test could be completely skipped.
There are some edge cases to this like environment variables and opened file pointers,
but that's the simple explanation.

That we need control of the build machine necessitated that the build machine be part of a CI pipeline,
which the company had custom-built.
This meant on top of the build tracing tech we needed to spin our own CI infrastructure from scratch,
including UI,
user and role management,
log streaming,
etc.

Worse,
it was a monumental effort to convince customers to move entire CI systems,
especially when our target customers by definition had large and monolithic code bases;
there was little expectation that any would be nimble enough to swap entire CI platforms.
Hindsight is 20/20, but it's clear now that it took a special type of company to agree to become a customer of ours.

Our salesperson had identified this problem,
and was desperate for a solution.

## Part 2

On that fateful post-holiday Sunday,
I had an itch.
Our kernel-level tricks were brilliant and effective,
but the effort of onboarding was stopping most prospective customers from even _seeing_ the benefit,
let alone diving head first into swapping over to it.

I began wondering,
could I get an 80/20 tradeoff here?
Are there enough tracing tools,
or at least undocumented hooks,
in Python for me to reproduce even a fraction of our C code's effectiveness?

Through some research and thought experiments I arrived at an unexpected starting point:
code coverage.
Clearly, the existence of any tool that generates code coverage reports—line-by-line analysis of which lines of code that ran during a test—implied the existence of tracing tools friendly enough to allow another body of code to skip tests.
