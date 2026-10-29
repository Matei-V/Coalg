# Coalg

Coalg is a web application backend built with Go, backed by a MySQL/MariaDB relational database, and integrated with a Python/Pygments-based code syntax highlighting service. It manages educational modules, chapters, lessons, associated resources, and problem sets, along with user authentication and an administrative approval workflow.

---

## Technical Architecture

The backend consists of four primary subsystems:

1. **HTTP Server & Routing (`net/http`)**: Handles incoming HTTP requests, session management via HTTP cookies, form processing, dynamic template rendering, and static asset delivery.
2. **Database Layer (`database/sql` + `go-sql-driver/mysql` + `godotenv`)**: Manages relational data persistence, schema queries, record insertion, auto-increment tracking, and moderation state transitions. Database credentials are read dynamically from an environment configuration file (`.env`).
3. **Authentication & Session Manager (`golang.org/x/crypto/bcrypt`)**: Handles credential hashing, verification, cryptographically secure 32-byte session token generation (`crypto/rand`), and role-based access control (standard user vs. administrator).
4. **Code Syntax Highlighting Pipeline**: Intercepts C++ code blocks within submitted lesson content and invokes Pygments via a Python subprocess to output formatted code before storing it in the database.

---

## Technology Stack & Dependencies

- **Language**: Go (v1.23+)
- **Database Engine**: MySQL 8.x or MariaDB 10.x / 11.x
- **Database Driver**: `github.com/go-sql-driver/mysql` (v1.8.1)
- **Environment Management**: `github.com/joho/godotenv` (v1.5.1)
- **Password Hashing**: `golang.org/x/crypto/bcrypt` (v0.31.0)
- **Template Engine**: Standard library `html/template`
- **Subprocess Utilities**: Python 3 with the `Pygments` package

---

## Project Structure

```
Coalg/
├── main.go                     # Server entry point, route definitions, and page handlers
├── go.mod / go.sum             # Go module definition and dependency checksums
├── .env                        # Environment variables (ignored by git: contains DB credentials)
├── db_setup.sql                # MariaDB/MySQL database schema, constraints, and seed data
├── api/
│   ├── api_handler.go          # HTTP API handlers for navigation, lessons, chapters, and approval
│   └── auth/
│       └── auth_hadler.go      # User registration, authentication, cookies, and session verification
├── backend/
│   └── db_handler/
│       ├── db_handler.go       # Core database queries, CRUD helpers, and connection pool initialization
│       └── auth_db_handler.go  # User account queries, session tokens, and admin operations
├── utils/
│   ├── code_parse.go           # Subprocess runner for code block extraction and parsing
│   └── hilite.me/
│       ├── run_python.bash     # Shell wrapper script executing Pygments
│       └── tools.py            # Python wrapper utilizing Pygments with C++ lexer
└── static/                     # Server-side HTML templates and assets
```

---

## Database Architecture (`coalgdb`)

The relational schema is configured with the following tables and relationships:

```
  ┌──────────────┐
  │   unitate    │
  │ (Grade Level)│
  └──────┬───────┘
         │ 1:N
         ▼
  ┌──────────────┐
  │   capitol    │
  │  (Chapter)   │
  └──────┬───────┘
         │ 1:N (CASCADE)
         ▼
  ┌──────────────┐
  │    pagina    │
  │   (Lesson)   │
  └──────┬───────┘
         ├──────────────────────────┐
         │ 1:N (CASCADE)            │ 1:N (CASCADE)
         ▼                          ▼
  ┌──────────────┐           ┌──────────────┐
  │   resursa    │           │   problema   │
  │  (Resource)  │           │  (Problems)  │
  └──────────────┘           └──────────────┘

  ┌──────────────┐
  │     user     │
  │  (Accounts)  │
  └──────────────┘
```

### Table Definitions

