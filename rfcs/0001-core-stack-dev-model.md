CMU Insta RFC 1
# Core Stack and Development Model
 - Author: @etashj
 - Created: 2026-05-07
 - Updated: 2026-05-07

## Overview
This RFC works to establish the general guidelines of the tech stack for CMU Insta. Specifically, this includes Svelte for the frontend including mobile-responsive design; a performant backend which relies primarlity on a RESTful API written with Go; and finally a relational database using SQLite. 

## Motivations
CMU Insta is a realtively simple project, which can prove to be usefulf or several years without much modification if planned properly. It needs to be maintainable and well designed to prevent the hassle of having to do a large refactor/reimplemntation.

## Goals
 - Define the backend, frontned, and database
 - Support development for the future

## Non Goals
 - Specific details will be covered in later RFCs

## Detailed Design

### Frontend
We will use Svelte with Vite since it requires low boilerplate and allows static site generation. Svelte is also more develpoer friendly than react. The style will be implemented with tailwind css due to the convenience of compoenents, inline styling, and utility classes. This will also allow mobile-responsiveness using compoenet based, reactive design. 

We will also need a form of verification, for whcih we will use andrew IDs when they are available. This is to be addressed in RFC 4. For students who submit before andrew IDs are released, we have to consider how we can verify admission, albeit manually. Additionally, we should verify the validity of instagram username submitted since they will be a collaborator or tagged on a post, this will also be addressed later, but is likely to be done via some form of user interaction on the external platform. 

### Backend
We will use Go to write the backend RESTful API due to its performance and ease of writing ocde for a small task. It's stadnard library should provide most of the HTTP features required for this task.  

The GO backend will be designed to handle user submissions, verifications if needed, interactions with the database, and scheduled posts. 

### Database
The database will have to host data about users and post data. This includes names, andrew ids, instagram usernames, post content (caption and photos), and internal status details. The database details will be covered in RFC 3, but will consist of some ocmbination of a SQLite relational database and file storage for images and captions. This is due to the simplcity and portability of SQLite and since our requests are limited to the size of the incoming class. This can reasonable be capped to 1,750 students and is almost guaranteed to be below 2,000 students, of whcih only a fraction will submit data. 

### Development Environment
We will be using Nix as a pacakge manager with direnv to ensure a consistent development environment across all contributors.

### Alternatives Considered
 - For frontend React was considered, but Svelte was chosen due to its low boilerplate and static site generation capabilities.
 - For backend, we considered Rust but avoided it due to complexity and chose against Typescript and Node due to speed
 - We also considered Docker for development but decided against it due to the pivot in hosting for ScottyLabs. 

## Open Questions
 - What will the authetnciation flow be? How can we verify admissions and instagram usernames?
 - What is the database structure? 
 - How can we implement the backend and secure endpoints to avoid false submissions? 

## Implementation Phases
1. Implementing the backend RESTful API
2. Implementing the databse
3. Implementing the frontend and mobile responsiveness
