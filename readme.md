# Resource Monitor

This application pings various applications to check uptime.

## Description

In some enterprise environments, an application may be deployed across several servers behind a load balancer. 

Consider the following example: ApplicationA implements a feature-critical search endpoint. ApplicationA is deployed across 4 servers:

- devappserver01
- devappserver02
- devappserver03
- devappserver04

ApplicationA may have health check pings, but these often determine only if the application started without errors. Should an error arise during the search endpoint (ex. 3rd party outages, db connection issues, etc.), the health check ping would not reveal this.

This tool provides a way to test every server instance of an app with a `testUrl`. Results are displayed on a webpage.

See the `apps.json` as an example.