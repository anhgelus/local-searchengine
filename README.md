# Local SearchEngine

Powerful customizable local search engine based on [Grafisearch](https://github.com/Grafikart/grafisearch).

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

### Docker

If you want to try the searchengine or to use it in docker, you can!

#### Setup the environment

Clone the repository and enter in it.
```bash
$ git clone https://github.com/anhgelus/local-searchengine
$ cd local-searchengine 
```

Copy the `example/config.toml` at the root of the project.
```bash
$ cp example/config.toml config.toml 
```

Create a new csv file for the stats at the root of the project.
```bash
$ touch stats.csv 
```

#### Creating the compose file

Create a new compose file.
```bash
$ touch compose.yaml 
```

Put this information in the compose file.
```yaml
version: '3'

services:
  se: # se is for searchengine
    build: . # Build the Dockerfile
    ports:
      - "8042:8042" # The port to connect to
    volumes:
      - "./config.toml:/root/.config/local-searchengine/config.toml" # Path to the config.toml
      - "./stats.csv:/root/.local/share/local-searchengine/stats.csv" # Path to the stats file
      # And put here a path to the local wallpaper 
```

#### Building and Starting

Run and build the image.
```bash
$ docker compose up -d --build
```

Stop the container if you want to stop using this.
```bash
$ docker compose down
```

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

##### Where the is stored?

The data is stored in `~/.local/share/local-searchengine/`

## How does it work?

The goal is to create a binary starting a webserver and the search engine. 
This server will take the results from some search engines and will reformat this results.

## Technologies

- Go 1.20
- PuerkitoBio/goquery v1.8.0
- fatih/color v1.14.1
- pelletier/go-toml/v2 v2.0.6

