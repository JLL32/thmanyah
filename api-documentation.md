# Thmanyah API Documentation

## Base URL
All API endpoints are prefixed with `/v1`

## Authentication
*Authentication details to be added when implemented*

## Response Format
All responses are returned in JSON format with the following structure:

### Success Response
```json
{
  "video": {...},           // For single video operations
  "videos": [...],          // For list operations
  "metadata": {...},        // For paginated list operations
  "message": "..."          // For delete operations
}
```

### Error Response
```json
{
  "error": {
    "message": "Error description",
    "details": {...}        // Optional validation errors
  }
}
```

## Video Data Model

### Video Object
```json
{
  "video_id": "string",
  "title": "string",
  "description": "string",
  "type": "string",
  "language": "string",
  "length": 0,
  "published_at": "2023-01-01T00:00:00Z",
  "version": 1
}
```

### Field Descriptions
- `video_id`: Unique identifier for the video
- `title`: Video title
- `description`: Video description
- `type`: Video type/category
- `language`: Video language code
- `length`: Video duration in seconds
- `published_at`: Publication timestamp in RFC3339 format
- `version`: Version number for optimistic locking (read-only)

## Endpoints

### 1. Create Video
**POST** `/v1/videos`

Creates a new video record.

#### Request Body
```json
{
  "video_id": "vid_123",
  "title": "Sample Video",
  "description": "This is a sample video description",
  "type": "educational",
  "language": "en",
  "length": 300,
  "published_at": "2023-01-01T00:00:00Z"
}
```

#### Response
**Status: 201 Created**
```json
{
  "video": {
    "video_id": "vid_123",
    "title": "Sample Video",
    "description": "This is a sample video description",
    "type": "educational",
    "language": "en",
    "length": 300,
    "published_at": "2023-01-01T00:00:00Z",
    "version": 1
  }
}
```

#### Headers
- `Location`: `/v1/videos/{video_id}`

### 2. Get Video
**GET** `/v1/videos/{id}`

Retrieves a specific video by its ID.

#### Parameters
- `id` (path): Video ID

#### Response
**Status: 200 OK**
```json
{
  "video": {
    "video_id": "vid_123",
    "title": "Sample Video",
    "description": "This is a sample video description",
    "type": "educational",
    "language": "en",
    "length": 300,
    "published_at": "2023-01-01T00:00:00Z",
    "version": 1
  }
}
```

#### Error Responses
- **404 Not Found**: Video not found

### 3. Update Video
**PATCH** `/v1/videos/{id}`

Updates an existing video. Only provided fields will be updated.

#### Parameters
- `id` (path): Video ID

#### Headers (Optional)
- `X-Expected-Version`: Expected version number for optimistic locking

#### Request Body
```json
{
  "title": "Updated Video Title",
  "description": "Updated description",
  "type": "tutorial",
  "length": 450,
  "language": "es",
  "published_at": "2023-02-01T00:00:00Z"
}
```

#### Response
**Status: 200 OK**
```json
{
  "video": {
    "video_id": "vid_123",
    "title": "Updated Video Title",
    "description": "Updated description",
    "type": "tutorial",
    "language": "es",
    "length": 450,
    "published_at": "2023-02-01T00:00:00Z",
    "version": 2
  }
}
```

#### Error Responses
- **404 Not Found**: Video not found
- **409 Conflict**: Version mismatch (when using optimistic locking)

### 4. Delete Video
**DELETE** `/v1/videos/{id}`

Deletes a specific video.

#### Parameters
- `id` (path): Video ID

#### Response
**Status: 200 OK**
```json
{
  "message": "video successfully deleted"
}
```

#### Error Responses
- **404 Not Found**: Video not found

### 5. List Videos
**GET** `/v1/videos`

Retrieves a paginated list of videos with optional filtering.

#### Query Parameters
- `title` (string, optional): Filter by title (partial match)
- `description` (string, optional): Filter by description (partial match)
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Number of items per page (default: 20)
- `sort` (string, optional): Sort field (default: "id")

#### Sort Options
Available sort fields (prefix with `-` for descending order):
- `video_id`
- `title`
- `description`
- `length`
- `type`

