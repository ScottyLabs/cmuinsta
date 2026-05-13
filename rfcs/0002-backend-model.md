CMU Insta RFC 2
# Backend Model
 - Author: @etashj
 - Created: 2026-05-07
 - Updated: 2026-05-12

## Overview
This RFC will outline the backend model for CMU Insta, including endpoints, possible authentication mechanisms, and some implementation details. 

## Motivations
Defining the endpoints here will serve as a point of documentation for developers: 
 - Allowing for an easy reference when implementing
 - For future developers, understanding the design choices made and how endpoints interact
 - Matching development patterns and adapting to future changes. 

## Goals
 - Define a method of authentication for people who make calls to the backend
 - Define the endpoints available for the CMU Insta frontend to call
 - Define the data models used for requests and responses
 - Define the Go libraries that are planned to be used. 

 ## Non Goals
 - Define specific implementation procedures
 - Lay out the exact database schema

## Detailed Design
Implementing the backend involves implementing a secure RESTful API which interacts with the Instagram API for posting and the database for tracking. 

### Endpoint Authentication
CMU's authentication flow uses SAML and ScottyLabs includes a Keycloak instance with SAML to OIDC. Thus we will be able to use some token-based authentication methods, likely JWT tokens which will be covered in RFC 4. 

### Endpoints
All endpoints will be available at `/api/v1/`. For example the `create` endpoint will be at `/api/v1/create/`.

#### 2.1 Instantiating a User
```
POST /me
```
**Authentication:** Required (Bearer JWT)

**Description:** Initializes a new user with null fields in the database and filesystem

**Request Body:**
 - Content-Type: `application/json`
 - Body:
 ```json
 {
   "instagram_id": "unique identifier to link to an instagram account", 
   "name": "display name for the user", 
   "major": "integed major caption content", 
   "hometown" : "desired hometown"
 }
 ```
**Constraints:**
 - `instagram_id`: Must be a unique string identifier for the user's Instagram account that is accepted by the backend (this is defined in RFC 4 and mentioned in part in RFC 3). 
 - `hometown`: The location should always be "City, State" for U.S. residents (no "USA" afterward), and "City/Region, Country" for internationals.
 - `hometown`: No short hands for cities like NYC or countries like UAE and always two digit codes for states like PA

**Response:**
 - Status: `201 Created`
 - Body:
```json
{
  "andrewid": "etashjha",
  "name": "Etash", 
  "major": "Computer Science", 
  "hometown": "Pittsburgh, PA", 
  "instagram_username": "etashjha",
  "caption": "The user's caption capped at 2,200 characters",
  "image_count": 10
  "queued_at": "2026-05-12T23:48:32.257Z",
  "queue_position": 3
  "posted_at": null,
}
```
 **Response Codes:**
 | Code | Description |
 |------|-------------|
 | 201 Created | User successfully initialized; returns the new user object. |
 | 400 Bad Request | JSON is malformed, `instagram_username` is missing, or fails validation. |
 | 401 Unauthorized | Missing, expired, or invalid JWT. |
 | 409 Conflict | A user with the derived andrewid already exists. |
 | 500 Internal Error | Database write failed or server-side exception. |


#### 2.2 Updating a Post
##### 2.2.1 Updating an Image
```
PUT /image/:image_idx
```
**Authentication:** Required (Bearer JWT)
**Description:** Replace the image at the specified index with the new image provided in the request body.
**Path Parameters:**
 - `image_idx`: The zero-indexed integer (0-9) of the image to update
**Request Body:**
 - Content-Type: `image/jpeg` or `image/png`
 - Header: `Authorization: Bearer <JWT>`
 - Body: Raw binary image

 **Response Format:**
 - Status: `200 OK`
 - Body:
 ```json
 {
   "andrewid": "etashjha",
   "image_idx": 3,
   "image_url": "https://example.com/images/etashjha/3.jpg",
   "updated_at": "2026-05-08T15:38:40Z"
 }
 ```
  
 **Response Codes:**
  
 | Code | Description |
 |------|-------------|
 | 200 OK | Image successfully updated; returns the updated image metadata. |
 | 400 Bad Request | `image_idx` is out of range or Content-Type is not `image/jpeg` or `image/png`. |
 | 401 Unauthorized | Missing, expired, or invalid JWT. |
 | 403 Forbidden | Authentication successful, but user does not have permission to modify this post. |
 | 500 Internal Error | Filesystem write failed or server-side exception. |


