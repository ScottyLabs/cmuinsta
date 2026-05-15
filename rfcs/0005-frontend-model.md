CMU Insta RFC 5
# Frontend Model
 - Author: @etashj
 - Created: 2026-05-09
 - Updated: 2026-05-09

## Overview
This RFC will layout the model for the frontend of CMU Insta.

## Motivations
Making an intutive, easy-to-use frontend is integral to using this app. Thus we will define some design parameters and a general layout

## Goals
 - Define the frontend model for CMU Insta
 - Provide a clear, intuitive layout for the frontend
 - Ensure mobile compatibility
 - Offload computation to the frontend and minimize server load

## Non Goals
 - Layout a full design for the frontend

## Detailed Design
### Technologies
 - We will be using Svelte due to its developer friendliness and maturity
 - Typescript will be used for type safety
 - Svelte also provides us a way to compile to JS and create a static site
 - React was cosnidered but abandoned due to compelxity and other frameworks were rather immature

 ### Functionalities
 - The user must be able to register for an account with their andrew ID
 - The user must be able to link their instagram account in the flow described in RFC 4
 - The experience should occur largely on instagram with the exception of verification

## Open Questions
 - Should content filtering be automated or manual? Frontend or backend?
 - How can we ensure a smooth experience and support? 

## Implementation Phases
 - The frontend will be implemented after the databse and backedn since it largely must jsut interact with them
