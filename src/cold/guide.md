---
date: 2024-08-05
filename: guide.html
type: page
updated: 2023-03-28
---

# User Guide

For the benefit of my coworkers and future coworkers,
this is how I work.
This is descriptive of me, not prescriptive to you;
treat it as a set of FYIs, not requirements.
Use or ignore at your discretion,
and share your own guide with me!

## Meetings
Because of my schedule I'm not a fan of morning meetings,
especially recurring ones. I can do one-offs when needed.

I like meetings with well-defined outputs
("At the end of this meeting, we will have decided…").
They help me prepare and keep the meeting on-task.

If the meeting output is arrived at quickly,
I would rather end the meeting early than use all time available.

I prefer asynchrony over ad hoc, ad hoc over scheduled:

<pre aria-label="A flow chart that details my communicaation preferences. It starts with messaging (e.g. Slack). That flows to a recording (e.g. Loom) if the explanation is difficult in text or if it is a presentation or demo. From there flows to an ad hoc meeting if there is lots of back-and-forth; the first node (messaging) also flows to ad hoc meeting if there is a lot of back and forth. Ad hoc meeting flows to scheduled meeting if the participants are not online at the same time. Scheduled meeting flows to recurring scheduled meeting if all of the above consistently happen on a recurring basis." style="line-height: 1rem; font-size: 0.8rem">
                            ┌───────┐
                            │ Start │
                            └───────┘
                               │
                      Need for communication
                               │
                               ▼
                         ┌────────────┐
     Presentation        │  Message   │
   ┌───or demo───────────│(e.g. Slack)│─────────────────────┐
   │                     └────────────┘   Lots of back-and-forth
   │                            │                           │
   │                            │                           │
   │                            │                           │
   │            Explanation     │                           │
   │         difficult in text──┘                           │
   │         │                                              │
   │         │                                              │
   │         │                                              │
   │         ▼                                              ▼
   │  ┌─────────────┐                                  ┌─────────┐
   │  │  Recording  │     Lots of back-and-forth       │ Ad hoc  │
   └─>│ (e.g. Loom) │─────────────────────────────────>│ meeting │
      └─────────────┘                                  └─────────┘
                                                            │
                                                            │
                                                      Not online at
                                                      the same time
                                                            │
                                                            │
                                                            ▼
                                                      ┌───────────┐
                                                      │ Scheduled │
                                                      │  meeting  │
                                                      └───────────┘
                                                            │
                                                            │
                                                      All above happen
                                                     consistently on a
                                                      recurring basis
                                                            │
                                                            │
                                                            ▼
                                                      ┌───────────┐
                                                      │ Scheduled │
                                                      │ recurring │
                                                      │  meeting  │
                                                      └───────────┘
</pre>

## Schedule
My body demands a lot of sleep.
By default I do not use alarms and let my body wake itself when it's ready.
After my morning routine,
the bell curve for my coming online has a median around 9:30AM pacific,
with a ~45 min standard deviation.

Inspiration strikes me randomly—sometimes during a workday,
but also sometimes on a Sunday, or at 4 AM, or for 16 hours straight.
I'm very productive during these random bursts,
but having a habit of overworking leads me to burnout.

To counterbalance, I shorten my standard workday a bit and the math roughly works out.
If you see emails or commits from me with unhealthy-feeling timestamps,
keep in mind my default workday is at ~70% capacity.

## My Promises to You

1. I will respond to you.
2. I will respond to you within one business day.
3. I promise to put my ego in the back seat to learning.
4. I promise to be receptive and thankful when you give me a piece of feedback.
5. I promise to share the _whys_ behind any feedback I give (whether personal, professional, or code).
6. I promise to be supportive and (if you want) helpful when you tell me about something that makes you feel uncomfortable.

## My Weaknesses

- I am weak to heat. My mood and productivity negatively coorelate with the temperature outside.
- Social interaction [costs me brainpower](autism.html#masking)
  so I am better at problem-solving when I have time alone to digest.
- I don't handle interruptions well, regardless of duration.^[[twos.dev/img/guide-programmerinterrupted.png](img/guide-programmerinterrupted.png)]^[[paulgraham.com/makersschedule.html](http://www.paulgraham.com/makersschedule.html)]
  Notifications for my Slack are often off while I am in flow.
  (If it's urgent, use the on-call flow.)
- I go rogue.
  Sometimes I'm so passionate about fixing or building something that
  I'll do it without tracking it because the overhead—or its emotional burden—
  is bigger than the job itself.
  I'll still run it by people before shipping.
<!-- Commented because the above was migrated from "Communication",
     but it doesn't make sense to brag about this here in Weaknesses.
  Some
  [adored features](https://twitter.com/search?q=https%3A%2F%2Ftwitter.com%2Fglcls%2Fstatus%2F720689621466619904&src=typed_query)
  have come out of having the freedom to do this;
  see `1` above.
-->


## Communication

<!-- Commenting because I like this idea, but it just doesn't belong here. Maybe somewhere else.
   **Succinct**: The more people expected to read what I'm writing, the higher the
   cost/benefit of spending time honing it. For widespread pieces, 90% of my
   time is editing. [More](http://www.paulgraham.com/simply.html), [even
   more](http://www.paulgraham.com/useful.html)
-->

#### 1. I would rather align on vision than work in lock-step.

Coordination is slow. As long as we share a common vision,
I trust you to make
[two-way door](https://shit.management/one-way-and-two-way-door-decisions/)
decisions in motion without blocking on consensus.
I ask that you trust me to do the same.
We can continue to share progress in nonblocking ways.^[[communitywiki.org/wiki/DoOcracy](https://communitywiki.org/wiki/DoOcracy)]


#### 2. I love when information has a home.

I'm forgetful.
I can more easily (re)gain context for an issue when discussion is concentrated in one place,
like comments on the relevant ticket or pull request.
I will leave breadcrumbs in these places for posterity if I feel things are missing.

#### 3. I focus on customer problems.

Even if a solution "feels right" be prepared for me to ask what customer problem it solves.

#### 4. [Don't say hello](https://nohello.net/).

<!-- Commented to help build this a bit before publishing.
## Interesting Things About Me

- As a kid, I was selectively mute for 9 years.
-->
---

_This page is in the form of
[Abby Falik's user guide](https://www.linkedin.com/pulse/leaders-need-user-manuals-what-i-learned-writing-mine-abby-falik/)
and inspired by [Dina Levitan](http://dinalevitan.com/)'s user guide._