##### 2.2.2 Swapping Images
```
PATCH /swap/:image_idx1/:image_idx2
```
**Authentication:** Required (Bearer JWT)
**Description:** Swap the image at the specified indices with each other.
**Path Parameters:**
 - `image_idx1`: The zero-indexed integer (0-9) of the first image to swap
 - `image_idx2`: The zero-indexed integer (0-9) of the second image to swap

 **Request Body:**
 - Header: `Authorization: Bearer <JWT>`
 - Body: None
 
 **Constraints:**
 - Both indices must be integers in the range 0–9, otherwise returns `400 Bad Request`.
 - `image_idx1` and `image_idx2` must not be equal, otherwise returns `400 Bad Request`.
 **Response Format:**
 - Status: `200 OK`
 - Body:
 ```json
 {
   "andrewid": "etashjha",
   "image_idx1": 2,
   "image_idx2": 5,
   "updated_at": "2026-05-08T15:38:40Z"
 }
 ```
 **Response Codes:**
  
 | Code | Description |
 |------|-------------|
 | 200 OK | Image successfully updated; returns the updated image metadata. |
 | 400 Bad Request | `image_idx1` or `image_idx2` is out of range |
 | 401 Unauthorized | Missing, expired, or invalid JWT. |
 | 403 Forbidden | Authentication successful, but user does not have permission to modify this post. |
 | 500 Internal Error | Filesystem write failed or server-side exception. |

##### 2.2.3 Removing an Image
```
DELETE /image/:image_idx
```
**Authentication:** Required (Bearer JWT)
**Description:** Remove the image at the specified index from the user's post. This WILL NOT update indices, it will leave a sparse spot. 

**Path Parameters:**
 - `image_idx`: The zero-indexed integer (0-9) of the image to remove, must be an index that is present on the server
**Request Body:**
 - Header: `Authorization: Bearer <JWT>`
 - Body: None

 **Constraints:**
 - Index must be integers in the range 0–9 and present on the server, otherwise returns `400 Bad Request`.

 **Response Format:**
 - Status: `200 OK`
 - Body:
 ```json
 {
   "andrewid": "etashjha",
   "image_idx": 2,
   "deleted_at": "2026-05-08T15:38:40Z"
 }
 ```

 **Response Codes:**
  
 | Code | Description |
 |------|-------------|
 | 200 OK | Image successfully updated; returns the updated image metadata. |
 | 400 Bad Request | `image_idx` is out of range |
 | 401 Unauthorized | Missing, expired, or invalid JWT. |
 | 403 Forbidden | Authentication successful, but user does not have permission to modify this post. |
 | 500 Internal Error | Filesystem write failed or server-side exception. |

##### 2.2.4 Updating the Caption 
```
PUT /caption
```
**Authentication:** Required (Bearer JWT)
**Description:** Replace the description of the user's post with the new description provided in the request body.
**Request Body:**
 - Content-Type: `application/json`
 - Header: `Authorization: Bearer <JWT>`
 - Body: 

 ```json
 { 
   "caption": "The user's caption capped at 2,200 characters" 
 }
```
**Constraints:**
  - `caption` is capped at 2,200 characters and will return a `400 Bad Request` if exceeded.

**Response Format:**
 - Status: `200 OK`
 - Body: 
 ```json
 {
   "andrewid": "etashjha",
   "caption": "The user's caption capped at 2,200 characters",
   "updated_at": "2026-05-08T15:38:40Z"
 }
 ```

**Response Codes:**
| Code | Description |
|------|-------------|
| 200 OK | Description successfully updated; returns the updated post object. |
| 400 Bad Request | Request body exceeds 2,200 characters or JSON is malformed. |
| 401 Unauthorized | Missing, expired, or invalid JWT. |
| 403 Forbidden | Authentication successful, but user does not have permission to modify this post. |
| 500 Internal Error | Database update failed or server-side exception. |.

#### 2.3 Retrieving User/Post Data
```
GET /me
```
**Authentication:** Required (Bearer JWT)
 
**Description:** Returns the full profile and post state for the authenticated user, including identity fields, ordered images, caption, queue status, and posted status. The andrewid is derived from the `preferred_username` claim in the verified JWT.
 
**Request Body:**
- Header: `Authorization: Bearer <JWT>`
- Body: None
**Response Format:**
- Status: `200 OK`
- Body:
```json
{
  "andrewid": "etashjha",
  "name": "Etash", 
  "major": "Computer Science", 
  "hometown": "Pittsburgh, PA", 
  "instagram_username": "etashjha",
  "caption": "The user's caption capped at 2,200 characters",
  "image_count": 10
  "queued_at": "2026-05-12T23:48:32.257Z",
  "queue_position": 3
  "posted_at": null,
}
```
 
**Notes:**
- `queue_at` and `posted_at` may be null if the action has not been performed
- `queue_position` is the number of posts that must be completed before this user is posted, or negative one if not queued OR already posted

**Response Codes:**
 
