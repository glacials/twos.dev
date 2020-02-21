<template>
  <b-container>
    <h2 class="title">Job history</h2>
    <b-media>
      <template v-slot:aside>
        <img src="../assets/splitsio.svg" class="company-logo d-none d-md-flex" alt="Splits.io logo" />
      </template>
      <h5 class="mb-0 justify-content-between d-flex">
        <span>
          <b-link href="https://splits.io" class="text-dark">Splits.io</b-link>
          <small class="text-muted">2019 - present</small>
          <small class="text-secondary d-block">Founder & CEO</small>
        </span>
        <img src="../assets/splitsio.svg" class="company-logo d-flex d-md-none" alt="Splits.io logo" style="height: 2em; width: 2em" />
      </h5>
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
            in 2014, it's in the early stages. Much of our effort is on growing the market itself.
          </b-card-text>
        </b-card>
        <b-card title="Building an open startup">
          <b-card-text>
            Our core product is open source, forcing an interesting balancing act of risk against community
            goodwill. Because we're in a grassroots market, openness is valued highly by our customers.
          </b-card-text>
        </b-card>
      </b-card-group>
    </b-media>
    <b-media>
      <template v-slot:aside>
        <img src="../assets/twitch.svg" class="company-logo d-none d-md-flex" alt="Twitch logo" />
      </template>
      <h5 class="mb-0 justify-content-between d-flex">
        <span>
          <b-link href="https://twitch.tv" class="text-dark">Twitch</b-link>
          <small class="text-muted">2014 - 2019</small>
          <small class="text-secondary d-block">Senior Software Engineer</small>
        </span>
        <img src="../assets/twitch.svg" class="company-logo d-flex d-md-none" alt="Twitch logo" style="height: 2em; width: 2em" />
      </h5>
      <p>
        As the 91st employee at Twitch I built many formational tools and products, including ones that directly led to
        acquisition, and others that followed up on integrating into Amazon's tech stack.
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
            One night at 9pm three team members and I received a business-critical work text. Over the next 14
            hours we built and launched host mode from scratch, a feature still loved by users today.
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
            I built the first spam filter on Twitch, still in operation today. I can't really talk about it
            because spam filters are one of those things that stop working when you talk about them.
          </b-card-text>
        </b-card>
      </b-card-group>
      <b-card-group>
        <b-card title="Social merging">
          <b-card-text>
            When Twitch acquired Curse, a company with social features nearly identical to ours, I merged the
            code, products, and data powering social features like Friends without downtime for either platform.
          </b-card-text>
        </b-card>
        <b-card title="Mentorship">
          <b-card-text>
            By the time I left Twitch, I was mentoring two full-time engineers&mdash;one new grad and one self-taught
            internal transfer from customer support. Both had unique challenges in mentorship, and I'm proud to say both
            are still thriving in their roles today.
          </b-card-text>
        </b-card>
      </b-card-group>
      <h5 class="mt-4">...and a lot more.</h5>
      <p>
        Especially while small, Twitch consumed a large part of my life and I spent blood, sweat, and tears building and
        fixing things I'd been passionate about as a user without guidance. Some made <b-link
        href="https://twitter.com/glcls/status/720689621466619904">big splashes</b-link>; others hummed silently along
        forever. 
      </p>
    </b-media>
    <b-media>
      <template v-slot:aside>
        <img src="../assets/llnl.gif" class="company-logo d-none d-md-flex" alt="Lawrence Livermore National Laboratory logo" />
      </template>
      <h5 class="mb-0 justify-content-between d-flex">
        <span>
          <b-link href="https://llnl.gov" class="text-dark">Lawrence Livermore National Laboratory</b-link>
          <small class="text-muted">2010 - 2013</small>
          <small class="text-secondary d-block">Software Engineering Intern</small>
        </span>
        <img src="../assets/llnl.gif" class="company-logo d-flex d-md-none" alt="Lawrence Livermore National Laboratory logo" style="height: 2em; width: 2em" />
      </h5>
      <p>
        As a repeat intern at LLNL, I worked both on closed source tools for the communications systems team and later
        <b-link href="https://github.com/esgf">open source tools</b-link> for the climate science team.
      </p>
    </b-media>
  </b-container>
</template>

<script>
export default {
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
