---
filename: death.html
date: 2022-05
shortname: TODO
---

---
filename: death.html
date: 2022-05
---

# Death Planning for Engineers
In Sweden there is a concept called döstädning, or death cleaning. This is the practice of cleaning your house and belongings so that if you died tomorrow, your next-of-kin would have a pleasant, not burdensome, time cleaning up your things. It can be used in the contexts of end-of-life planning, sure, but I’ve found it useful even as an exercise in increasing what I’d call “meaning density” among my things.

"End of life planning", or continuity, is something familiar to the professions of law and medicine. But these professions deal with inheritances and end-of life care, not with day-to-day items like your six spatulas or the scattered change in your kitchen drawer.

As an engineer I want to have the more technical aspects of my life handled in a more technical manner.

My [side][splits.io] [projects][whatsinstandard], for example, serve a combined ~25k users per month. I would hate to get hit by a bus and have their respective stacks slowly fall over as bills go unpaid and content stops getting updated.

At the same time, I don't want to burden others with the responsibility of caring for a project they don't care about. For these reasons, I've decided to get serious about "continuity planning" for the tech I've written after I'm gone.

[splits.io]: https://splits.io
[whatsinstandard.com]: https://whatsinstandard.com

# Access Needed by Successors
- Code
- Infrastructure 
	- 
- Money
	- Discovery: where is our money?
	- Instruction: why was I doing what I was doing with our money?
	- Access: unneeded because law handles this, but can be expedited by me

# Storage
I'm not looking to keep an ever-growing file folder of 2FA backup codes.

# Security
One of the most interesting questions is how to give someone full administrator access after you've gone, but not before. Lawyers and wills are great at *delivering* postmortem information reliably, but I wouldn't trust that system to *secure* it reliably, at least to my standards.

We will need automated mechanism to release the keys. It's okay for this mechanism to be triggered by law such as a will, but not for it to be easy to silently steal (such as a private key written on a piece of paper in a filing cabinet).

## Impersonation vs. authentication
There is the choice of whether to deliver your own private keys and passwords to a trustee, or to do it the "proper" way by promoting their own account to administrator and then self-destructing yours.

Although proper authentication feels like the correct way to do things, it comes with significant initial & ongoing effort, and adds a whole lot of complexity, for not much tangible benefit other than improved UX for the trustee. The complexity is reason enough to say no here—we want our delivery to be as simple as possible to reduce the chance of errors.

## 

# Testing

# Deploying
## How to kick off a deploy
# Health Checks