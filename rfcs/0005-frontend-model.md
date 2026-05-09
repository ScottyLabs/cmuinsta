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
 - The user should be able to preview their post almost exactly as it would appear on Instagram
 - Images should be edited (into the cover template) and downscaled/cropped for the instagram format on the frontend
 - The user should be able to save their progress automatically (i.e. uploads go straight to server immediately) and resume when coming back
 - The user should be able to post and view their status on the queue when they have already posted
 - They should not be able to post multiple times or modify their post after it has been posted
 - They should be able to cacnel their post if possible given the constaints of RFC 2
 - There should be a filter on posts to ensure no explcit content is made through this platform (that is, text or graphics, frontend or backend)
 - All this should be fully mobile compatible to ensure smoothness of the experience

 ### Implementation
  - Likely this will be one dashboard in an SPA
  - One form, one status page, registration page

## Open Questions
 - Should content filtering be automated or manual? Frontend or backend?
 - How can we ensure a smooth experience and support? 

## Implementation Phases
 - The frontend will be implemented after the databse and backedn since it largely must jsut interact with them