- **`unitate`**: Represents class / curriculum grade levels (e.g., 9, 10, 11, 12).
  - `id_unit` (INT, Primary Key)
- **`capitol`**: Represents educational chapters belonging to a curriculum grade.
  - `id_cap` (INT, Primary Key, Auto Increment)
  - `titlu` (VARCHAR(64), Not Null)
  - `id_unit` (INT, Foreign Key -> `unitate.id_unit`)
- **`pagina`**: Stores lesson titles, formatted content descriptions, author names, and moderation approval status.
  - `id_pag` (INT, Primary Key, Auto Increment)
  - `titlu` (VARCHAR(64), Not Null)
  - `descriere` (TEXT)
  - `autor` (VARCHAR(64), Defaults to `"anonim"` if user is not authenticated)
  - `aprobat` (TINYINT(1), Default `0`): `0` for unapproved / pending moderation, `1` for approved / published.
  - `id_cap` (INT, Foreign Key -> `capitol.id_cap`, ON DELETE CASCADE, ON UPDATE CASCADE)
- **`resursa`**: External educational links or references tied to a specific lesson.
  - `id_res` (INT, Primary Key, Auto Increment)
  - `titlu` (VARCHAR(64))
  - `link` (VARCHAR(256))
  - `id_pag` (INT, Foreign Key -> `pagina.id_pag`, ON DELETE CASCADE, ON UPDATE CASCADE)
- **`problema`**: Practice programming problems associated with a lesson.
  - `id_prob` (INT, Primary Key, Auto Increment)
  - `titlu` (VARCHAR(64))
  - `link` (VARCHAR(256))
  - `id_pag` (INT, Foreign Key -> `pagina.id_pag`, ON DELETE CASCADE, ON UPDATE CASCADE)
- **`user`**: User identity, credential store, and role privileges.
  - `id_user` (INT, Primary Key, Auto Increment)
  - `name` (VARCHAR(256), Unique user identifier)
  - `parola` (VARCHAR(256), Bcrypt-hashed password string)
  - `admin` (TINYINT, Default `0`): `0` for regular user, `1` for administrator privileges.
  - `sessiontoken` (VARCHAR(255), Active session token string)

---

## Authentication & Security Details

- **Password Storage**: Passwords are encrypted before database insertion using Bcrypt with a cost factor of 10 (`bcrypt.GenerateFromPassword`).
- **Session Tokens**: On successful authentication (`POST /login`), the backend generates a 32-byte cryptographically secure pseudo-random token via `crypto/rand`, base64 URL-safe encodes it, stores it in `user.sessiontoken`, and transmits it via `HttpOnly` cookies scoped to `Path: "/"`.
- **Session Expiry**: Login cookies (`session_token` and `username`) are configured with a 24-hour expiration window; logout expires both cookies immediately (`Max-Age: -1h`).
- **Authorization Enforcement**: Restricted routes (`/add-page`, `/account`) verify the caller via `auth.Authorize(r)`, checking that the request's `username` cookie matches an existing user whose `sessiontoken` in the database matches the provided `session_token` cookie.
- **Role Permissions**: Handlers distinguish administrators via `dbh.IsAdmin()`, granting access to unapproved lessons and moderation endpoints (`/aproved`, `/reject`).

---

## Code Parsing & Highlighting Pipeline

When a lesson is submitted via `POST /add_new_page`:
1. The backend inspects `descriere` for `<code>...</code>` blocks (`utils.Parse_code`).
2. Code contents are isolated and formatted.
3. The helper script path is resolved dynamically via `filepath.Join("utils", "hilite.me", "run_python.bash")`.
4. A subprocess is spawned via `os/exec` executing `run_python.bash`, which changes directory to its own location (`SCRIPT_DIR`) and invokes `tools.py`.
5. `tools.py` invokes `Pygments` with a C++ lexer (`get_lexer_by_name('cpp')`) to generate formatted HTML.
6. **Graceful Fallback**: If the Python execution encounters an error, the system logs the failure and falls back to `<pre><code>...</code></pre>` blocks rather than terminating the server.
7. The resulting HTML replaces the raw code block and is stored directly in `pagina.descriere`.

