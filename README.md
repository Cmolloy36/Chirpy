# Chirpy: an HTTP Server

Welcome to Chirpy! This is an HTTP server that definitely is not built off of old Twitter.

## Endpoints
- `/app/`
- `GET /api/healthz`
- `POST /api/chirps`
- `DELETE /api/chirps/{chirpID}`
    - Description: Delete chirp with specified ID from database.
    - Request format: `delete http://localhost:8080/api/chirps/{chirpID}`
    - Input body format: N/A
    - Arguments: `{chirpID}`
- `GET /api/chirps/{chirpID}`
- `GET /api/chirps`
    - Description: Retrieve chirps from database.
    - Request format: `get http://localhost:8080/api/chirps`
    - Input body format: N/A
    - Optional Queries: `sort={"asc" or "desc"}`, `author_id={user uuid}`
        - `sort`: sorts chirps by ascending or descending chronology (asc=oldest first). If no parameter is provided, chirps will be returned in ascending order.
        - `author_id`: returns chirps from author specified by author_id. If no chirps are found from that author, returns `404 Not Found`. 

- `POST /api/login`
- `POST /api/polka/webhooks`
- `POST /api/refresh`
- `POST /api/revoke`
- `POST /api/users`
- `PUT /api/users`
- `GET /api/users/{userID}`
- `GET /admin/metrics`
- `POST /admin/reset`

## Future Improvements

- [ ] Finalize Endpoint descriptions
- [ ] Add more REST Client tests to better understand how it works (learn to automate this?)