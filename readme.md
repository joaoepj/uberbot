# uberbot

The bot that watches over the Ottomated discord server.

## Description

uberbot is a custom bot written in discordgo. It has a custom data handler, command handler, and argument parser. 
Some main features of uberbot are moderation and meme like functions

## Dependencies
* Go 1.16
* Custom [discordgo](https://github.com/bwmarrin/discordgo) [fork](https://github.com/qpixel/discordgo) <!--Try to switch back to the old one sooner than later-->
* OrderedMap
* GodotEnv
* [tinylog](https://github.com/ubergeek77/tinylog)

## Development

### Setup
Make sure you have an ssh key attached to github and have configured it with git. Please also have a gpg key that works with git and can sign commits.

1. Fork the repository
2. Clone your fork of the repository
```sh
git clone git@github.com:<username>/uberbot.git
cd uberbot
```
2. Update submodules
```sh
git submodule update --init --recursive
```
3. Resolve dependencies
```sh
go mod tidy
```
4. Build uberbot
```shell
go build 
```
5. Create a .env file and add these lines to it
```shell
UBERBOT_TOKEN=<discordtoken>
ADMIN_IDS=<yourdiscordid>
```
6. Run uberbot
```shell
./uberbot
```

## Contributing
Any contributions you make are greatly appreciated.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