---

## HTTP API & Route Map

### Page Routing
| Method | Route | Description |
|---|---|---|
| `GET` | `/` | Serves main landing page (`index.html`). |
| `GET` | `/unit?clasa={id}&id={id}` | Serves curriculum unit and lesson view (`unit.html`). |
| `GET` | `/add-page` | Serves lesson submission page (requires valid session cookie via `auth.Authorize`). |
| `GET` | `/account` | Serves user account / moderation dashboard (requires valid session cookie via `auth.Authorize`). |
| `GET` | `/register-page` | Serves account registration form. |
| `GET` | `/login-page` | Serves account login form. |

### Content & Navigation Endpoints
| Method | Route | Parameters | Description |
|---|---|---|---|
| `GET` | `/load-menu` | `clasa` | Returns grade unit navigation shell. |
| `GET` | `/get-menu` | `clasa`, `id` | Fetches chapters and approved lessons for a grade; renders navigation tree. |
| `GET` | `/load-lesson` | `id`, `clasa` | Retrieves lesson data, description, external resources, and problems. |
| `GET` | `/load-chapters` | None | Returns existing chapters across grades 9, 10, and 11. |
| `GET` | `/add_new_chapter` | None | Template helper for chapter creation. |
| `GET` | `/add_to_existing_chapter` | None | Template helper for appending to an existing chapter. |
| `POST` | `/add_res` | `nr_res` | Increments dynamic resource counter and returns field template. |
| `POST` | `/update_nr_res` | `nr_res` | Updates the hidden resource count parameter. |
| `POST` | `/add_pb` | `nr_pb` | Increments dynamic problem counter and returns field template. |
| `POST` | `/update_nr_pb` | `nr_pb` | Updates the hidden problem count parameter. |

### Lesson Management & Moderation
| Method | Route | Description |
|---|---|---|
| `POST` | `/add_new_page` | Processes lesson payload, highlights code blocks, persists resources and problems, sets `aprobat=0`. Falls back to author `"anonim"` if unauthenticated. |
| `GET` | `/my-less` | Returns lessons submitted by the current user (or all pending lessons if caller is admin). |
| `GET` | `/aproved` | Accepts query param `id`; marks lesson as approved (`aprobat = 1`). |
| `GET` | `/reject` | Accepts query param `id`; removes unapproved lesson from database and normalizes sequence IDs. |

### Authentication Endpoints
| Method | Route | Parameters | Description |
|---|---|---|---|
| `POST` | `/register` | `username`, `password` | Validates input (username >= 4 chars, password >= 8 chars), hashes credentials, creates user record, redirects to `/login-page?registered=true`. |
| `POST` | `/login` | `username`, `password` | Validates credentials, issues session cookies (`session_token`, `username`) with `Path: "/"`. |
| `GET` | `/logout` | None | Invalidates active database session token and expires cookies. |
| `GET` | `/auth-cont` | None | Evaluates session token and renders authenticated or unauthenticated nav state. |

---

## MySQL / MariaDB Server Setup & Database Initialization

### 1. Install MariaDB / MySQL

On Debian / Ubuntu:
```bash
sudo apt update
sudo apt install -y mariadb-server mariadb-client
# Or for MySQL:
# sudo apt install -y mysql-server
```

On Arch Linux:
```bash
sudo pacman -S mariadb
sudo mariadb-install-db --user=mysql --basedir=/usr --datadir=/var/lib/mysql
```

On Fedora / RHEL:
```bash
sudo dnf install -y mariadb-server
```

### 2. Start the Database Service

```bash
sudo systemctl start mariadb     # Use 'mysql' if using MySQL server
sudo systemctl enable mariadb
```

