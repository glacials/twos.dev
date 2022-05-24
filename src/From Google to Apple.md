---

date: 2022-05
---

filename: apple.html
created: 2022-05
updated: 2022-05
---

TODO: USB-C (maybe?), Touch Bar, SMS in macOS, 

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

My biggest fear with this change was losing a decade of Google Nest ecosystem benefits, but running Homemanager on a Raspberry Pi made it all work seamlessly, even down to our Nest doorbell camera showing on our Apple TV picture-in-picture when the doorbell rings. (If you’re not keen to se up a raspi, try Starling Home Hub instead.)

Ecosystem benefits:

- Can use Shortcuts for automation

### Starling Home Hub
- When my doorbell rings, my Nest camera shows up in picture-in-picture

### Apple TV (service)

Replacements afforded:

- None

Thoughts:

Apple TV feels like early Netflix productions—everything is high quality (both technically and narratively). Ted Lasso is the most wholesome show of the last decade. See and WeCrashed are shows I’m working through but have been good—about 70-80% as good as Breaking Bad, which is my barometer these days.

Ecosystem benefits:

- None

### AirPods Pro 

Replaces Google Pixel Buds 2 and Bose QC35 II.

Thoughts:

Everything you’ve heard about the AirPods Pro noise cancellation is true. More than once I’ve forgotten they’re in my ears; transparency mode is good.

Ecosystem benefits:

- Automatically switches between MacBook, iPad, and iPhone based on heuristics like screen-on or lid-closed
- Supports Siri, but I don’t use it

### iMac

Replacements afforded:

- Google Code cron → macOS’s Shortcuts.app
- Gmail → iCloud Mail (see below)

Thoughts:

This is a used iMac from 2017; no fancy M1 chip. 

Ecosystem benefits:

- One more device AirPods Pro can switch between
- 

### Mail

Replaces Gmail UI.

Thoughts:

Mail's junk filter is bad. It has sorted roughly one legitimate email per week into the junk folder, even after multiple months of correcting it. Sometimes two similar emails from the same mailing list will come in around the same time, one into junk, one not.

Mail's search feature is bad. On macOS, it will assume you want to search the current mailbox, e.g. Inbox. Because I archive emails that don't need my attention, this is never what I want to search. (Mail on iOS and iPadOS correctly searches everything by default.) When you are done searching and click Inbox to "return home", the search field doesn't clear on its own, leading you to believe your Inbox is empty.

Even beyond that, Gmail unsurprisingly wins search by an order of magnitude. As one example, Gmail searches the contents of PDFs attached to emails.

Mail is okay at displaying emails as conversations. Once in a while, it misses something. When you add an email to a thread and click send, it doesn't immediately appear below the thread, causing you to believe something went wrong. Presumably it's not discovered until the screen re-scans the Sent mailbox for related emails.

Mail is bad at dealing with calendar invitations. They appear as normal attachments -- two per invitation, an ICS file and one generically called "Mail Attachment". You have to learn by trial and error that the ICS file is not the one you should use, although it appears to work, because it holds no ties back to the creator. Instead, you should open "Mail Attachment", which will allow you to RSVP.

Ecosystem benefits:

- Share menu feels more natively “part of” iOS, iPadOS
- Gesture support on macOS 
- Drag-and-collect support on iOS

Notes:

- If you use Gmail and want the UX of Mail.app, I suggest Mimestream. Created by the ex-Apple developer of Mail.app to use Gmail APIs and features instead of SMTP.


### iCloud Mail

Replaces Gmail backend.

Thoughts:

After seeing many harrowing stories of Google users' accounts being inexplicably banned I decided I'd start the years-long process of changing my email address, to a domain I control. Luckily, Apple released custom domain support for iCloud Mail around the same time.

The iCloud Mail web app experience is bad. With Apple's quality bar they set for themselves, I'm surprised it exists. I quickly switched to Thunderbird on my gaming computer.

Although the macOS Mail app supports auto sorting rules, these apply locally. The iCloud Mail service has a separate rules system with the same effect, but the two have disparate feature sets and don't synchronize. For example, iCloud rules can match one condition per rule -- from: or subject:, but not both.

Later, I inherited an old iMac and I have it always running so my mail always gets sorted.

Ecosystem benefits:

- Push works for Mail when used with iCloud Mail

### Calendar

Replaces Google Calendar.

Thoughts:

Before transitioning my custom-domain email address to iCloud Mail, it was serviced by Google Workspace; I had this Google Workspace account added to my iPhone and iPad. This got me into a weird state where Calendar would not let me add an event to my iCloud calendar if the invitee was my new email address, if that Google Workspace account was present on my device, but it would allow me to add it to that Google account's calendar. When I removed that Google account from my device, I could add the event to my iCloud calendar fine.

I don't recommend using the Week view in Calendar. Similar to Google Calendar, Calendar uses a horizontal red bar to represent the current time; however this red bar extends to occupy 100% of the width of the week, and does not do a good enough job of showing you which day today is. This has led me to mistake the wrong day of the week for today multiple times, inducing panic about being late for meetings and reaching out to people to reschedule other meetings. I now use Day view, where no red bar shows unless viewing today.

Ecosystem benefits:

- Quick look is a blast, the same way it is for previewing files in Finder

### Maps

Replaces Google Maps.

### Safari

Replaces Google Chrome.

### Numbers, Pages

Replaces Google Sheets, Google Docs.

Thoughts:

Ecosystem benefits:

- Siri suggestions take me to specific sheets, docs

### Notes

Replaces Google Keep.

Thoughts:

Ecosystem benefits:

- I use Shortcuts to start my Chinese lesson, which creates a scratch note with the current date as its title and opens my Zoom meeting

### Reminders

Replaces Google Tasks and the reminders feature of Google Assistant.

Thoughts:

Ecosystem benefits:

- Siri

### Siri

Replaces Google Assistant.

Thoughts:

Siri is nearly strictly worse than Google Assistant. It can't answer questions like "What temperature do I need to cook chicken to?" or "Who played Alan in Tron Legacy?".

Otherwise, the user experience of Siri has the classic Apple pleasantness.

Ecosystem benefits:


### Shortcuts

Replaces Google Assistant routines and apps like Tasker.

Thoughts:

I can’t say enough good things about Shortcuts. It is the most power user friendly thing about iOS and goes against all expectations I had around Apple not allowing a hacking mentality on iOS. I prefer it to shell scripts on a cron job.

As one example, I like to write in an app called iA Writer, and I use Shortcuts to automatically commit and push those to this website daily. Because iA Writer stores files in iCloud, and another app called Working Copy can interact with Git repositories, Shortcuts lets me glue them together:

1. (Working Copy) Pull from TWOS.DEV remote
2. (iOS) Get contents of folder ICLOUD/IA WRITING/PUBLISHED
3. (Working Copy) Write CONTENTS OF FOLDER to ./SRC in TWOS.DEV
4. (Working Copy) Stage ./SRC in TWOS.DEV
5. Commit TWOS.DEV with message AUTOMATIC COMMIT BY IA WRITER SYNC JOB
6. Push TWOS.DEV to remote

I’m a software engineer and am comfortable coding, but the fact that I could do all this without any was impressive.