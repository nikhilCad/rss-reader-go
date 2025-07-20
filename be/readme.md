# RSSReaderGo Backend

## Caching & Refresh Behavior

- **All article data is always served from the local SQLite database (the cache).**
- **The `/posts` endpoint** returns articles from the cache. It never fetches live from the internet.
- **The `/refresh` endpoint** (POST) is the only way to fetch new data from the internet. When called, it fetches all feeds, parses them, and updates the database with new or changed articles.
- After a refresh, the next `/posts` call will return the updated articles from the cache.
- This design avoids rate limits and ensures fast, reliable article loading for the frontend.
- No live fetches happen on normal article loads; only the refresh endpoint triggers a real fetch.

## Endpoints

- `GET /posts`  
  Returns all cached articles as JSON.  
  Response format:

  ```json
  {
    "fromCache": true,
    "articles": [ ... ]
  }
  ```

- `POST /refresh`  
  Triggers a background fetch of all feeds and updates the cache. Returns 204 No Content.

## Why this design?

- Prevents hitting feed sites too often (avoids rate limits).
- Ensures fast article loading for users.
- Allows for offline reading of previously fetched articles.

---
