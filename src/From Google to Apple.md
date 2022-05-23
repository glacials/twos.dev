---
filename: apple.html
date: 2022-05
---

# From Google to Apple

I have been an Android user for 10 years -- Google Android, at that. Starting with the 2011 Galaxy Nexus, I've been a Nexus or Pixel user my entire adult life.

I bought an iPhone. I try to be an informed person, and that means knowing what I’m missing out on, to understand why (or whether) I prefer the thing I chose.

For all the hate the walled garden gets, I understand I can't appreciate it perched on the wall. In Apple land, one device on its own doesn’t make a convincing argument—the ecosystem cohesion does.

Knowing the garden may or may not let me out, I risk being trapped but knowing both sides, over letting fear of entrapment trap me into knowing only one. This is my story of hopping the fence.

## The Plan

My last Apple device, other than MacBooks used for development, was an iPhone 3G in 2010. Starting with iPad, I’ve tiptoed my way into the ecosystem, checking my satisfaction as I go. Each component I add theoretically adds value to each existing component. For example, having an iPad and a Mac enables using the iPad as an external display for the Mac. I will document here these benefits and my thoughts about each product I add.

## Execution 

### Precondition: MacBook

I have used a MacBook running macOS as my development machine since I switched from ArchLinux on my Cr-48 in 2013.

### iPad

The iPad is my slow ease into the iOS ecosystem. In the past, I’ve used a Nexus 7 as a household tablet. To me, a the point of a tablet is to replace all uses of my phone while at home. This gives an upgraded experience and drains a separate battery than the one I may need later.

However, within about a year of buying my Nexus 7 I’d stopped my habit of carrying it room to room. It afforded too little benefit and took too much effort to maintain separately from my phone. Apps didn’t expect one user to be using two devices; bugs were frequent. Almost zero apps used a different design for the tablet form factor.

I expected a similar experience from iPad.

Replacements afforded:

- None

Thoughts:

With iPad, I’ve carried it room to room for two years now. The app ecosystem better understands that tablets exist. For example, Octal, an iOS HN viewer, will show the front page and the current article side by side in independent scroll views. The multitasking experience is opinionated, but phenomenal. 

Ecosystem benefits:

- Sidecar (use iPad as external monitor for macOS)
- Universal Control (use mouse and keyboard from macOS to control iPad)

### iPhone

Replacements afforded:

- Hangouts/Chat → iMessage
- Google Photos → iOS Photos + iCloud
- Google Maps → Apple Maps
- YouTube Music → Apple Music
- Google Podcasts → Apple Podcasts
- Google Play Books → Apple Books

Thoughts:

Things immediately become clear in how Google and Apple approach customer problems. In the Google ecosystem, 

Ecosystem benefits:

- macOS’s Preview.app will let me use my iPhone as a cursor to annotate or sign a document, in two clicks.
- For some friends, SMS becomes iMessage
- For SMS and iMessage, I can interact from my MacBook

### Apple TV (hardware)

Replaces Chromecast.

Thoughts:

Given that we don’t have HomePods, and my spouse forbids me to get any after filling our home with more Google Home Minis than we have rooms (they give them out like candy to people deep in the ecosystem), we’ve had to get used to using a remote again.

It’s a high quality remote and feels good in the hand, but it’s a remote and it can be lost. We’ve 3D printed a holder for it that attaches to the coffee table.

Ecosystem benefits:

- When an onscreen keyboard shows up, a notification shows on my iPhone and iPad allowing me to type on that device instead. This lets me use 1Password autofill for passwords on those services like Hulu that haven’t yet added QR code authentication or similar.

### Apple Home

Replaces Google Home.

Thoughts:

My biggest fear was that a decade of building up Google Nest devices would be for naught, but Starling Home Hub and later Homemanager with a Nest plugin made it all work seamlessly.

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

Replaces Google Assistant routines.

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