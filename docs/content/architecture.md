---
title: "Architecture"
date: 2023-09-05T02:24:03+01:00
draft: false
---

> NOTE
> This document is technical in nature, and serves as architecture guide to
> those who wish to improve SaferPlace. If you don't know what this means, or
> are not interested, feel free to ignore this page.

## Components

SaferPlace backend is composed of multiple modular components. You can run all
the components in a single binary (monolithic mode, useful for development) or
you can run each component in its own process.

### Report

This component listens to new incidents and pushes them to the queue

### Uploader

Pushes Uploaded Images to the Storage

### Consumer

Headless component consuming the incident from the queue, inserting the data
into the database and notifies the reviewer about a new incident

### Review

Allow the reviewer to interact with the incident and updates the data in the
database.

### Viewer

Lists all incidents for a region as well as individual incidents.

---

## Interaction Diagram

SaferPlace is divided into multiple components, each serving their own role.
Below is a diagram of  all the components and their interactions. We will
discuss each point individually.

![SaferPlace Architectural Diagram](/images/architecture-diagram-v0.1.svg)

### 1a - Image Upload

Image upload is an optional feature in SaferPlace, designed to provide a better
reference when reporting. The user submits the image to the uploader service,
which in turn generates a UUID for the image and writes it to a storage bucket.
It then returns the UUID which is then added to the report data.

### 1b - Report Incident

User then submits the report, with all relevant data to the report service. The
report service performs basic validation of the report, ensuring that the
location is in on earth, and removes any fields which we might not want, but
were included in the report for whatever reason. The UUID of the incident is
then returned to the user, so they can track it in the application.

Note: the image UUID and the incident UUID are not related.

### 2a - Bucket Upload

The image is uploaded to the storage bucket, using the internally configured
credentials.

### 2b - Push incident onto the queue

The incident is then pushed onto the incoming incident queue. The queue serves
as a buffer for all incoming incidents in cases where the rest of the
infrastructure is not ready to consume the incident.

### 3 - Consume Incident

The incident is then consumed from the queue by the Consumer. The consumer
inserts data into the database, as well as notifies the reviewer that a new
incident is up for review

### 4a - Incident Database Insertion

Incident is inserted into the database

### 4b - Notify Reviewer

The consumer sends a message to the notifier which is responsible sending a
message that a new incident is up for review to third parties (instant messaging
platform, push notifications, email etc.)

### 5 - Reviewer Is notified

Reviewer is notified about the incident, and is provided with a link to access
the review UI.

### 6 - Incident Review

The reviewer uses the review UI to see the incident details. They judge the
incident based on review criteria, and can add the following resolution:

| Action | Description |
|--------|-------------|
| Reject | The incident report is not up to our standards |
| Ignore | Keep the incident up for review, but comment on the incident for having a second opinion. |
| Accept | Incident is accepted, and will show up for users. |
| Alert  | Incident has implications right now to people in the area, and if they can, they should take caution. |

The reviewer can also add further comments to each incident.

### 7 - Update Incident Details

The reviewer added their resolution, and the incident data is updated in the
database (resolution, comments).

### 8 and 9a - User Views the incident

User views incidents in the area. A series of requests are made for each region.
A region is a 2D grid of boundary boxes which show all incidents in them. This
is so that we don't get a specific user location, the results can be cached,
and split the requests between different viewers.

### 9b - View Incident Image

When a user is viewing a specific image, they can see the image uploaded. They
contact the bucket directly as it has anonymous read access.
