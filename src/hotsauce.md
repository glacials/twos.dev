---
filename: hotsauce.html
date: 2022-06-04
---

# When the April Fool's Day Prank is Real

I love it when companies call their own bluffs on April Fool's Day. I like to click into their pranks to see where the seam is. Are they "launching a new product" but the only page to be found is a press release? Or, is there a product page with a disabled Buy button? If they go far the Buy button works, but brings you to a page that says "Fooled you!".

I was never satisfied by these seams after witnessing first-hand the gold standard: On April 1, 2004, a popular search engine pranked that it was pivoting into becoming an email service—and then [it did](http://googlepress.blogspot.com/2004/04/google-gets-message-launches-gmail.html). This was especially surprising in the age before every startup was a data company.

Of course, Gmail was not created to remove a seam in an April Fool's Day prank. It was made independently for April Fool’s, and someone decided it would be funny to disguise its launch as a prank. But as a witness the effect was the same; the bluff gets called again and again. There's never a "Fooled you!" page.

## Calling Your Bluffs

I always wanted more experiences like this, so I built one. I run a website called [Splits.io](https://splits.io) whose goal is to help speedrunners (people who time themselves playing video games) play faster. It's a niche site, but we have paying customers and about 10k monthly active users.

For April Fool's Day in 2020, we launched a full-page takeover introducing a feature to help people speedrun faster: the urgent need to use the bathroom.

We joked that we would ship them a spicy hot sauce which they could use as a sort of gastrointestinal negative reinforcement. But when people clicked Buy, they did not see "Fooled you!", but a form to enter their payment information. And when they entered their payment information and pressed submit, they received a receipt in their inbox.  And a charge to their card. A week later, a package would show up at their door.

We called all our bluffs.

## The Set Up

I knew I wanted people to receive real hot sauce, and I knew it couldn't just be any hot sauce from the grocery store. The April Fool's takeover was going to advertise a Splits.io-branded hot sauce, and we were going to give it to them. I called it _Race to the Bottom_.

### The Sauce

If you don't know anything about making or bottling hot sauce but you still want one branded with your own logo or whatever, you probably want to work with a _white label company_, a company that already sources ingredients, creates recipes, cooks, and bottles hot sauce, but stops just short of slapping a label on it.

A surprising number of brands you use every day operate via white label companies like this, namely store brands. The customers in these relationships are often large businesses, so there's not as much self-service as consumer retail enjoys. That means phone calls and email threads to get things hammered down.

Luckily this is starting to change; swag vendors are the farthest ahead here, but for hot sauces the best we could hope for was partial automation. I engaged with several vendors and ordered sample packs from two -- one self-service for $25 and one via phone call for free.

Next I made some fries, invited some friends over (pre-pandemic), and got to work.

[insert image - see favorites in Photos]

The hot sauce we selected was called Pale Ale Chipotle and had an 8/10 heat score (scoville rating nowhere to be found).

### The Label

We're a small team with no designer, so I cracked open Acorn and spent way too long making way too ugly a label.

[todo - label pic]

## Launch Day

[todo: show tweets from others]

### Shipping

I went into this expecting to make bulk trips to the post office, but I actually never had to leave the house—usps.com allows you to schedule a mail carrier to pick up packages during their mail route. I ordered padded flat-rate envelopes from their website. On April 1st, here was my pipeline:

1. (Trigger) Receive “payment completed” email from Stripe
2. Pay for a shipping label at usps.com (~$9)
3. Print shipping label, affix to a padded envelope with scotch tape
4. Pack in hot sauce and extra packing material
5. Schedule a pickup at usps.com
6. (On pickup day) Put queued package(s) on the porch

## Final Surprises

I built the full-page takeover to check if the day was April 1 when deciding whether to take over or lie dormant, but I didn't think about checking the year. So one year later, we suddenly got another influx of orders. Luckily I still had plenty of stock with expiration dates still over a year away, so spent a little surprise day packing hot sauce again -- and again a third time, another year later.

# Finances

After the sample packs and the bulk order of 48 bottles (the minimum), upfront costs came to $192.05. Shipping was $9 per package. Over three April Fool's Days we sold 28 bottles among 25 packages at $15 / bottle.

The cost of my labor and the boon to the brand are harder to calculate. I'll call them a wash.

```plain
Revenue
Orders: $15*28 =    $420.00

Costs
Samples:           ($ 25.00)
Hot sauce:         ($167.05)
Shipping: $9*25  = ($225.00)

Profit/loss ----------------
                    $  2.95
```

I'll take it.