Verify the service is active:
```bash
sudo systemctl status mariadb
```

### 3. Create the Database and Configure Root Access

Log into the database server as root:
```bash
sudo mysql -u root
# or:
# sudo mariadb -u root
```

Create the `coalgdb` database:
```sql
CREATE DATABASE IF NOT EXISTS coalgdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Set or update the root password to match your configuration:
```sql
ALTER USER 'root'@'localhost' IDENTIFIED BY 'your_mariadb_password';
GRANT ALL PRIVILEGES ON coalgdb.* TO 'root'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

---

### 4. Running `db_setup.sql` to Setup the Database

The database schema, constraints, cascading foreign keys, and seed data are fully contained in `db_setup.sql`.

#### Option A: Run directly from the command line (Recommended)

From the project root directory, execute the script directly into `coalgdb`:

```bash
# If using MySQL:
mysql -u root -p coalgdb < db_setup.sql

# If using MariaDB:
mariadb -u root -p coalgdb < db_setup.sql
```
You will be prompted to enter your database root password.

#### Option B: Run from inside the MySQL/MariaDB interactive shell

Alternatively, you can import the script using the `SOURCE` command:

```sql
-- 1. Connect to the server
mariadb -u root -p

-- 2. Select the database
USE coalgdb;

-- 3. Execute the SQL file
SOURCE db_setup.sql;
```

#### What `db_setup.sql` Configures:
- **Drops and Recreates Tables**: `capitol`, `pagina`, `problema`, `resursa`, `unitate`, and `user`.
- **Relational Constraints**: Establishes foreign keys with `ON DELETE CASCADE ON UPDATE CASCADE`.
- **Initial Grade Units**: Inserts curriculum units `9`, `10`, `11`, and `12` into `unitate`.
- **Placeholder Rows**: Inserts zero-indexed base records (`id=0`) for chapters, pages, resources, and problems.
- **Initial Chapters**: Seeds base chapters (e.g. `Divide et impera`, `Introducere in C++`).
- **Default User Accounts**:
  - **Administrator**: `admin` (stored with Bcrypt hash, `admin = 1`).
  - **Standard User**: `matei` (stored with Bcrypt hash, `admin = 0`).

#### Verifying the Setup:
Check that all tables were created properly:
```sql
USE coalgdb;
SHOW TABLES;
```
Expected output:
```
+--------------------+
| Tables_in_coalgdb  |
+--------------------+
| capitol            |
| pagina             |
| problema           |
| resursa            |
| unitate            |
| user               |
+--------------------+
```

---

## Environment Configuration (`.env`)

The backend uses `godotenv` to dynamically load environment variables at runtime in `backend/db_handler/db_handler.go`:

```go
godotenv.Load(".env")
db, err := sql.Open("mysql", "root:"+os.Getenv("PASSWORD")+"@tcp(localhost:3306)/coalgdb")
```

Create a `.env` file in the project root (this file is excluded from git via `.gitignore`):

```bash
touch .env
```

Add your database root password:
```env
PASSWORD=your_mariadb_password
```

---

## Server Setup and Execution

### 1. Prerequisites

1. **Go Toolchain**: Go 1.23 or newer (`go version`).
2. **Python 3 & Pygments**: Required by the backend syntax highlighting utility:
   ```bash
   pip install Pygments
   ```
   Ensure the wrapper script is executable:
   ```bash
   chmod +x utils/hilite.me/run_python.bash
   ```

### 2. Install Go Dependencies

Download module dependencies specified in `go.mod` (including `mysql` driver and `godotenv`):
```bash
go mod download
```

### 3. Run the Server

#### Direct Execution:
```bash
go run main.go
```

#### Compile and Run Standalone Binary:
```bash
go build -o coalg_server main.go
./coalg_server
```

Once running, the backend initializes the database connection pool, registers all API and page routes, and starts listening on:
```
http://localhost:8080
```
