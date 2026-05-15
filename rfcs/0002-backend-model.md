CMU Insta RFC 2
# Backend Model
 - Author: @etashj
 - Created: 2026-05-07
 - Updated: 2026-05-15

## Overview
This RFC will outline the backend model for CMU Insta, including endpoints, possible authentication mechanisms, and some implementation details. 

## Motivations
Defining the endpoints here will serve as a point of documentation for developers: 
 - Allowing for an easy reference when implementing
 - For future developers, understanding the design choices made and how endpoints interact
 - Matching development patterns and adapting to future changes. 

## Goals
 - Define a method of authentication for people who make calls to the backend
 - Define the endpoints available for the Instagram DM hooks to call
 - Define the data models used for requests and responses
 - Define the Go libraries that are planned to be used. 

 ## Non Goals
 - Define specific implementation procedures
 - Lay out the exact database schema

## Detailed Design
Implementing the backend involves implementing a secure RESTful API which interacts with the Instagram API for DMs and posting plus the database for tracking. The majority of the user interaction should occur within the instagram DM with the sole exception of andrew ID verification. 

### Endpoint Authentication
CMU's authentication flow uses SAML and ScottyLabs includes a Keycloak instance with SAML to OIDC. Thus we will be able to use some token-based authentication methods, likely JWT tokens which will be covered in RFC 4. 

### Endpoints
All endpoints will be available at `/api/v1/`, but we will only expose one endpoint for the hook to fire at on a user DM. For example the `hook` endpoint will be at `/api/v1/hook/`.

#### Instagram Hooks
```
POST /hook
```
**Description:** The endpoint for the Instagram API to fire when a linked account receives a DM. 

**Request Body:**
 - Content-Type: `application/json`
 - Body: Follows the Instagram specification
 ```json
 {
   "object": "instagram", 
   "entry": [{
     "id": "BOT_IGSID"
     "time": 1748000000, 
     "messaging": [
       {
         "sender":     { "id": "USER_IGSID"    },
         "recipient":  { "id": "BOT_IGSID"     },
         "timestamp":  1748000000, 
         "message":    { /* Details omitted */ }
       }
     ]
   }]
 }
 ```
The `message` body can have many shapes based on the message context. We consider plain text and image formats. 
*Plain text*
```json
"message": {
  "mid": "aGVsbG8gd29ybGQK", 
  "text": "hello"
}
```
*Image(s)*
```json
"message": {
  "mid": "aGVsbG8gd29ybGQK", 
  "attachments": [{
    "type": "image", 
    "payload": {
      "url": "https://cdn.instagram.com/..."
    }
  }]
}
```
*Quick Reply*
```json
"message": {
  "mid": "aGVsbG8gd29ybGQK", 
  "text": "Edit",
  "quick_reply": {
    "payload": "EDIT"
  }
}
```



**Constraints:**
 - We will ignore other forms of input as we only consider when a message sis received. 
 - This is the only exposed endpoint for the instagram API

**Response:**
 - Status: `200 OK`
 - Body: The code stops instagram from refiring the hook, error handling will be handled separately
 - Note returning a `4xx` code will cause instagram to treat as a permanent failure and `5xx` will be temporary and trigger a retry

#### OIDC Authentication
```
POST /auth/callback
```
**Description:** Converts a code to a JWT using the Keycloack proxy (see RFC 4)
 
**Request Body:**
- Header: None
- Body: 
```json
{
  "code": "provided code"
}
```
**Response Format:**
- Status: `200 OK`
- Body:
```json
{
  "token": "JWT for authentication",
  "expires_in": "5000",
}
```

**Response Codes:**
| Code | Description |
|------|-------------|
| 200 OK | User data successfully retrieved; returns the full user object. |
| 400 Bad Request | Missing code in request body |
| 401 Unauthorized | Missing, expired, or invalid code/JWT. |
| 500 Internal Error | Database read failed or server-side exception. |

### Go Implementation Details
To implement the cmuinsta authentication flow using Gin, we will primarily rely on the framework's middleware pattern to enforce the security requirements of RFC 4. By integrating `golang-jwt/jwt/v5` for cryptographic parsing and `coreos/go-oidc/v3` for automated public key discovery from the ScottyLabs Keycloak instance, you can create a stateless "gatekeeper" that intercepts requests to your `/api/v1/` route group. This architecture ensures that sensitive endpoints, such as `/create`, never receive spoofed data; instead, they securely extract the verified preferred_username claim directly from the Gin context, allowing your business logic to remain clean, modular, and focused on core functionality. 

We will use a state machine where we keep track of a position in the process so that users can send indivdiaul messages for each step. This will be tracked in the database. 

### Queue Structure
The queue will be implemented likely using the same database system. An important consideration is buffering between posts to avoid overwhelming the instagram API whether intentionally in an attack or an accidental influx of users. 

If we consider about 2,000 students for the sake of simplicity, with 30% being ED and 70% non ED, we find about 600 students starting in late December and 1,400 starting in mid-March. This is of course an absolute maximum and the patterns will not be uniform. As a result, instead of attempting to evenly distribute posts, we will simply add a timer between users to maximize the number of posts per day being spread evenly.

That is, one student may post, but then the posting process is locked for a buffer period.

Via the API, we have up to 50 posts per 24 hours As a result we will place a 20 minute buffer between posts and running for 50 posts starting at 7:00 AM. This will run us to 11:40 PM with the maximum number of posts. 

The queue will be managed via the Database and described in RFC 3. 


## Open Questions
 - How will the database be structured for this?
 - How will authentication be integrated from the frontend?
 - Is there a more ideal queue structure? 
 
## Implementation Phases
 - This will be implemented in parts parallel to the database schema which is covered in RFC 3. 
 - Expect many small commits addressing each numbered endpoint.
