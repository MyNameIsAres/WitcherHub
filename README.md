## WitcherHub

WitcherHub is a project that aims to create more interactivity between streamers on Twitch and viewers. Viewers can spawn in the game as a companion for the player, play Gwent with (or for) them, create challenges, and even quests.

WitcherHub is a combination of several components such as:

### WitcherConnect
WitcherConnect is responsible for connecting the The Witcher 3 to the outside world and allowing requests to flow from different sources to the game. Heavily inspired by the approach currently used in WolvenKit, it has enabled WitcherHub to send exec functions to the game.

### A GUI desktop (in progress - now showcased... yet)
An easy to use GUI desktop tool created with Golang Wails to setup and manage the WitcherHub bot and Twitch extension.

### A combination of different (with permission modified) The Witcher 3 mods to create the interactivity in the game. Such as:
    * Multiplayer Gwent (extension)
    * Create challenges / quests through the Challenge and Quest creator, featured in the WitcherHub extension
    * Spawn monsters
    * Spawn the viewer as a companion
    