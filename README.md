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

## How does it work?

The goal is to create a binary starting a webserver and the search engine. 
This server will take the results from some search engines and will reformat this results.

## Technologies

- Go 1.20
- PuerkitoBio/goquery v1.8.0
- beevik/etree v1.1.0
- fatih/color v1.14.1

