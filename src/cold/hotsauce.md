---
date: 2022-06-05
filename: hotsauce.html
type: post
---

# When the April Fool's Day Prank is Real

I love it when companies call their own bluffs on April Fool's Day. I like to
click into their pranks to see where the seam is. Are they "launching a
new product" but the only page to be found is a press release? Or, is
there a product page with a disabled Buy button? Sometimes the Buy button
works, but brings you to a page that says "Fooled you!".

I was never satisfied by these seams after witnessing first-hand the gold wacky
standard: On April 1, 2004, a popular search engine pranked that it was pivoting
into becoming an email service -- and then [it
did](http://googlepress.blogspot.com/2004/04/google-gets-message-launches-gmail.html).
This was especially surprising in the age before every startup was a data
company.

Of course, Gmail was not created to remove a seam from an April Fool's Day
prank. It was made independently from April Fool’s Day, and someone decided it
would be funny to disguise its launch as a prank. But as a bystander the
effect is the same: it seems that they are calling their own bluffs. There's
no "Fooled you!" page.

## Calling Your Own Bluffs

I always wanted more experiences like this, so I built one. I run a website
called [Splits.io](https://splits.io) whose goal is to help
speedrunners (people who time themselves playing video games) play faster.
It's a niche site, but we have paying customers and about 10k monthly active
users.

For April Fool's Day in 2020, we launched a full-page takeover introducing a
feature to help people speedrun faster: the urgent need to use the bathroom.

We joked that we would ship them a spicy hot sauce which they could use as a
sort of gastrointestinal negative reinforcement. But when people clicked Buy,
they did not see "Fooled you!", but a form to enter their payment
information. And when they entered their payment information and pressed
submit, they received a receipt in their inbox. And a charge to their card. A
week later, a package would show up at their door.
We called all our bluffs.

## The Set Up

I knew I wanted people to receive real hot sauce, and I knew it couldn't just
be any hot sauce from the grocery store. The April Fool's takeover was going
to advertise a Splits.io-branded hot sauce, and we were going to give it to
them. I called it _Race to the Bottom_.

### The Sauce

If you don't want to deal with actually making or bottling a product like hot
sauce but you still want one sporting your own brand, the main option is to
work with a
**white label** company, i.e. one that does every step of making a
traditional product except for the branding. Instead of selling to consumers
or retailers, they sell to those with a brand and nothing else. This could be
someone who's just amazing at branding and needs a supplier, or it could be
your employer making some commemorative swag, or it could be an engineer
trying to make a dumb April Fool's Day joke.

A surprising number of household brands use white label products behind the
scenes, namely store brands. Because the customers of these white label
companies are often large businesses, there's not a pressing need for
self-service like in consumer retail. That means it can take some phone tag
and email threads to get the relationship hammered out. But still, it works --
19% of the US economy is made up of white label goods
([source](https://www.statista.com/topics/1076/private-label-market/#topicHeader__wrapper)).

Luckily the last few years have seen some self-service white label products
enter the market, propelled I think by the demand for startup swag. For
something as niche as hot sauce, though, the best we could hope for was
partial automation. I engaged with several vendors and ordered sample packs
from two -- one via self-service for $25 and one via phone call for free.

Then I made some chips, invited some friends over (pre-pandemic), and got to
work.

{{ img
  "Taste testing hot sauces."
  "trialing"
  "several hot sauce bottles around a plate, with part of each bottle poured into a spot on the plate that allows for dipping"
}}

We selected a hot sauce called Pale Ale Chipotle and the company rated its
spiciness as 8/10 (no Scoville rating).

### The Label

We're a small team with no designer, so I cracked open Acorn and spent way too
long making way too ugly a label.

{{ img
  "The final label for Race to the Bottom. We were also trialing an energy drink, but scrapped it at the last minute."
  "label"
  "a hot sauce label reading 'Race to the bottom: A speedrunner's hot sauce' in the Splits.io colors"
}}

The white label company sent me a template with requirements such as how to
print the volume of the sauce within compliance, and the right proportions for
the label. They took care of the nutrition facts and bar code.

## Launch Day

We set up everything to happen automatically. The code had been deployed days
earlier, with JavaScript just waiting for the day to tick over before doing
the takeover. The announcement tweet was scheduled ahead of time as well.
The response when the clock struck midnight was exactly as we hoped.

> its actually real wtf

&mdash; Dark (@DarkRTA)
[April 1, 2020](https://twitter.com/DarkRTA/status/1245495897875759108)

{{ imgs
  "Feedback during and after launch."
  "feedback-1"
  "Screenshot a Discord channel where someone asks 'it is at this point that I begin to wonder how far I'll go for the meme, with a screenshot of them at the last step of the hot sauce payment flow, with an unclicked button that says 'Pay'. Below that, they say 'I had to know' and attached another screenshot of the order being completed. They then said 'this is gonna be like that movie Accepted where you accidentally make a clickable button on a fake website and then end up learning how to meet the demand, isn't it?'"
  "feedback-2"
  "Screenshot a Discord channel where someone asks 'did anyone get the item' and then someone else replies with a photo of their hot sauce bottle saying 'Yep lmao'. A third person says 'whoa you got yours already? nice. i can't believe it's real lol. @Glacials the madman'"
}}

### Shipping

I went into this expecting to make bulk trips to the post office, but I
actually never had to leave the house -- usps.com allows you to schedule a
mail carrier to pick up packages during their mail route. I had padded
flat-rate envelopes ready, also delivered from their website. On April 1st,
here was my pipeline:

1. (Trigger) Receive “payment completed” email from Stripe</li>
2. Pay for a shipping label at usps.com (~$9)</li>
3. Print shipping label, affix to a padded envelope with scotch tape</li>
4. Pack in hot sauce and extra packing material</li>
5. Schedule a pickup at usps.com</li>
6. (On pickup day) Put queued package(s) on the porch</li>

## Final Surprises

I built the full-page takeover to check if the day was April 1 when deciding
whether to take over or lie dormant, but I didn't think about checking the
year. So one year later, we suddenly got another influx of orders. Luckily I
still had plenty of stock with expiration dates still over a year away, so
spent a little surprise day packing hot sauce again -- and again a third time,
another year later.

## Finances

After the sample packs and the bulk order of 48 bottles (the minimum), upfront
costs came to $$192.05. Over three April Fool's Days we sold 28 bottles in 25
orders at $$15 / bottle.

The cost of my labor and the boon to the brand are harder to calculate. I'll
call them a wash.

```plain
Revenue
Orders: $15*28   =  $ 420.00

Costs
Samples:           ($  25.00)
Hot sauce:         ($ 167.05)
Shipping: $9*25  = ($ 225.00)

Profit/loss ----------------
                    $   2.95
```

I'll take it.

## Conclusion

Getting paid $2.95 to live out a minor dream of building a seamless April Fool's
Day joke was a great deal. I'd have paid a lot more just for the joy of seeing
people's reactions in real time.

On the way, I learned how white labels work and how big a slice of the economy
they make up, which has changed how I think about brands. Sometimes you have to
be a little wacky to learn something new about the world.
