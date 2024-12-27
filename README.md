<h1 align="center">
  pagode
  <br>
</h1>

<h4 align="center">Fast passive google dorking enumeration tool.</h4>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Install</a> •
  <a href="#running-pagode">Usage</a> •
  <a href="#post-installation-instructions">API Setup</a> •
</p>

---

`pagode` is a google dorking enumeration tool that returns valid links for websites, using passive queries to google search API. It has a simple, modular architecture and is optimized for speed. `pagode` is built for
doing one thing only - google dorking, and it does that very well.

# Features

- Fast and powerful google dorking tool
- Multiple output formats supported (file, stdout)
- Optimized for speed and **lightweight** on resources
- **STDIN/OUT** support enables easy integration into workflows

# Usage

```sh
pagode -h
```

This will display help for the tool. Here are all the switches it supports.

```yaml
Usage: ./pagode [flags]

Flags:
INPUT: -g, -dork string[]  dork to search
  -gL, -list string   file containing a list of dorks to search
  -d, -domain string  domain to search with "site:" operator

OUTPUT: -o, -output string  file to write output to

CONFIGURATION:
  -c, -config string  config file for API keys (default "config.yaml")
  -p, -proxy string   HTTP(s)/SOCKS5 proxy to use with pagode
  -m, -max int        maximum number of results per search
  -t, -threads int    number of concurrent goroutines for querying (default 10)

DEBUG: -silent         show only results in output
  -version        show version of pagode
  -v              show verbose output
  -nc, -no-color  disable color in output
```

# Installation

`pagode` requires **go1.20** to install successfully. Run the following command to install the latest version:

```sh
go install -v github.com/ergors/pagode/cmd/pagode@latest
```

## Post Installation Instructions

For using pagode you need the API keys and Search ID from Google Search API, to obtain them [follow this guide](https://developers.google.com/custom-search/v1/overview).

config.yaml

```yaml
google:
  - apikey:searchid
```

# Running Pagode

To run the tool with a dork, just use the following command.

```console
pagode -g ext:php

                                 __
    ____  ____ _____ _____  ____/ /__
   / __ \/ __ `/ __ `/ __ \/ __  / _ \
  / /_/ / /_/ / /_/ / /_/ / /_/ /  __/
 / .___/\__,_/\__, /\____/\__,_/\___/
/_/          /____/

		projectdiscovery.io

[INF] Loading config from /home/renato/.config/pagode/config.yaml
https://snskeyboard.com/instafont.php
https://www.lacentrale.fr/lacote_origine.php
https://www.194964.com/top.php
http://www.flugzeuginfo.net/table_airlinecodes_airline_en.php
https://www.marketo.com/privacy.php
https://lotto.auzonet.com/bingobingo.php
https://www.livesudoku.com/indexit.php
```

The links discovered can be piped to other tools too. For example, you can pipe the discovered links to [`httpx`](https://github.com/projectdiscovery/httpx) which will then find
valid HTTP links.

```console
echo "ext:jsf site:com" | pagode -silent | httpx -title

    __    __  __       _  __
   / /_  / /_/ /_____ | |/ /
  / __ \/ __/ __/ __ \|   /
 / / / / /_/ /_/ /_/ /   |
/_/ /_/\__/\__/ .___/_/|_|
             /_/

		projectdiscovery.io

[INF] Current httpx version v1.6.9 (latest)
[WRN] UI Dashboard is disabled, Use -dashboard option to enable
https://meteor.springer.com/login.jsf [Login - Meteor]
https://github.com/bisqwit/that_editor/blob/master/conf.jsf [that_editor/conf.jsf at master · bisqwit/that_editor · GitHub]
https://jaegermeister-learning.com/pages/login.jsf
https://secure.magnushealthportal.com/sslogin/login.jsf [Login to Magnus Health | MyMagnus.com]
https://retirementplans.vanguard.com/VGApp/pe/pubnews/SpecialNeedsKids.jsf [Vanguard - Page not found]
https://annuities.talcottresolution.com/asc/ContactUs.jsf [Contact Us]
https://ppl.aiacompanystore.com/index.jsf [Home - PPL (PPL)]
https://idahocom.aiacompanystore.com/index.jsf [Home - Idaho College of Osteopathic Medicine (ICOM)]
https://colgatepalmolive.yet2.com/res/innovation-portal/submission-form.jsf [Colgate Palmolive Open Innovation Portal -- Submission Form]
https://arabicstorypedia.com/itemdetails.jsf?itemIDNew=3210 [404 Not Found]
```

<table>
<tr>
<td>

# License

`pagode` is made with 🖤 and based on `subfinder` by the [projectdiscovery](https://projectdiscovery.io) team.
