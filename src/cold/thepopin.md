---
date: 2021-05-01
filename: thepopin.html
preview: How to pop in on someone when fully distributed.
type: post
---

# The Remote Pop-In

This page describes a method of dealing with one of the missing features of
fully remote teams: the pop-in.

## The Pop-In

The pop-in, back before open floorplans, is when you physically walk into
someone's office to ask them a question or spitball an idea. In an open
floorplan you can instead swivel your chair or walk to their desk.

The pop-in is valuable because it's as informal as email, but as
high-bandwidth and low-latency as meetings. It's interruptive to be popped in
on, but that interruption power is self-limiting: the effort of physically
going somewhere is high, and the tools available to the target like doors,
signs, and social cues are numerous and don't require much thought to
implement.

Slack DMs come close to replacing the pop-in, but they reduce the effort too
much. This leads to more interruptions (a little lower in intensity, but a lot
higher in number) and even to the replacement of trackers, emails, and wikis
with DMs.

There's a lot wrong with the pop-in, but like any tool it has its uses. Many
great products—companies, even—have been built on the shoulders of
the pop-in.

## The Remote Pop-In

There's really no remote equivalent to the pop-in. There are products that
help you "hang out" like Facebook Portal, but nothing that implements the
pop-in.

At [YourBase](https://yourbase.io) we decided we needed one. Rather
than build another tool, we experimented with Discord. I'm happy to say we
found something that works, and even solves the worst problem of the
open-office pop-in: the inability to shut your damned door.

### The Setup

![A screenshot of a Discord server with several voice channels, some empty, some having one person in them.](/img/thepopin-0-dark.png)

_The company Discord server on a workday._

In our company Discord server, we give everyone one voice channel of their
own. This is your "office". When you are online and able to be disturbed, you
join your office, alone. You stay in there as long as you are "in".

Because of how Discord works, everyone can see who is in what office at all
times. But you can only hear conversations for the room you're in, and you can
only be in one room at a time. Joining a room plays an obvious audio chime and
is not something you can do sneakily.

When someone wants to pop in, they join your voice channel. The chime plays,
they say hello, you unmute and turn on video if you like, and have your
pop-in. And that's really it.

If you're going to the restroom and will be AFK, you can join the restroom
channel. You're in, but away. No one is actually wearing a headset on the
toilet, it's just a fun way to feel like a real office.

We've even fallen into some typical office behaviors:

- We've started having meetings in people's offices, just because it's easier.
- We have a "cafeteria" room where you can eat your lunch if you've got nothing going on in real life. We also hold lunch-and-learns here.
- We have a "rooftop" and other common spaces where some people chill, which
  signals that they're open for a chat.

  The oddest one is that when it's slow around the office, you can tell. You can
  almost hear how quiet it is when you're the first one in or the last one out.

  The pop-in is valuable because it's informal, high bandwidth, and low latency.
  Enabling it using Discord has worked very well for us over the last six
  months, and we didn't even have to build anything.

  Thanks to <a href="https://discordapp.com" target="_blank">Discord</a> for making such a
  flexible platform. Thanks to <a href="https://yourbase.io" target="_blank">YourBase</a> for
  being so open to crazy ideas. (Check us out if your CI is slow.)

  ["The Pop-In"](https://www.youtube.com/watch?v=KzOv2jrC1I8)
