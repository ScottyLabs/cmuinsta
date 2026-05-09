CMU Insta RFC 3
# Database Model
 - Author: @etashj
 - Created: 2026-05-08
 - Updated: 2026-05-09

## Overview
Storing student data and their images requires a robus databse and storage model. This RFC will outline that process using file storage on our server for images and a realtional database with ORM to keep track of registered students. 

## Motivations
In order to crate an effective database for this project, we must lay out a clear model. This includes, 
 - The relational database details (ORM models, columns)
 - The storage model for images and descriptions (file storage on our server)
 - Describing design choices and limitations where necesary
This will not only allow for a clear plan to reference when implementing the databse, it will serve as a ground truth for the database schema and storage model for debugging and testing.

## Goals
 - Define the database schema and storage model
   - This includes filepaths for images
   - The row and column definitions for the relational database (already chosen to be SQLite)
- Choosing an ORM for use in Go with the relational database. 

## Non Goals
Code and implementation details will not be described in this RFC. It will provide a higher level overview of the database schema and storage model to begin with. 

## Detailed Design

### Database Layout
Since the database is rather simple, we will forgo a relational flowchart in favor of a simple table schema. Here is the schema for the `students` table with some sample entries. 

| AndrewID | Instagram | Queued | Posted | 
|-----------|-----------|-------|-------|
| etashjha |  etashj         |  NULL     | NULL      |
| kdass | kritd          |  2026-05-08T08:00:00Z     |  NULL     |
| anishp | ap-1          |  2026-05-06T08:00:00Z     |   2026-05-07T11:20:00Z     |

Using this schema, we can store the queue and post history for each student, as well as their Instagram username. A NULL field indicates that an event has not yet occurred, it follows that if the Posted field is non-NULL, then the Queued field must also be non-NULL. This will be abstracted using an ORM. 

### File Storage
On the server side, there will exist some `store` directory in which subdirectories corresponding to andrew IDs will be created. Each subdirectory will contain the student's images and descriptions. Images will be stored as enumeratedJPEG files at the isntagram resolution and the description as plain text. We will have a tree like this, with up to 10 images per student.

```
store/
├── etashjha/
│   ├── 0.jpg
│   ├── 1.jpg
│   └── description.txt
├── kdass/
│   ├── 0.jpg
│   ├── 1.jpg
│   └── description.txt
└── anishp/
    ├── 0.jpg
    ├── 1.jpg
    └── description.txt
```

### ORM
We will use [GORM](https://gorm.io/) since it is a popular ORM for Go and supports SQLite out of the box. Providing a simple, schema-based interface for database operations. We are not concerned with performance for this task. 

## Open Questions
 - Will we need to concern ourselves with speed? Should this be optimized for speed?
 - When should we clean up the files? After posting or after all posts are processed?

## Implementation Phases
 - The database will be implemented in parallel with the backend and tested simultaneously.
