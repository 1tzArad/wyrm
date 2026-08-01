# Wyrm Auth Test Client

A minimal standalone web client for testing the authentication endpoints of the Wyrm backend. It only contains a login page and a register page.

## Backend API

The client targets the backend's auth handlers:

- `POST /register` — body: `{ "username": "...", "password": "..." }`
- `POST /login` — body: `{ "username": "...", "password": "..." }`

Responses use the backend `Response` envelope:

```json
{ "success": true, "data": { "message": "..." }, "error": null }
{ "success": false, "data": null, "error": { "code": "...", "message": "..." } }
```

Login sets a JWT in an HttpOnly cookie named `token`. The client sends requests with `credentials: "include"` so the browser stores and sends the cookie automatically; it cannot be read from JavaScript.

## Getting started

### 1. Install dependencies

```sh
npm install
```

### 2. Configure the backend URL

Copy `.env.example` to `.env` if you want to override the default:

```
VITE_API_BASE_URL=
```

- **Empty (default):** requests go to the Vite dev server, which proxies `/register` and `/login` to `http://localhost:8080`. Recommended for local development, because the backend does not enable CORS.
- **Absolute URL:** e.g. `http://localhost:8080` calls the backend directly. This requires the backend to allow cross-origin requests.

### 3. Run the dev server

```sh
npm run dev
```

Open the printed URL (default `http://localhost:5173`).

### Build

```sh
npm run build
```

The output is written to `client/dist`.

### Other scripts

- `npm run preview` — serve the built files locally.
- `npm run typecheck` — type-check the client with `tsc`.
