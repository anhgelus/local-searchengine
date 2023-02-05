# Local SearchEngine

Powerful customizable local search engine based on [Grafisearch](https://github.com/Grafikart/Grafisearch).

## Features

- Nice wallpaper background (like Bing)
- Fast
- Support bangs (!rt, !wrenfr...)
- Customizable theme
- Instant response in some cases (like timer, calculatrice)
- Have information in the full screen
- Block some trash websites or "pay to view" website (pinterest, allocine, jeuxvideos.com...)
- No ads

## Installation

### Linux

Download the [latest binary](https://github.com/anhgelus/local-searchengine/releases/latest)

Or you can build it yourself

#### Build the binary

Clone the repository
```bash
$ git clone https://github.com/anhgelus/local-searchengine 
```

Build the repository with Go 1.20
```bash
$ go build -o local-searchengine .
```
#### Launch the binary for the first time

Run the binary and use the argument "install"
```bash
$ ./local-searchengine install
```

Exit (ctrl+c)
Start the service
```bash
$ systemctl start --user local-searchengine.service
```

Configure your browser with this ip address : `https://localhost:8042/?q=`

#### Configuration

Open the configuration file `~/.config/local-searchengine/config.toml`
```toml
Version = '0.1'
AppName = 'Local SearchEngine'
BlockList = ['pinterest.com', 'allocine.com', 'jeuxvideo.com', 'lemonde.fr', 'w3schools.com', 'pinterest.fr']
WallpaperPath = ''
LogoPath = ''
```

And customize what do you want.

## How does it work?

The goal is to create a binary starting a webserver and the search engine. 
This server will take the results from some search engines and will reformat this results.

## Technologies

- Go 1.20
- PuerkitoBio/goquery v1.8.0
- beevik/etree v1.1.0
- fatih/color v1.14.1

