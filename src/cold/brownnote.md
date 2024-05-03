---
date: 2024-05-02
filename: brownnote.html
type: post
---

<link rel="stylesheet" type="text/css" href="/devices.min.css" />

# Brown Note

I've been working on an iOS app recently,
which is new for me.
The app is called [Brown Note][appstore]
and it is a food journal designed to solve the biggest problem I have had with other food journal apps:
logging entries is so complex that my brain unconsciously stops doing it.

## Customer Need

I'm slow to realize biological correlations.
There are several weird things about my body,
so it's not always obvious that when I eat `A` it causes poor experience `B`.
Lots of things cause poor experience `B`,
and also poor experiences `C`, `D`, and `E`,
so I'm not always on the lookout.

It took me months just to realize that I didn't always have "stomach issues",
a phrase I'll use as a shorthand for some combination of urgent, painful, and/or frequent poops,
and that they might be fixable.

When my doctor suggested food journaling,
I was diligent about looking for a high-tech solution in the form of a mobile app,
and I found several.

Thus started a cycle:

1. I choose a fantastic mobile app to use for food journaling.
2. I use it to log every meal and poop.
3. ?
4. Weeks later, I realize I have stopped logging and completely forgetten about it.
5. Repeat from 1.

As I tried more apps this way,
and became more conscious over time about my falling off the wagon,
I realized my issue was the friction of log entry.
I would be eating among friends,
or sitting down in pain,
and a lightweight thought of "I should log this" would brush over my arm
only to be dismissed—"I'll log it when I'm free".
Only when later came, it was still too heavy a process
(and too delayed a gratification)
to convince myself to do.

I realized I needed to remove as much friction from the process as I could.
I would simplify my process so much that that light brushing over my arm would be a strong enough signal to log immediately.

## Product Discovery

I carried a pen and a folded up piece of paper in my pocket,
writing down everything that went into and out of me,
occasionally copying them into a spreadsheet by hand,
and doing some eyeball analysis—and later some Python scripting—to find problem foods.

I did this because the physical pen would literally poke me whenever I sit down to do my business,
but there was something magical about how the small size of the folded piece of paper forced each log entry to be brief.
That magic allowed me to build a routine around it.

Version 2 was opening the aforementioned spreadsheet on my phone and directly logging there,
but this slog quickly led me to version 3,
an iOS Shortcut placed on my Home Screen that prompts for a single line of text,
which it appends to the spreadsheet.

## Simple Works

All I'm entering in these lines is plain text like
"Good poop",
"Bad poop",
or a few words describing a food item.
In some iterations,
I didn't even include a date;
order was enough.

> - Bad poop
> - Mango lassi
> - Chinese
> - Taco Bell
> - CookUnity
> - Soylent
> - Good poop
> - Taco Bell
> - Whiskey
> - Good poop
> - Chinese
> - Bubble tea

_An excerpt from my v2 spreadsheet._

I'm not separating "foods" and "poops" into two categories of thing.
I'm not typing out a bunch of tags or ingredients with every food item.
If I zero in on a food that's too ambiguous—say, "pizza"—only then will I start logging its components:
cheese, bread, tomatoes.
Until and unless that happens,
simplicity wins.
Keeping the routine going is more important than anything.

## Into Swift

In the background of all this I'd become more and more interested in native iOS development,
for [several](/apple.html) [unrelated](https://whatsinstandard.com/) reasons.

After stalling out trying to make an unrelated app with a lot of complexity,
I figured this was the perfect do-over to start simple.
I started learning Swift and SwiftUI again, this time with SwiftData too.

The result is [Brown Note][appstore].

<div class="device device-iphone-14-pro" style="zoom: 0.65; margin-bottom: 2rem">
  <div class="device-frame">
    <img
    alt="A screenshot of Brown Note showing how a specific food, chai in this case, impacts poops, including percentages of good vs. bad poops that come after."
    class="device-screen"
    src="/img/brownnote-correlate.png"
    style="margin-top: 0"
    />
  </div>
  <div class="device-stripe"></div>
  <div class="device-header"></div>
  <div class="device-sensors"></div>
  <div class="device-btns"></div>
  <div class="device-power"></div>
</div>


It was insanely easy to get Brown Note to a stage where I can use it for personal logging.
The complexity all lies in the analytics and insights,
which aren't part of its everyday use cases and so didn't block me from onboarding myself as a user.


<div class="device device-iphone-14-pro" style="zoom: 0.65; margin-bottom: 2rem">
  <div class="device-frame">
    <img
    alt="A screenshot of Brown Note showing the user some meals they have tracked."
    class="device-screen"
    src="/img/brownnote-trackmeals.png"
    style="margin-top: 0"
    />
  </div>
  <div class="device-stripe"></div>
  <div class="device-header"></div>
  <div class="device-sensors"></div>
  <div class="device-btns"></div>
  <div class="device-power"></div>
</div>

This helped immensely with design,
as I would personally use the app several times a day,
noting every rough edge and feature need.

<div class="device device-iphone-14-pro" style="zoom: 0.65; margin-bottom: 2rem">
  <div class="device-frame">
    <img
    alt="A screenshot of Brown Note showing the user a poop they have tracked, and the estimated inputs (recent meals) to that poop."
    class="device-screen"
    src="/img/brownnote-trackpoops.png"
    style="margin-top: 0"
    />
  </div>
  <div class="device-stripe"></div>
  <div class="device-header"></div>
  <div class="device-sensors"></div>
  <div class="device-btns"></div>
  <div class="device-power"></div>
</div>

I'm still in constant development,
but the simple use cases are covered.
As I build out more advanced ones, I'm being cognizant about keeping the logging experience dead simple.

For example,
the app can understand food components
(e.g. `latte` contains `milk` and `coffee`)
but I'm shielding the logging experience from being impacted by it.
When the app needs more info,
it asks asynchronously.

<div class="device device-iphone-14-pro" style="zoom: 0.65; margin-bottom: 2rem">
  <div class="device-frame">
    <img
    alt="A screenshot of Brown Note asking the user for the ingredients of saag paneer."
    class="device-screen"
    src="/img/brownnote-honein.png"
    style="margin-top: 0"
    />
  </div>
  <div class="device-stripe"></div>
  <div class="device-header"></div>
  <div class="device-sensors"></div>
  <div class="device-btns"></div>
  <div class="device-power"></div>
</div>

Let me know what you think,
or join the [TestFlight][testflight] group to give me some feedback before proper release.

<!-- _Brown Note is available on the [App Store][appstore] for iOS._ -->

[appstore]: https://apps.apple.com/us/app/brown-note/id6479333983
[testflight]: https://testflight.apple.com/join/ww7RII5M
