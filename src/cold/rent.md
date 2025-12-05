---
date: 2025-12-04
filename: rent.html
preview: A simple framework for deciding who should pay what for which room.
type: post
---

# How to Split Rent

Given a home with $n$ unequal rooms and $n$ roommates with unique preferences,
what's the right way to split rent and assign rooms?

Should the room with the attached bathroom cost more, or the one with the balcony?

Should roommate A get the room with the view, or roommate B?

How much should each room cost?

Read below, or [skip to the calculator](#calculator) to enter your own numbers.

## The Algorithm

Each roommate privately writes down the price they’d be **happy to pay** for each bedroom,
making sure the numbers they choose sum to $Rent$.

For example,
in a 2-bedroom home with $Rent = 2000$,
I might write:

```plain
North bedroom: $1200
South bedroom: $800
```

i.e. I strongly prefer the north bedroom,
but I'd happily take the south bedroom if it were cheap.

If a coin were flipped to decide which room I get at the price I listed,
I'd be happy either way.
It is crucial that each roommate decide their numbers like this,
because that is functionally what will happen.^[If a person cannot make such a split, then either rent is out of proportion with their expectations or they have budget constraints that limit their options. This algorithm cannot solve for these, but a transparent discussion can.]

Gather all slips of paper and average the numbers for each room.

| Person  | North Bedroom | South Bedroom |
| ------- | ------------- | ------------- |
| Arnav   | $1000         | $1000         |
| Carol   | $1300         | $700          |
|         |               |               |
| Average | $1150         | $850          |

### Result

The highest bidder gets the room.
They pay the average.

| Person | Bedroom       | Willing to Pay | Will Pay |
| ------ | ------------- | -------------- | -------- |
| Arnav  | South Bedroom | $1000          | $850     |
| Carol  | North Bedroom | $1300          | $1150    |

The averages always sum to $Rent$.

Almost always, each person pays less than what they would have been happy with;
everyone feels like they got a deal.

This works for any number of roommates.

#### Edge Cases {#edge}

- If rounding down makes you come up short,
  someone needs to pony up the extra $1 or 1¢.
- If multiple people bid the same amount for a room,
  or otherwise if "highest bidder gets it" becomes unclear,
  [here's](/rent_algorithm.html) the nuanced algorithm.
- In extreme situations,
  a person may pay more for a room than their chosen price.
  I've never seen this happen.

## The Calculator {#calculator}

Remember the rules:

1. Each roommate should **privately** write down a price for each room that:
   - they would be happy to pay if assigned that room randomly, and
   - together sum to $Rent$.
2. Once everyone has written these numbers down,
   collect them and enter them below to see room assignments and prices.
3. Each roommate pays the $Will\ Pay$ amount for $Room$.
   The numbers will sum to $Rent$.

{{ template "_rent-widget-styles.html.tmpl" }}

<div class="rent-widget" id="rent-widget">
  <div class="rent-widget__controls">
    <label>
      Roommates
      <input type="number" min="1" max="8" value="3" data-rent-count />
    </label>
    <label>
      Rent
      <input
        type="number"
        min="0"
        step="10"
        value="2500"
        data-rent-total
      />
    </label>
    <button type="button" class="rent-widget__button" data-rent-reset>
      Reset to example
    </button>
  </div>
  <div class="rent-widget__table-wrap" id="rent-widget-grid"></div>
  <div class="rent-widget__messages" id="rent-widget-messages" aria-live="polite"></div>
  <div class="rent-widget__results" id="rent-widget-results" aria-live="polite"></div>
</div>
{{ template "_rent-widget.html.tmpl" }}

## Why Does This Work?

Each room is a one-of-a-kind product in a marketplace.
The marketplace is "rooms available in this home".

Each buyer values each room differently.
So we have to answer two questions:

1. Who gets which room?
2. How much does each room cost?

So we turn to the standard way to discover these things for a one-of-a-kind items.

### Auctions

But there are some weird constraints.
First, each roommate must buy exactly one room.

We can't hold a standard auction;
when any room is bid on,
the other rooms' prices would change.^[This can technically work,
but in practice would give a bad social foundation for a roommate relationship
Imagine having a roommate that pays $0 because the other roommates got into a bidding war.]
What we can do is hold $n$ auctions _simultaneously_.
Numbers are recorded privately and in one turn
to make it difficult to manipulate dependent prices.

The second constraint is that the sum of all rooms' prices must total to $Rent$.

If we were to simply charge the highest bidder their bid amount,
the total would almost always exceed $Rent$.
Some common strategies exist to lower the amount the winner pays,
such as having them pay the [second-highest amount](https://en.wikipedia.org/wiki/Generalized_second-price_auction).
But we need a guarantee of a specific sum, not just "lower than the sum of the highest bids".

This is why each person's chosen prices must sum to $Rent$.
If you have a 2-dimensional table of numbers where the sum of every row is the same number $T$,
then the sum of the averages of every column is also $T$.

This is called _linearity of addition_ and is how we can end up with a total of exactly $Rent$
without throwing a wrench in anyone's expectations of price.

## Conclusion

Everyone has different preferences when it comes to room selection.
I value sound isolation and light,
but you may value an attached bathroom or extra space.

Instead of picking randomly and silently settling to avoid tension,
consider using basic economics to do proper price discovery and hopefully everyone can be a bit happier with how it turns out.

I'd be delighted to hear if this helps you out.
Reach me at [ben@twos.dev](mailto:ben@twos.dev).
I'm not around much on social media these days.
