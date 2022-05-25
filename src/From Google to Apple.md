---
filename: apple.html
date: 2022-05
updated: 2022-05
---

# From Google to Apple

This is a living document.

My life is in Google more than anyone I know. The usual suspects, sure -- Gmail, Photos, Chrome, and the like, but also the ones people stay away from—Tasks, Keep, YouTube Music, Fit, Pay, Play Books, Play Movies & TV, Podcasts. Our house has more Google Home Minis than rooms. I was a beta user for the first Chromebook, the Cr-48. I have more scar tissue than most from services being sent to the “Google graveyard”. For a few years I was gung-ho enough to overlay my Facebook profile picture with a Google+ logo and the text “I’ve moved”. It’s no surprise I’ve been an Android user for 12 years, 10 of them with Nexus and Pixel devices.

Once in a while I challenge my beliefs by changing seats. This is how I understand why—or whether—I really prefer what I chose. But when I've previously trialed a borrowed iPhone for a few months (thanks [@bensw](https://twitter.com/bensw)), my thoughts amounted to "It's not that different, but the Google apps don't have OS integration".

I've realized since then it’s hard to know the Apple value proposition without being immersed in the ecosystem. For the same reason the Google Assistant doesn't shine iOS, iOS doesn't shine if you only use Google apps.

After a buggy summer with my Pixel 3a, I decided now’s a good time to give Apple an honest shot. The infamous walled garden may not let me out, but better to trap myself on one side and know both, than trap myself on the other and know one. These are my notes on hopping the fence.

### ArchLinux → macOS

I have used one or another Mac as a development machine since I switched from the Cr-48 in 2013, which was running ArchLinux at the time. Having gone straight from Windows to Linux due to a love of the POSIX command line and the simplicity of the system (e.g. the “everything is a file” notion), I found macOS to scratch the same itches minus the days spent customizing things.

### Nexus 7 → iPad

My minimum goal for a tablet is to replace all my at-home phone use. It gives an upgraded experience and as a bonus, drains a different battery bank than the one I may need later. As a habit, I carry it room-to-room.

With the Nexus 7, within a year I’d stopped that habit. It afforded too little benefit and took too much effort to maintain separately from my phone. Apps didn’t expect one user to be using two devices, so bugs were frequent. Almost zero apps used the tablet form factor well (or at all), making the UX worse than a phone. I expected a similar experience from iPad.

Instead, I’ve carried it room to room for two years now. It’s fulfilled the minimum goal and has even replaced some light computer usage, such as writing. The app ecosystem better understands the tablet form factor; many apps use it to display more information or at least differently laid out information. The multitasking experience is opinionated but helpful.

The ecosystem benefits I receive from iPad are the ability to use it as a second monitor, and the ability to use my Mac’s mouse and keyboard to control the iPad (a subtle but important distinction). macOS also allows iPad to be used as an input device for signing documents and marking up images. This is something I did before by manually transferring files, but now it’s two clicks with instant sync.

### Pixel → iPhone

iPhone makes clear the differences in the Google and Apple attitudes around user identity. In Google a user is a Google account; an Android phone is a window into it. In Apple a user is an iPhone, and the iPhone is possibly synchronized with an Apple account.

For example, in Google Docs documents are saved to your account, presumably to a row in a database somewhere. In Pages documents are saved to the filesystem, defaulting to the iCloud Drive directory, which synchronizes them to your Apple account. The filesystem is the database.

There’s a part of this that’s charming, but it’s not for everyone. I regularly experience synchronization delays and conflicts.

I was initially surprised iPhone and iPad lack feature parity. This first came up trying to swipe-type on iPad. Another is iPad can show multiple apps on the screen at a time, but iPhone cannot. iPad has no Wallet or Apple Pay. I felt confused at first, but I see Apple’s angle—none of these quite make sense on the other form factor. They’d rather deny you the feature than allow using it to be confusing or unintuitive. Infamously, there is no calculator app for iPad because Apple hasn’t found the time to adapt the iPhone calculator design.

The more native apps I use on iOS, the better the UX gets. For example, you can use multi-touch to drag-and-drop items between apps, such as from Photos into Mail. Third-party apps only implement this functionality sometimes.

#### Notifications 

Notifications on iOS behave differently than on Android, and that difference makes them feel a lot worse at first. After some time I’d call them a sidegrade.

