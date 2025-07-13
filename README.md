1.  **Project setup & run instructions**

    ```bash
    git clone https://github.com/yourorg/backend-challenge.git
    cd backend-challenge
    cp .env.example .env          # fill in MONGO_URI, JWT_SECRET
    go mod download
    go run main.go
    ```

2.  **JWT token usage guide**

    To access protected endpoints, you need a JWT token. Follow these steps:

    **Step 1: Obtain a JWT token by logging in**

    Use the following cURL command (replace credentials as appropriate):

    ```bash
    curl -X POST http://localhost:3000/api/v1/users/login \
      -H 'Content-Type: application/json' \
      -d '{"email":"test2@gmail.com","password":"123456"}'
    ```

    The response will include a `token` field. Copy the value of this token.

    **Step 2: Use the JWT token**

    - **In cURL requests:**  
      Add the following header to your requests:

      ```
      Authorization: Bearer <token>
      ```

      Example:

      ```bash
      curl -X GET http://localhost:3000/api/v1/users \
        -H 'Authorization: Bearer <token>'
      ```

    - **In Swagger UI:**
      1. Open Swagger UI in your browser (see section below).
      2. Click the **Authorize** (lock) icon.
      3. Enter your token as:
         ```
         Bearer <token>
         ```
      4. Click "Authorize". Now you can try out protected endpoints.

3.  **Sample API requests & responses**

    ### Register

    ```bash
    curl -X POST http://localhost:3000/api/v1/users/register \
      -H 'Content-Type: application/json' \
      -d '{"name":"Alice","email":"alice@example.com","password":"secret"}'
    ```

    **Response** (201):

    ```json
    {
      "id": "...",
      "name": "Alice",
      "email": "alice@example.com",
      "createdAt": "2025-07-13T..."
    }
    ```

    ### Login

    ```bash
    curl -X POST http://localhost:3000/api/v1/users/login \
      -H 'Content-Type: application/json' \
      -d '{"email":"alice@example.com","password":"secret"}'
    ```

    **Response** (200):

    ```json
    {
      "token": "<jwt-token>"
    }
    ```

    ### Protected endpoints (e.g., List users)

    ```bash
    curl -X GET http://localhost:3000/api/v1/users \
      -H 'Authorization: Bearer <token>'
    ```

    **Response** (200):

    ```json
    [
      {
        "id": "...",
        "name": "Alice",
        "email": "alice@example.com"
        // ...other fields
      }
    ]
    ```

4.  **🛠️ Using Swagger UI**

    The API provides an interactive Swagger UI for exploring and testing endpoints.

    - **Open Swagger UI:**  
      Navigate to [http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html) in your browser.

    - **Authorize with JWT:**

      1. Click the **Authorize** (lock) icon at the top right of the Swagger UI.
      2. In the popup, enter your JWT token prefixed with `Bearer ` (including the space):
         ```
         Bearer <token>
         ```
      3. Click "Authorize" to save the token for your session.

    - **Try out endpoints:**  
      Click on any endpoint, then click **"Try it out"**. Fill in any required parameters and click **"Execute"** to make requests directly from the browser, using your authorized token for protected endpoints.

5.  **Architecture Overview**
    We follow a **hexagonal (ports & adapters)** architecture to keep business logic decoupled from external frameworks:
    ```
    .
    ├── config/                       # load environment and DB configuration
    ├── domain/                       # domain entities and their validation
    ├── application/                  # application layer (use-cases / business logic)
    │   ├── ports/                    # port interfaces (e.g., UserRepository)
    │   └── usecases/                 # implementation of use-cases
    ├── infrastructure/               # adapters for external systems
        ├── logger/                   # log each time when call api in mongoDB
    │   └── mongo/                    # MongoDB repositories and logging middleware
    ├── delivery/                     # delivery layer: HTTP handlers, middleware, router
    │   ├── handlers/                 # Fiber HTTP handlers
    │   ├── middleware/               # custom middleware (JWT, logging)
    │   └── router/                   # route definitions
    ├── docs/                         # auto-generated Swagger/OpenAPI docs
    ├── tests/                        # unit tests using mock implementations
    ├── .env                          # environment variable definitions
    ├── Dockerfile & docker-compose.yml # containerization configurations
    ├── go.mod & go.sum               # Go module dependencies
    └── main.go                       #
    ```
