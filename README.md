# Browser tab groups

Save and open **links** from the command line with ease. A tap group is a collection of links (urls) that belong together.

Example: `work` tap group would contain links for work. `uni` tap group would contain links for uni etc.

## Installation

1. Manually by downloading the binary from the release page
1. Using go with `go install github.com/magdyamr542/browser-tab-groups@latest  `


## Features

1. Group `urls` with a label
   - use case: for work i want to quickly open `[Gitlab, Jira, Github, ...]`
   - use case: for a given issue i want to quickly open its `[jira link, bitbucket pr, bitbucket branch, ...]`
   - use case: for uni i want to quickly open `[Moodle, Web mailer, ...]`
1. Open a group of `urls` from the cli in the browser
1. Open a single `url` from a group of `urls`
   - use case: for a given issue i saved its urls `[Bitbucket, Jira, Github, ...]` but want to quickly open only its `Jira link` without the rest of urls because i don't need them right now.
1. Remove a group of `urls`
   - use case: after being done with a ticket. i want to remove all of its saved links

## Usage

1. `br` will print the usage
1. `br list` to list all saved tab groups
1. `br add <tap group> <url>` to add the `url` to the tap group `tap group`
1. `br open <tap group>` to open all `urls` in the tap group `tap group` in the browser
1. `br open <tap group> <url matching string>` to open the url(s) that _fuzzy match_ `url matching string` in the browser

## Workflow looks like this

1.  `br add express-routing https://github.com/expressjs/express`
1.  `br add express-routing https://expressjs.com/en/guide/routing.html`
1.  `br ls`

    ```
    uni:
    https://webmail.tu-dortmund.de/roundcubemail/

    express-routing:
    https://github.com/expressjs/express
    https://expressjs.com/en/guide/routing.html
    ```

1.  `br open express` would open the two links under the `express-routing` group in the browser
1.  `br open express git` would open the link for **express github** because it uses `fuzzy finding` to filter for links based on the user's input