Android notifications are stateful. If you receive a Gmail notification on your phone then read that email on any platform, the notification goes away. Otherwise, the notification will be there even days later. In this way, Android notifications can be used as a todo list of sorts.

iOS notifications are a feed. In the same scenario on iOS, the notification is (usually) never revoked remotely. Instead, when the phone is unlocked the notification moves to a secondary location “below” (via swipe) the lockscreen. I do inbox zero, so this is strictly worse. I forget to address notifications by committing the cardinal sin of unlocking my phone.

However, iOS mostly fills the gap with badges. Badges are red numbers in the corners of app icons that are stateful in the same way Android notifications are, so can be used as a todo list. I’ve found these a decent enough substitute, although I still prefer Android’s notifications.

#### iMessage

I finally know why none of my iPhone friends were as excited as me about “solving” mobile chat—it’s been solved for them for years. iMessage is the gold standard of chat apps. I’m even feeling the guilty urge to nudge friends towards iPhone so I can use it with them.

In classic Apple fashion iMessages aren’t between Apple accounts, but devices. I recently changed my email address but I need to keep my old one in my Apple account indefinitely so my friends’ threads with me won’t retroactively split into two.

### Chromecast/Google TV → Apple TV

We don’t have HomePods and my spouse forbids me to get any after filling our home with Google Home Minis, so we’ve had to migrate away from “hey Google, play $SHOW”. 

#### The Remote
Apple TV can initiate play from a phone like Chromecast, but beyond that we’ve used the remote. The touchpad on the remote is awkward and hard to use; I brush it when trying to click it, causing me to click on the wrong thing.

Otherwise it has a good build quality and feels good in the hand. But it’s a remote and can be lost. We’ve 3D printed a holder for it that attaches to the coffee table. I miss the remote-free life.

When a text input is selected on the Apple TV, a notification shows on my iOS devices allowing me to use them as a keyboard. Password manager support works as normal; this is helpful for invoking 1Password to fill logins.

### Home → Home

My biggest fear with this change was losing a decade of Google Nest ecosystem benefits, but running Homemanager on a Raspberry Pi made it all work seamlessly, even down to our Nest doorbell camera showing on our Apple TV picture-in-picture when the doorbell rings. (If you’re not keen to se up a Raspberry Pi, try Starling Home Hub instead.)

The UX of Shortcuts with Home is miles better than Google Home’s long device list that feels like a web page. This was a strict upgrade.

### Pixel Buds → AirPods Pro

Everything you’ve heard about the AirPods Pro noise cancellation is true. Transparency mode is so good that more than once I’ve forgotten they’re in my ears. They’ve,ore than replaced my Bose QuietComfort 35 IIs.

The ecosystem starts to shine here. When I am working at my Mac with my AirPods connected to it, then I walk away and start using my iPhone, they switch automatically.- 

### Gmail → iCloud Mail

The big one. I’ve been wanting to switch my email to a domain I control anyway, and iCloud+ supports that.

I don’t recommend importing Gmail archives into iCloud Mail; the experience was fraught with landmines and didn’t achieve the desired result. I’ve chosen to live the life of searching in two places when I need something.

The iCloud web interface is bad. On my Windows gaming computer I’ve installed Thunderbird to get by. Mail itself integrates very well with the rest of the ecosystem and has a solid UX at its core, but is plagued paradoxically by usability issues.

Mail's junk filter sorts ~one legitimate email per week into the junk folder, even after months of correcting it.

Mail's search is bad. On macOS, it will assume you want to search the current mailbox. Because I do inbox zero, this is never what I want to search. (iOS and iPadOS correctly search everything.) When you are done searching and click Inbox to "return home" the search field doesn't clear, leading you to believe your Inbox is empty.

Ignoring that, Gmail still wins search by an order of magnitude. For example, Gmail searches the contents of PDFs attached to emails; I’ve found this invaluable finding old leases and whatnot.

Mail is okay at displaying emails as conversations. Once in a while, it omits something it shouldn’t. When I click Send while replying to a thread it doesn't immediately append my message, causing me to believe something went wrong.

Mail is bad at dealing with calendar invitations. One invitation appears as two attachments: one ICS file and one generic "Mail Attachment". Each opens Calendar with an identical event draft, but I’ve learned to use “Mail Attachment” because it is a two-way connection with the sender; ICS is only a copy of the initial event details.