| Code | Description |
|------|-------------|
| 200 OK | User data successfully retrieved; returns the full user object. |
| 401 Unauthorized | Missing, expired, or invalid JWT. |
| 404 Not Found | No user record exists for the derived andrewid. |
| 500 Internal Error | Database read failed or server-side exception. |


#### 2.4 Queuing a Post
```
POST /enqueue
```
**Authentication:** Required (Bearer JWT)
 
**Description:** Uses JWT to queue a post for Instagram posting. The queue status is updated to `"pending"` and the post is added to the queue. Returns the updated queue status.
 
**Request Body:**
- Header: `Authorization: Bearer <JWT>`
- Body: None
**Response Format:**
- Status: `200 OK`
- Body:
```json
{
  "andrewid": "etashjha",
  "queue_status": "pending",
  "posted": false,
  "queue_position": 3,
  "queued_at": "2026-05-08T15:38:40Z",
  "scheduled": "2026-05-08T16:40:00Z"
}
```
 
**Notes:**
- The `queue_position` is the position of the post in the queue, starting from 0 and will strictly decrease as posts are made.
- The `queued_at` timestamp indicates when the post was added to the queue and is fixed.
- The `scheduled` timestamp indicates when the post is scheduled to be posted and may be updated based on internal queue considerations covered in the "Queue Processing" section of this RFC.

**Response Codes:**
| Code | Description |
|------|-------------|
| 200 OK | User data successfully retrieved; returns the full user object. |
| 401 Unauthorized | Missing, expired, or invalid JWT. |
| 404 Not Found | No user record exists for the derived andrewid. |
| 409 Conflict | The user is already queued for a post. |
| 500 Internal Error | Database read failed or server-side exception. |

#### 2.5 Retrieving a Post's Queue Status
```
GET /status
```
**Authentication:** Required (Bearer JWT)
 
**Description:** A convenience method which returns the queue status of a post, a subset of the information available in the user metadata endpoint `/me`. 
 
**Request Body:**
- Header: `Authorization: Bearer <JWT>`
- Body: None
**Response Format:**
- Status: `200 OK`
- Body:
```json
{
  "andrewid": "etashjha",
  "queue_status": "pending",
  "posted": false,
  "queue_position": "3",
  "queued_at": "2026-05-08T15:38:40Z",
  "scheduled": "2026-05-08T16:40:00Z"
}
```
 
**Notes:**
- The `queue_position` is the position of the post in the queue, starting from 0 and will strictly decrease as posts are made.
- The `queued_at` timestamp indicates when the post was added to the queue and is fixed.
- The `scheduled` timestamp indicates when the post is scheduled to be posted and may be updated based on internal queue considerations covered in the "Queue Processing" section of this RFC.

**Response Codes:**
| Code | Description |
|------|-------------|
| 200 OK | User data successfully retrieved; returns the full user object. |
| 401 Unauthorized | Missing, expired, or invalid JWT. |
| 404 Not Found | No user record exists for the derived andrewid. |
| 500 Internal Error | Database read failed or server-side exception. |

#### 2.6 Canceling a Post
```
DELETE /dequeue
```
**Authentication:** Required (Bearer JWT)
 
**Description:** Uses JWT to dequeue a post for Instagram posting. The queue status is updated to `"none"` and the post is removed from the queue. If the post is already posted or not on the queue, the request is ignored and error is returned.
 
**Request Body:**
- Header: `Authorization: Bearer <JWT>`
- Body: None
**Response Format:**
- Status: `200 OK`
- Body:
```json
{
  "andrewid": "etashjha",
  "queue_status": "none",
  "posted": false
}
```

#### 2.7 Authenticating OIDC with a callback
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
 
**Notes:**
- Note that the error code varies for jobs that were already posted and jobs that were not on the queue.

**Response Codes:**
| Code | Description |
|------|-------------|
| 200 OK | User data successfully retrieved; returns the full user object. |
| 400 Bad Request | The user was not found on the queue.  |
| 401 Unauthorized | Missing, expired, or invalid JWT. |
| 404 Not Found | No user record exists for the derived andrewid. |
| 409 Conflict | The post was already posted and thus not on the queue. |
| 500 Internal Error | Database read failed or server-side exception. |

### Go Implementation Details
To implement the cmuinsta authentication flow using Gin, we will primarily rely on the framework's middleware pattern to enforce the security requirements of RFC 4. By integrating `golang-jwt/jwt/v5` for cryptographic parsing and `coreos/go-oidc/v3` for automated public key discovery from the ScottyLabs Keycloak instance, you can create a stateless "gatekeeper" that intercepts requests to your `/api/v1/` route group. This architecture ensures that sensitive endpoints, such as `/create`, never receive spoofed data; instead, they securely extract the verified preferred_username claim directly from the Gin context, allowing your business logic to remain clean, modular, and focused on core functionality.

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
