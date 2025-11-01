# Churches API Documentation

This document describes the church-related endpoints added to the songbooks-of-praise-backend API.

## Overview

The Churches API allows you to manage churches, their locations (countries and states), elders, services, and events. All endpoints follow RESTful conventions.

## Database Schema

The church system is based on the following tables from `churches.sql`:

- `countries` - List of countries with ISO codes
- `states` - Administrative divisions within countries
- `churches` - Church information and contact details
- `elders` - Church leadership information
- `church_services` - Service schedules for churches
- `church_events` - Church events and activities

## Endpoints

### Location Endpoints

#### Countries

- `GET /countries` - Get all countries
- `GET /countries/with-states` - Get all countries with their states/provinces
- `GET /countries/:id` - Get a specific country by ID
- `GET /countries/:country_id/states` - Get all states for a specific country

#### States/Provinces

- `GET /states` - Get all states/provinces
- `GET /states/:id` - Get a specific state by ID
- `GET /states/:state_id/churches` - Get all churches in a specific state

### Church Endpoints

#### Church CRUD Operations

- `GET /churches` - Get all churches
- `GET /churches/:id` - Get a specific church by ID
- `POST /churches` - Create a new church
- `PUT /churches/:id` - Update an existing church
- `DELETE /churches/:id` - Delete a church

#### Church Elders

- `GET /churches/:church_id/elders` - Get all elders for a church
- `GET /elders/:id` - Get a specific elder by ID
- `POST /elders` - Create a new elder
- `PUT /elders/:id` - Update an existing elder
- `DELETE /elders/:id` - Delete an elder

#### Church Services

- `GET /churches/:church_id/services` - Get all services for a church
- `GET /services/:id` - Get a specific service by ID
- `POST /services` - Create a new service
- `PUT /services/:id` - Update an existing service
- `DELETE /services/:id` - Delete a service

#### Church Events

- `GET /churches/:church_id/events` - Get all events for a church
- `GET /events/:id` - Get a specific event by ID
- `POST /events` - Create a new event
- `PUT /events/:id` - Update an existing event
- `DELETE /events/:id` - Delete an event

## Request/Response Examples

### Create Church

**POST** `/churches`

```json
{
  "name": "Grace Community Church",
  "address": "123 Main Street, Springfield, IL 62701",
  "phone": "+1-217-555-0123",
  "email": "info@gracecommunity.org",
  "description": "A welcoming community focused on worship and service",
  "website": "https://gracecommunity.org",
  "established": "1985-03-15",
  "facebook": "https://facebook.com/gracecommunity",
  "instagram": "https://instagram.com/gracecommunity",
  "youtube": "https://youtube.com/gracecommunity",
  "spotify": "https://spotify.com/gracecommunity",
  "state_id": 15
}
```

**Response:**

```json
{
  "church": {
    "id": 1,
    "name": "Grace Community Church",
    "address": "123 Main Street, Springfield, IL 62701",
    "phone": "+1-217-555-0123",
    "email": "info@gracecommunity.org",
    "description": "A welcoming community focused on worship and service",
    "website": "https://gracecommunity.org",
    "established": "1985-03-15",
    "facebook": "https://facebook.com/gracecommunity",
    "instagram": "https://instagram.com/gracecommunity",
    "youtube": "https://youtube.com/gracecommunity",
    "spotify": "https://spotify.com/gracecommunity",
    "state_id": 15,
    "state": {
      "id": 15,
      "name": "Illinois",
      "country_id": 1,
      "country": {
        "id": 1,
        "name": "United States",
        "iso_alpha2": "US",
        "iso_alpha3": "USA",
        "iso_numeric": "840"
      }
    },
    "elders": [],
    "services": [],
    "events": []
  }
}
```

### Create Elder

**POST** `/elders`

```json
{
  "name": "John Smith",
  "email": "john.smith@gracecommunity.org",
  "phone": "+1-217-555-0124",
  "picture": "https://example.com/photos/john-smith.jpg",
  "church_id": 1
}
```

### Create Service

**POST** `/services`

```json
{
  "service_type": "Sunday Morning Worship",
  "schedule": "Sundays at 10:00 AM",
  "church_id": 1
}
```

### Create Event

**POST** `/events`

```json
{
  "name": "Annual Church Picnic",
  "start_date": "2025-07-15T12:00:00Z",
  "end_date": "2025-07-15T18:00:00Z",
  "location": "Central Park",
  "image": "https://example.com/events/picnic2025.jpg",
  "color": "#4CAF50",
  "recurrent": false,
  "description": "Join us for our annual church picnic with food, games, and fellowship!",
  "church_id": 1
}
```

## Data Models

### Church Model

```go
type Church struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Address     string    `json:"address"`
    Phone       *string   `json:"phone"`
    Email       string    `json:"email"`
    Description *string   `json:"description"`
    Website     *string   `json:"website"`
    Established *string   `json:"established"`
    Facebook    *string   `json:"facebook"`
    Instagram   *string   `json:"instagram"`
    YouTube     *string   `json:"youtube"`
    Spotify     *string   `json:"spotify"`
    StateID     int       `json:"state_id"`
    State       *State    `json:"state,omitempty"`
    Elders      []Elder   `json:"elders,omitempty"`
    Services    []Service `json:"services,omitempty"`
    Events      []Event   `json:"events,omitempty"`
}
```

### Country Model

```go
type Country struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    ISOAlpha2   string  `json:"iso_alpha2"`
    ISOAlpha3   string  `json:"iso_alpha3"`
    ISONumeric  string  `json:"iso_numeric"`
    States      []State `json:"states,omitempty"`
}
```

### State Model

```go
type State struct {
    ID        int      `json:"id"`
    Name      string   `json:"name"`
    CountryID int      `json:"country_id"`
    Country   *Country `json:"country,omitempty"`
    Churches  []Church `json:"churches,omitempty"`
}
```

## Geographic Coverage

The system includes comprehensive global coverage with:

- **118 countries** from all continents
- **2,317 administrative divisions** (states, provinces, regions, etc.)
- Complete coverage for:
  - North America (US states, Canadian provinces, Mexican states, etc.)
  - South America (Brazilian states, Argentine provinces, etc.)
  - Europe (German länder, French regions, Spanish communities, etc.)
  - Asia (Chinese provinces, Japanese prefectures, Indian states, etc.)
  - Africa (South African provinces, Nigerian states, etc.)
  - Oceania (Australian states, New Zealand regions, etc.)

## Error Responses

All endpoints return appropriate HTTP status codes:

- `200 OK` - Successful GET/PUT operations
- `201 Created` - Successful POST operations
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server errors

Error responses include a JSON object with an `error` field:

```json
{
  "error": "Invalid church ID"
}
```

## Authentication

Currently, the church endpoints don't require authentication. You may want to add authentication middleware to protect certain endpoints (POST, PUT, DELETE) in a production environment.
