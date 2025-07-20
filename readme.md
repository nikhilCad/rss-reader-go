# RSS Reader Go

## Backend (Go)

1. Install Go dependencies:
   ```sh
   go mod tidy
   ```
2. Run the backend:
   ```sh
   cd be
   go run *.go
   ```

## Frontend (Preact + Vite)

We are using pre-react to decrease bundle size

1. Install Node.js (if not already installed).
2. Go to the frontend directory:
   ```sh
   cd fe
   ```
3. Install dependencies:
   ```sh
   npm install
   ```
4. Start the development server:

   ```sh
   npm run dev
   ```

   The frontend will be available at the URL shown in the terminal (usually http://localhost:5173).

5. To build for production:
   ```sh
   npm run build
   ```
   The static files will be in `fe/dist/`.

## Notes

- The frontend talks to the Go backend at http://localhost:8080 by default. You may need to configure CORS or proxy settings for API calls.
- For hot reload of the backend, use [Air](https://github.com/cosmtrek/air) or similar tools.