macOS Mail supports filters, but only locally. The iCloud Mail service has a separate filter system with the same effect, but the two have disparate feature sets and don't synchronize. iOS Mail does not support filters. I’ve opted to use Mail’s local-only filters, and have my iMac stay awake 24/7.

If you want a Mail-like experience with Gmail as a backend, I recommend [Mimestream](https://mimestream.com) instead; it is written by a former Apple engineer on Mail but uses Gmail APIs instead of SMTP.

### Calendar → Calendar

After migrating calendars but before transitioning my custom-domain email address to iCloud Mail, my email was serviced by Google Workspace. I had this Google Workspace account added to my iPhone and iPad. This got me into a weird state where Calendar would not let me add an event to my iCloud calendar if the invitee was my new email address, if that Google Workspace account was present on my device, but it would allow me to add it to that Google account's calendar. When I removed that Google account from my device, I could add the event to my iCloud calendar fine.

I don't recommend using the Week view in Calendar. Similar to Google Calendar, Calendar uses a horizontal red bar to represent the current time; however this red bar extends to occupy 100% of the width of the week, and does not do a good enough job of showing you which day today is. This has led me to mistake the wrong day of the week for today multiple times, inducing panic about being late for meetings and reaching out to people to reschedule other meetings. I now use Day view, where no red bar shows unless viewing today.

Quick look is a blast, the same way it is for previewing files in Finder.

### Google Maps → Apple Maps

Apple Maps gets a lot of flak for its initial release state, rightly so. But that was years ago, and they’ve kept at it. It’s improved a lot.

Unfortunately it’s only gotten good, not great. There are a lot of small UX benefits over Google Maps, such as the spoken direction “Go past this light, then at the next one, turn `$DIRECTION`”, but about 5% of my trips to new places end me at the right area, but wrong specific destination. I’ve been routed to the delivery entrance for a museum, the back entrance for an airport, and the wrong parking lot (15 minutes of walking wrong) of a large 360° shopping center.

I’m continuing to give it chances because I know more data helps, but for risky timelines I go back to Google Maps.

The ecosystem benefits of Apple Maps include rich links in iMessage and Calendar, Siri suggestions for destinations based on those and other sources, and Apple Watch vibration patterns to indicate upcoming turns.

### Chrome → Safari

This was a far less noticeable change than I expected. Everything from bookmark sync to my extensions to rendering works the same.

In macOS the Safari chrome fades into the background better than the Chrome chrome, such that websites feel a bit closer to standalone applications even when in a tabbed window. The touchpad gesture to go back or forward produces a natural transition where the whole web page slides off screen to reveal the target web page, but it backfires and feels artificial on single-page apps where the original navigation event from the target page was no so heavyweight as that.

The ecosystem effects of Safari include Handoff, which allows a browsing session to smoothly transition from device to device.

#### iCloud Keychain

While moving to Safari, I replaced 1Password with iCloud Keychain. It serves basic needs, but that’s it. It can store a username, a password, a 2FA code, and a domain name for each entry; nothing more.

It cannot store two domains for one entry, e.g. gmail.com and google.com. It cannot store arbitrary notes on an entry, e.g. the PIN that T-Mobile customer service agents ask for. It cannot name an entry, e.g. Washington Corporations and Charities System instead of ccfs.sos.wa.gov. It cannot store non-login entries like documents, ID numbers, or insurance information.

Using the 2FA field involved manual effort. Scanning almost never worked, so I would enter the 2FA secret by copying and pasting the code. Some sites provide a raw code as the code, while others wrap the code in a URL that contains other metadata. 1Password accepts either, but iCloud Keychain assumes you always hand it a code; if you hand it a URL it will silently accept it but produce incorrect codes.

I moved back to 1Password.

### Sheets → Numbers,

As a casual spreadsheet user, Sheets and Numbers are nearly identical. Numbers has nicer UX when editing formulae that visualizes any cell(s) being referenced. It also tries to humanize references, e.g. “Ben age” for a cell in a row with header “Ben” and a column with header “Age”, instead of something like A:123. This is nice until headers get long and multiworded. Overall it’s a wash.

### Docs → Pages

(To be filled in; I have not had much Pages experience.)

### Keep → Notes

I value simplicity and elasticity in notetaking—get out of my way and let me write, then let me deal with it later. Keep supplies that. Its layout is hard to browse, but it makes up for it with great search.

Notes is simple and elastic in a different way. Where Keep focuses on shortform sticky-style notes, Notes focuses on longform, with roughly the same text formatting options as Markdown (unfortunately without the markup). Instead of adding several notes to a category or color in Keep, I append to an existing note that contains several thoughts. This keeps the number of notes down, which makes categorization more reasonable, which improves browsing.

For me, it’s a wash between the two.

The ecosystem effects of Notes include Shortcuts and cross-app drag-and-drop. I use a shortcut to create a new note titled and sorted correctly before starting a Chinese lesson.

### Tasks / Google Assistant Reminders → Reminders

Reminders is one of the best-designed apps on iOS. Reminders can be scheduled to “pop” at a date, a date and time, a location, and/or when messaging `$PERSON`. Reminders can belong to lists (e.g. work vs. personal), lists can be shared (e.g. family chores), and reminders within shared lists can be assigned to people.

Reminders have a name, description, URL field, and priority (higher priority reminders are sorted higher and given special UI treatment). Reminders can have images attached to them, and any number of subtasks. They can be tagged and flagged.

It’s a powerful app, but everything is presented simply. There is a native macOS app that synchronizes, so I get proper notifications on most of my devices. In app, I use the “Today” view which shows reminders ready to be addressed.

As one use case, my spouse and I share a family reminders list. On that list, a reminder to take out the trash pops every trash day when I arrive home. It’s assigned to me, but if she happens to do it before I get home she can check it off; it won’t then pop for me.

The ecosystem benefits of Reminders include integration with the share sheet in native apps; e.g. sharing from Safari automatically fills in the URL field.

### Google Assistant → Siri

Siri is nearly strictly worse than Google Assistant. It can't answer questions like "What temperature do I need to cook chicken to?" or "Who played Alan in Tron Legacy?".

The ecosystem benefits of Siri include surface-level interaction with native apps: setting reminders, playing music, 

### YouTube Music → Apple Music

Apple Music is a great example of the big place Apple still struggles: services.

Importing my music was a headache and missed or incorrectly identified a lot of songs. The most reliable method I had involved an Automator workflow that would move the cursor to my existing library, select and copy a title, then move it to Apple Music’s search field, paste, and add the top result.

The UI of Apple Music for macOS is a let-down and is more on par with the iCloud Mail web interface than with any other native Apple app. Navigation is slow and unresponsive, which is compounded by the fact that it takes too many transitions to get where you’re going.

I miss Google Play Music.

The ecosystem benefits of Apple Music include tighter Shazam integration and better Siri support for when I’m driving and want to play something.

### Google Assistant Routines / Tasker → Shortcuts

I can’t say enough good things about Shortcuts. It is the most power-user-friendly thing about iOS and goes against all expectations I had around Apple and power users. I prefer it to shell scripts.

As one example, I like to write in an app called iA Writer. I use Shortcuts to automatically commit and push my writing to this website daily. Because iA Writer stores files in iCloud, and another app called Working Copy can interact with Git repositories, Shortcuts lets me glue them together:

1. (Working Copy) Pull from `twos.dev` remote
2. (iOS) Get contents of folder `iCloud/iA Writer/Published`
3. (Working Copy) Write CONTENTS OF FOLDER to `./src` in `twos.dev`
4. (Working Copy) Stage `./src` in `twos.dev`
5. Commit `twos.dev` with message “Automatic commit by iA Writer sync job”
6. Push `twos.dev` to remote

I’m a software engineer and am comfortable coding, but the fact that I could do all this without any was impressive. It’s also fun to say that [my phone is a vital part of my CI/CD pipeline][meta.html].

## Conclusion

This is a living document.

Its goal is to document whether the Apple ecosystem is bigger than the sum of its parts. It is. But the more interesting detail I’ve learned is that it’s the long tail of ecosystem benefits that makes up most of that excess. Not the two or three things per product I’ve mentioned above, but the dozens that happen without me noticing that add up to make a more enjoyable experience.

I equate it to working in a clean space vs. a messy space. There are functional benefits to working in a clean space—it’s easier to find things, spilling a liquid is not as destructive, you breathe in less dust—but the bulk of the benefit is in the hard-to-describe ways the space feels better and motivates more.

As of 2022 May, I’m overall happy with the Apple ecosystem and would not count it out from becoming my new preference.