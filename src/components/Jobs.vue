<template>
  <b-container>
    <h2 class="title">Job history</h2>
    <!-- Adding a new entry? Remember to update your resume too! -->
    <CareerEntry
      name="YourBase"
      url="https://yourbase.io/"
      logo="../assets/images/yourbase.svg"
      dates="2020 - present"
      role="Senior Software Engineer"
    >
      <p>
        I started in April 2020 as the fifth employee at YourBase (now nine), a
        simplified build system and accelerator for complex code bases.
      </p>
      <b-card-group>
        <b-card title="Delving into RSpec">
          <b-card-text>
            So far at YourBase I've dived deep into the internals of RSpec to
            update our test selection tech for Ruby, which infers dependency
            graphs to later skip tests that don't need to run.
          </b-card-text>
        </b-card>
        <b-card title="Optimizing APIs">
          <b-card-text>
            Using knowledge I gained from years working with ORMs, I've
            decreased latency across the board for the YourBase API by over an
            order of magnitude.
          </b-card-text>
        </b-card>
      </b-card-group>
    </CareerEntry>
    <CareerEntry
      name="Splits.io"
      url="https://splits.io/"
      logo="../assets/images/splitsio.svg"
      dates="2019 - 2020"
      role="Founder & CEO"
    >
      <p>
        After five years of bootstrapping to 10k monthly active users as a side project, I formed a company, convinced
        two people to come on board, and went full-time on my pet project: building the Strava of esports.
      </p>
      <b-card-group>
        <b-card title="Wearing many hats">
          <b-card-text>
            My biggest impact as CEO is to fill in the cracks. From product to social media to swag to customer support,
            every role I step into improves my perspective and the company's trajectory.
          </b-card-text>
        </b-card>
        <b-card title="Growing talent">
          <b-card-text>
            Having two reports (one non-engineering) is a new challenge for me. With no manager above me, I research on
            my own how to manage them, grow their talents, and improve their careers.
          </b-card-text>
        </b-card>
      </b-card-group>
      <b-card-group>
        <b-card title="Finding the market">
          <b-card-text>
            Our segment of esports is niche, with a market size of only a few hundred thousand. However like streaming
            in 2014, it's in the early stages. Much of our effort is in growing the market itself.
          </b-card-text>
        </b-card>
        <b-card title="Building an open startup">
          <b-card-text>
            Our core product is open source, forcing an interesting balancing act of risk against community goodwill.
            Because we're in a grassroots market, openness is valued highly by our customers.
          </b-card-text>
        </b-card>
      </b-card-group>
    </CareerEntry>
    <CareerEntry
      name="Twitch"
      url="https://twitch.tv"
      logo="../assets/images/twitch.svg"
      dates="2014 - 2019"
      role="Senior Software Engineer"
    >
      <p>
        As the 91st employee at Twitch (YC W07) I built many formational tools and products—including ones that directly led to
        its $1B acquisition by Amazon, and others that followed up on integrating into the Amazon stack—before leaving at
        around 1800 employees.
      </p>
      <b-card-group>
        <b-card title="Localization">
          <b-card-text>
            I championed and led a public rework of how localization is handled at Twitch. You can read my write-up
            <a href="https://blog.twitch.tv/the-innards-of-our-new-internationalization-systems-e84f99a27526">
            at the Twitch blog</a>.
          </b-card-text>
        </b-card>
        <b-card title="Host mode">
          <b-card-text>
            One night at 9pm three team members and I received a business-critical work text. Over the next 14 hours we
            built and launched host mode from scratch, a feature still loved by users today.
          </b-card-text>
        </b-card>
      </b-card-group>
      <b-card-group>
        <b-card title="Auth">
          <b-card-text>
            Following a
            <a href="https://blog.twitch.tv/important-notice-about-your-twitch-account-aaa876745856">security
            breach</a> I and two coworkers founded the Identity team where we replaced our authentication infrastructure
            with an OpenID Connect implementation, still in service today.
          </b-card-text>
        </b-card>
        <b-card title="Spam">
          <b-card-text>
            I built the first spam filter on Twitch, still in operation today. I can't really talk about it because spam
            filters are one of those things that stop working when you talk about them.
          </b-card-text>
        </b-card>
      </b-card-group>
      <b-card-group>
        <b-card title="Social merging">
          <b-card-text>
            When Twitch acquired Curse, a company with social features nearly identical to ours, I merged the code,
            products, and data powering social features like Friends without downtime for either platform.
          </b-card-text>
        </b-card>
        <b-card title="Mentorship">
          <b-card-text>
            By the time I left Twitch I was mentoring two engineers, teaching an internal Go class I'd created for new
            hires, advising for an internal postmortem council, and training engineering teams to build scalable
            GraphQL.
          </b-card-text>
        </b-card>
      </b-card-group>
      <h5 class="mt-4">...and a lot more.</h5>
      <p>
        Especially while small, Twitch consumed a large part of my life and I
        spent blood, sweat, and tears working towards things I'd always wanted
        as a user, including open source involvement, running
        <b-link href="https://www.twitch.tv/videos/93570710">panels</b-link>,
        and small features with
        <b-link href="https://twitter.com/glcls/status/720689621466619904">big splashes</b-link>.
      </p>
    </CareerEntry>
    <CareerEntry
      name="Lawrence Livermore National Laboratory"
      url="https://llnl.gov"
      logo="../assets/images/llnl.gif"
      dates="2010 - 2013"
      role="Software Engineering Intern"
    >
      As a repeat intern at LLNL, I worked both on closed source tools for the communications systems team and later
      <b-link href="https://github.com/esgf">open source tools</b-link> for the climate science team.
    </CareerEntry>
  </b-container>
</template>

<script>
import CareerEntry from './CareerEntry.vue'

export default {
  components: {
    CareerEntry,
  },
  name: 'Jobs',
}
</script>

<style scoped>
.company-logo {
  width: 5em;
}

.media {
  margin-bottom: 3em;
}
</style>
