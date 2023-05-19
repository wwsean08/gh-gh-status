# GitHub Status Checker

This `gh` extension can be used to monitor the current status of GitHub based on the status page.  It can either be run once, or run with the watch flag to poll the page at a 1 minute interval.

## Installation
```shell
gh extension install wwsean08/gh-gh-status
```

## Upgrade
As a note, GitHub really doesn't like your command starting with `gh` like this one does.  In order to upgrade the plugin you can run one of the following commands:

```shell
gh extension upgrade wwsean08/gh-gh-status
```

or

```shell
gh extension upgrade gh-gh-status
```

## Usage
### Run once
```shell
gh gh-status
```
![](docs/img/run-once.png)
### Poll constantly
```shell
gh gh-status --watch
```
![](docs/img/watch.png)
