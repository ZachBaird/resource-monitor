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

## Notes

- If an ApiKey is required for the endpoint, add its secret for the app's config with the `Name` of the app prefixing it (ex. APP1_SECRET for `app1`).
- If an ApiKey is required, make sure it's hashed and that your SECRET is added in your `.env`.
- When using an ApiKey, use the `Header` to determine what the HTTP header should be called. The default is `Authorization`, but some enterprises may want to use a custom header name (ex. `api-key`).
- If the ApiKey is not hashed, the SECRET is not necessary, but I do not recommend this approach.
- Only use GET endpoints for testing. A ping test should not mutate any data.
- The only exception to the above is if your RESTful API has POST endpoints that receive a body in the request for read purposes.