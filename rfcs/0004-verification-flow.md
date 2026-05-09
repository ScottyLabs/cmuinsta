CMU Insta RFC 4
# Verification Flow
 - Author: @etashj
 - Created: 2026-05-09
 - Updated: 2026-05-09

## Overview

## Motivations
We need to ensure students who post content are verified to be committed and can only modify their own content. As a result, we will be using andrew ID as a verfiication method since it is delivered within a few days of the deposit fee for committment. ScottyLabs provides an authentication flow which we will utilize. 

## Goals
 - Define the authentication flow to sign in with andrew ID

## Non Goals
 - Implement a custom verification flow

## Detailed Design 
We will use the authentication flow provided by ScottyLabs to sign in with andrew ID. This is a Keycloak OIDC proxy to CMU's identity provider which uses SAML. This OIDC flow will return a JWT which will be used for API access as well as identity verification.

Additionally, we must verify that a user's instagram account is linked to their andrew ID. Since Instagram requires some interaction from the user, we will require the user to DM the account to create the initial point of contact. From there they will get a unique link which can be used to link their andrew ID to their instagram account. This will be stored in the database separately as a unique verifier. This is addressed in RFC 3. 

## Open Questions
 - Will there be workarounds? 
 - Are we doing too much? 

## Implementation Phases
 - The authentication flow will likely be implemented with the backend, and then added to the frontend.