Examples:
- `sort=title` (ascending by title)
- `sort=-length` (descending by length)

#### Example Request
```
GET /v1/videos?title=tutorial&page=1&page_size=10&sort=-published_at
```

#### Response
**Status: 200 OK**
```json
{
  "metadata": {
    "current_page": 1,
    "page_size": 10,
    "first_page": 1,
    "last_page": 5,
    "total_records": 50
  },
  "videos": [
    {
      "video_id": "vid_123",
      "title": "Sample Video",
      "description": "This is a sample video description",
      "type": "educational",
      "language": "en",
      "length": 300,
      "published_at": "2023-01-01T00:00:00Z",
      "version": 1
    }
  ]
}
```

### 6. Health Check
**GET** `/v1/healthcheck`

Returns the health status of the API.

#### Response
**Status: 200 OK**
```json
{
  "status": "available",
  "system_info": {
    "environment": "development",
    "version": "1.0.0"
  }
}
```

## Error Codes

### HTTP Status Codes
- **200 OK**: Request successful
- **201 Created**: Resource created successfully
- **400 Bad Request**: Invalid request data
- **404 Not Found**: Resource not found
- **405 Method Not Allowed**: HTTP method not supported for this endpoint
- **409 Conflict**: Resource conflict (e.g., version mismatch)
- **422 Unprocessable Entity**: Validation errors
- **500 Internal Server Error**: Server error

### Common Error Response Examples

#### Validation Error (422)
```json
{
  "error": {
    "message": "Validation failed",
    "details": {
      "title": "must be provided",
      "length": "must be greater than zero"
    }
  }
}
```

#### Not Found Error (404)
```json
{
  "error": {
    "message": "the requested resource could not be found"
  }
}
```

## Rate Limiting
*Rate limiting details to be added when implemented*

## Pagination

The list videos endpoint supports cursor-based pagination:

- Use `page` and `page_size` parameters to control pagination
- The `metadata` object in the response contains pagination information
- Maximum `page_size` is typically limited (check with backend team for current limit)

## Best Practices

1. **Always handle errors gracefully** - Check HTTP status codes and parse error responses
2. **Use optimistic locking for updates** - Include the `X-Expected-Version` header when updating videos to prevent conflicts
3. **Implement proper pagination** - Don't assume all results fit in a single page
4. **Cache responses when appropriate** - Videos don't change frequently, consider caching GET requests
5. **Validate input on the frontend** - While the API validates input, frontend validation improves user experience

## Frontend Integration Examples

### JavaScript/Fetch Example
```javascript
// Create a video
const createVideo = async (videoData) => {
  const response = await fetch('/v1/videos', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(videoData)
  });
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  
  return await response.json();
};

// Get a video with error handling
const getVideo = async (videoId) => {
  try {
    const response = await fetch(`/v1/videos/${videoId}`);
    
    if (response.status === 404) {
      return null; // Video not found
    }
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    return data.video;
  } catch (error) {
    console.error('Error fetching video:', error);
    throw error;
  }
};

// Update video with optimistic locking
const updateVideo = async (videoId, updates, currentVersion) => {
  const headers = {
    'Content-Type': 'application/json',
  };
  
  if (currentVersion) {
    headers['X-Expected-Version'] = currentVersion.toString();
  }
  
  const response = await fetch(`/v1/videos/${videoId}`, {
    method: 'PATCH',
    headers,
    body: JSON.stringify(updates)
  });
  
  if (response.status === 409) {
    throw new Error('Video was modified by another user. Please refresh and try again.');
  }
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  
  return await response.json();
};
```

## Notes for Frontend Developers

1. **Date Handling**: All timestamps are in RFC3339 format. Use appropriate date parsing libraries in your frontend framework.

2. **Video Length**: The `length` field represents duration in seconds. You may want to format this for display (e.g., "5:00" for 300 seconds).

3. **Partial Updates**: The PATCH endpoint only updates provided fields. Omitted fields remain unchanged.

4. **Error Handling**: Always implement proper error handling for network requests and API errors.

5. **Loading States**: Implement loading states for all API calls to improve user experience.

6. **Version Tracking**: Store the video version when fetching videos to use for optimistic locking during updates.