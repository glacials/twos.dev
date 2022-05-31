---
filename: death.html
date: 2022-05
---

# Death Cleaning for Engineers

In Sweden there is a concept called döstädning, or death cleaning. This is the practice of cleaning your house and belongings so that if you die tomorrow, your next-of-kin would have a pleasant time, not a burdensome one, cleaning out your house. It can be used in the context of end-of-life planning, but I’ve found it useful as an exercise in increasing what I’d call “meaning density” in my possessions.

"End of life planning", or continuity, is something familiar to the professions of law and medicine. But these professions deal with inheritances and end-of life care, not day-to-day items like my six spatulas or the scattered change in my kitchen drawer.

As an engineer, big parts of my life exist digitally. I’d like to hand these off—or discontinue them—in a way that’s not burdensome, when the time comes.

My [side](https://splits.io) [projects](https://whatsinstandard.com), for example, serve ~25k users per month. I would hate to get hit by a bus and have them slowly fall over as bills go unpaid and infrastructure rots. However, I’m not looking to lay a ton of responsibility at someone’s feet.

## Assets to Consider

- Side projects
  - Infrastructure
  - Code
- Money
  - Discovery: where is it?
  - Instruction: why?
  - Access: how does one get at it?

# Storage

I'm not looking to keep an ever-growing file folder of 2FA backup codes.

# Security

One of the most interesting questions is how to give someone full administrator access after you've gone and not before. Lawyers and wills are great at _delivering_ postmortem information reliably, but I wouldn't trust that system to _secure_ it reliably, at least to my standards.

We will need automated mechanism to release the keys. It's okay for this mechanism to be triggered by law such as a will, but not for it to be easy to silently steal (such as a private key written on a piece of paper in a filing cabinet).

## Impersonation vs. authentication

There is the choice of whether to deliver your own private keys and passwords to a trustee, or to do it the "proper" way by promoting their own account to administrator and then self-destructing yours.

Although proper authentication feels like the correct way to do things, it comes with significant initial & ongoing effort, and adds a whole lot of complexity, for not much tangible benefit other than improved UX for the trustee. The complexity is reason enough to say no here—we want our delivery to be as simple as possible to reduce the chance of errors.

## Deploying

Triggering a deploy sounds difficult at first (how can a system detect my death?) but remember humans can be a part of this system. Setting up a deploy trigger can be as lightweight as telling someone you trust the location of a sealed envelope which should be opened in the event of your death. Inside it, place a letter that asks them to kick off other processes. For extra failover, make two envelopes in two locations for two people.

##

## Health Checks

A convenient time to run health checks against your pipeline is when you get a new computer. Instead of importing things normally, open your envelope digitally and try to “recover” your accounts.
