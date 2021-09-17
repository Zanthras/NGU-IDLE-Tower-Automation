# NGU Idle ITOPOD AP Idler

This program tries to take advantage of the odd ITOPOD quirk in regards to handing out AP/EXP gains. Per the 
[wiki note](https://ngu-idle.fandom.com/wiki/Arbitrary_Points#cite_note-3) you can vastly increase your AP gains by 
switching levels every time a different tier of the tower is about to produce rewards. Its a global kill counter but
with the first 20 tiers all having different kill requirements you can actually take advantage of all the tiers at the 
same time. 

## Requirements

- Go1.16 (probably will work with earlier just completely untested)
- tesseract (must be directly callable) https://tesseract-ocr.github.io/tessdoc/Installation.html#windows
- NGU Idle (steam version)
- GCC and stuff https://github.com/go-vgo/robotgo#requirements


## Basic operation

Start NGU Idle and be sure its the only window open called "NGU Idle" the save game folder is also named that, be sure
its not open in explorer. Then change the resolution setting to be 960*600 everything will fail at any other resolution.
Before any run update your PP/EXP bonus percentage to have accurate stats. You can find the values in ITOPOD.go line 
41/42. If the stats breakdown lists your total EXP bonus as 5560.69 put 55.6069 in the code. Recompile and you are ready
 to run the program. Position a cmd window in such a way that its not overlapping the ngu idle window and run the exe. 
When you are ready to stop the program hit the backtick key (next to the number 1) That will cause the program to stop, 
if that has issues for any reason locking your computer screen will also cause it to crash which conveniently also stops 
it. If you wish to temporarily pause the loop, hit the p key, to resume hit the p key again.
 
## Automation overview

The program is perfectly happy to run for arbitrary lengths of time, however if left running for more than 12 hours some
 convience functions will kick in. At 12 hours it will stop running the tower and eat an iron pill, followed by eating
 all maxed ygg fruit. Then another 12 hours of tower runs will commence. Once that completes a second pill will be eaten
 and then all ygg fruit (not just maxed). Finally it will spin the wheel and start the whole process over again. All
 tower runs start with the somewhat lengthy process of parsing the kill data out of the tool tips per tier, as this uses
 OCR via tesseract it has a non 0 chance of failing to parse correctly. Currently a rather simple workaround has been 
added which will on any parse failure trigger a single kill and restart the parsing process. This seems to work well 
enough for now.

Tower level selection is slightly complex but attempts to ensure the maximal pp/exp/ap gain. After every single kill the
algorithm will check if the closest tier is within 2 or 3 kills of generating rewards. If a tier is within that range it
 will move to the larger tier of the two. Level selection within a tier is always (tier*50)-5 or optimal level if its 
 the last tier. If there are multiple tiers 1 kill away from producing awards it will always pick the largest tier. 
Finally if no tiers are close it will switch to the max tier in order to have the highest PP farm potential. The reason
for not always going back to the max tier and instead having the kills +2/3 or logic is because tier switching is super
costly in terms of time, you want to be killing as fast as possible for maximal gains. Switching back to the max tier 
for a single kill before switching to the next target wastes a ton of time, this will happen more and more often as the
number of tiers you have access to increases. As all the relevent stats are constantly being printed out feel free to 
try things and see what the metrics print out.

## Metrics
During operation the program will constantly print out some internal tracking metrics. Example metrics
```
Hours: 0.68 Kills/KPM: 2268/55.38 EXP/EPM: 2847014/69521.60 AP/APM/KPA: 716/17.48/3.17 PP/PPPH: 242.5/355.37
```
- Hours: Tower runtime. As mentioned earlier the tower resets every 12 hours so this value will be larger than that
- Kills/KPM: Total kills and kills per minute. kpm is the main driving factor for getting more everything.
- EXP/EPM: Total exp gained and exp per minute. Be sure to set the exp bonus from total stats to be correct
- AP/APM/KPA: Total AP gained, ap per minute and kills per AP. When comparing two runs KPA measures how often you have tier transitions
- PP/PPPH: Total Perk Points gained, and perk points per hour gain. set the PP_BASE to be correct for your difficulty.
- FPS/Instant: Average framerate as measured by the duration it takes to get a pixel color. And the instant measurement of the same

The log file "itopod_rewards.log" will be generated and written to once per tower run (every 12 hours or when you hit `)
it contains the last status line printed on the console for long term tracking.

## Limitations
Click delay is vital to getting the application to work correctly. Testing on my single machine shows 100% of clicks
registering with a 25 ms delay after 5000 clicks tested. 20ms delay had only a 99% success rate. Other machines might 
have other success rates. Missing clicks is bad, sometimes I can account for it, sometimes I cant. 
 
This program is 100% based on reading pixel colors and controlling your mouse, using your computer while this is running
 does not work, turning your monitor off while this is running does not work. Locking your computer does not work. 

Ive only tested this with windows 7, probably will work with win10. 100% will not work with linux without modifications

## Disclaimer

Use at your own risk, it shouldnt do anything bad, and the code is hopefully super easy to read to validate that.