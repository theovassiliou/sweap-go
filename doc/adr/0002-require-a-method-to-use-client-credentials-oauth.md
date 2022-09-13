# 2. Require a method to use client_credentials OAuth

Date: 2022-07-27

## Status

Accepted

## Context

The Sweap API uses OAuth2 client_credentials method to authenticate access to the API. So a "make or buy"-decision is required on how to implement the token flow.

While it is straightforward to obtain the required token (POST request with credentials embedded as form data), the received token has to be used as Bearer access token with every request. In addition regular refresh of the bearer token is required.

A quick research identifies golang.org/x/oauth2/clientcredentials package as suitable for integration.

Light source code review showed no particular issues.

## Decision

We are using golang.org/x/oauth2/clientcredentials package as OAuth2 client credentials package, and avoid an own implementation of the token flow.

## Consequences

### Pro

Proven access to OAuth2 credentials.

### Cons

Have to monitor for possible security issues of the package.
