# Personal Programming Knowledge Base - Technical Specification

## 1. Project Overview

A web application for storing and organizing programming knowledge, recipes, tutorials, and code snippets. Users can create accounts, write articles, organize them with tags, and search through their personal knowledge base.

---

## 2. Technology Stack

### Backend
- **Language**: Go (Golang)
- **Web Framework**: `gorilla/mux` or `chi/chi` for routing
- **Database Driver**: `mattn/go-sqlite3` (SQLite) or `lib/pq` (PostgreSQL)
- **Template Engine**: `html/template` (Go standard library)
- **Session Management**: `gorilla/sessions`
- **Password Hashing**: `golang.org/x/crypto/bcrypt`

### Frontend
- **HTML5** for structure
- **CSS3** for styling (no frameworks required, but optional: simple.css or water.css)
- **No JavaScript required** (pure server-side rendering)

### Database
- **SQLite** (recommended for simplicity) or **PostgreSQL**

---

## 3. Database Schema

### 3.1 Users Table
```
Table: users
--------------------------------------------------
Column Name      | Type          | Constraints
--------------------------------------------------
id               | INTEGER       | PRIMARY KEY, AUTO_INCREMENT
username         | VARCHAR(50)   | UNIQUE, NOT NULL
email            | VARCHAR(100)  | UNIQUE, NOT NULL
password_hash    | VARCHAR(255)  | NOT NULL
created_at       | TIMESTAMP     | NOT NULL, DEFAULT CURRENT_TIMESTAMP
```

**Indexes:**
- Unique index on `username`
- Unique index on `email`

---

### 3.2 Articles Table
```
Table: articles
--------------------------------------------------
Column Name      | Type          | Constraints
--------------------------------------------------
id               | INTEGER       | PRIMARY KEY, AUTO_INCREMENT
user_id          | INTEGER       | NOT NULL, FOREIGN KEY → users(id)
title            | VARCHAR(255)  | NOT NULL
content          | TEXT          | NOT NULL
created_at       | TIMESTAMP     | NOT NULL, DEFAULT CURRENT_TIMESTAMP
updated_at       | TIMESTAMP     | NOT NULL, DEFAULT CURRENT_TIMESTAMP
```

**Indexes:**
- Index on `user_id`
- Index on `created_at` (for sorting)
- Full-text index on `title` and `content` (for search)

**Foreign Keys:**
- `user_id` REFERENCES `users(id)` ON DELETE CASCADE

---

### 3.3 Tags Table
```
Table: tags
--------------------------------------------------
Column Name      | Type          | Constraints
--------------------------------------------------
id               | INTEGER       | PRIMARY KEY, AUTO_INCREMENT
user_id          | INTEGER       | NOT NULL, FOREIGN KEY → users(id)
name             | VARCHAR(50)   | NOT NULL
created_at       | TIMESTAMP     | NOT NULL, DEFAULT CURRENT_TIMESTAMP
```

**Indexes:**
- Unique composite index on `(user_id, name)` - prevents duplicate tag names per user
- Index on `user_id`

**Foreign Keys:**
- `user_id` REFERENCES `users(id)` ON DELETE CASCADE

---

### 3.4 Article_Tags Table (Junction Table)
```
Table: article_tags
--------------------------------------------------
Column Name      | Type          | Constraints
--------------------------------------------------
article_id       | INTEGER       | NOT NULL, FOREIGN KEY → articles(id)
tag_id           | INTEGER       | NOT NULL, FOREIGN KEY → tags(id)
```

**Indexes:**
- Primary key on `(article_id, tag_id)`
- Index on `tag_id`

**Foreign Keys:**
- `article_id` REFERENCES `articles(id)` ON DELETE CASCADE
- `tag_id` REFERENCES `tags(id)` ON DELETE CASCADE

---

## 4. Application Routes

### 4.1 Public Routes (No Authentication Required)

| Method | Route        | Description                          | Template         |
|--------|--------------|--------------------------------------|------------------|
| GET    | /            | Landing page / home                  | home.html        |
| GET    | /register    | Display registration form            | register.html    |
| POST   | /register    | Process registration                 | (redirect)       |
| GET    | /login       | Display login form                   | login.html       |
| POST   | /login       | Process login                        | (redirect)       |

---

### 4.2 Protected Routes (Authentication Required)

| Method | Route                | Description                          | Template            |
|--------|----------------------|--------------------------------------|---------------------|
| GET    | /logout              | Logout and clear session             | (redirect)          |
| GET    | /dashboard           | User's article dashboard/list        | dashboard.html      |
| GET    | /articles/new        | Display new article form             | article_new.html    |
| POST   | /articles            | Create new article                   | (redirect)          |
| GET    | /articles/:id        | View single article                  | article_view.html   |
| GET    | /articles/:id/edit   | Display edit article form            | article_edit.html   |
| POST   | /articles/:id        | Update article                       | (redirect)          |
| POST   | /articles/:id/delete | Delete article                       | (redirect)          |
| GET    | /search              | Search articles                      | dashboard.html      |
| GET    | /tags/:name          | Filter articles by tag               | dashboard.html      |

**Note:** You can use DELETE method for `/articles/:id/delete` if you implement method override.

---

## 5. Page Templates Structure

### 5.1 Base Layout Template

**File:** `templates/layout/base.html`

**Components:**
- `<!DOCTYPE html>` declaration
- `<head>` section with:
  - Meta tags (charset, viewport)
  - Page title (dynamic)
  - CSS links
- `<body>` with:
  - Header/Navigation bar
  - Flash messages container
  - Main content area (block/placeholder)
  - Footer
- Template blocks:
  - `{{block "title" .}}` - Page title
  - `{{block "content" .}}` - Main content
  - `{{block "scripts" .}}` - Optional page-specific scripts

---

### 5.2 Navigation Bar

**Location:** Inside `base.html` or separate `nav.html`

**If User is Logged In:**
- App logo/name (links to dashboard)
- Dashboard link
- New Article button/link
- Search link
- Username display
- Logout button

**If User is Not Logged In:**
- App logo/name (links to home)
- Login link
- Register link

---

### 5.3 Authentication Templates

#### Register Page
**File:** `templates/register.html`

**Form Fields:**
- Username (input, required)
- Email (input, type="email", required)
- Password (input, type="password", required)
- Confirm Password (input, type="password", required)
- Submit button

**Validation Messages:**
- Display errors for: username taken, email taken, passwords don't match, weak password

---

#### Login Page
**File:** `templates/login.html`

**Form Fields:**
- Username or Email (input, required)
- Password (input, type="password", required)
- Submit button
- Link to registration page

**Validation Messages:**
- Display error for: invalid credentials

---

### 5.4 Article Templates

#### Dashboard
**File:** `templates/dashboard.html`

**Components:**
- Page title: "My Knowledge Base" or "Dashboard"
- Search bar (form with input field)
- "New Article" button (prominent)
- Tag filter section:
  - Display all user's tags as clickable buttons/links
  - "All Articles" option to clear filter
- Article list:
  - Each article displayed as card/row with:
    - Title (clickable, links to article view)
    - Content preview (first 150 characters)
    - Tags (as badges/chips)
    - Created date
    - Edit button (links to edit page)
    - Delete button (form with POST to delete route)
- Empty state message if no articles
- Pagination controls (if implementing pagination)

---

#### New Article Page
**File:** `templates/article_new.html`

**Form Fields:**
- Title (input, required, maxlength=255)
- Content (textarea, required, rows=20)
- Tags (input, placeholder: "linux, ubuntu, virtual-machine")
- Submit button ("Create Article")
- Cancel button (links back to dashboard)

**Helper Text:**
- Explain tag format: "Separate tags with commas"

---

#### View Article Page
**File:** `templates/article_view.html`

**Components:**
- Article title (h1)
- Metadata section:
  - Created date
  - Last updated date (if different from created)
  - Tags (as badges/chips, clickable to filter)
- Article content (formatted with proper spacing/line breaks)
- Action buttons:
  - Edit button (links to edit page)
  - Delete button (form with confirmation)
  - Back to Dashboard button

---

#### Edit Article Page
**File:** `templates/article_edit.html`

**Form Fields:**
- Title (input, pre-filled with current value)
- Content (textarea, pre-filled with current value)
- Tags (input, pre-filled with current tags as comma-separated)
- Submit button ("Update Article")
- Cancel button (links back to article view)

---

### 5.5 Error Templates

#### 404 Not Found
**File:** `templates/404.html`

**Content:**
- "Page Not Found" heading
- Friendly message
- Link back to dashboard or home

---

#### 403 Forbidden
**File:** `templates/403.html`

**Content:**
- "Access Denied" heading
- Explanation message
- Link back to dashboard or home

---

## 6. Application Logic & Workflows

### 6.1 User Registration Flow

**Steps:**
1. User navigates to `/register`
2. System displays registration form
3. User submits form with username, email, password, confirm password
4. Server validates input:
   - All fields present and not empty
   - Email format is valid
   - Username is 3-50 characters, alphanumeric with underscores
   - Password minimum 8 characters
   - Password and confirm password match
   - Username not already taken (check database)
   - Email not already taken (check database)
5. If validation fails:
   - Re-display form with error messages
   - Keep non-password fields filled
6. If validation passes:
   - Hash password using bcrypt (cost factor: 10-12)
   - Insert new user record into database
   - Set success flash message: "Registration successful! Please log in."
   - Redirect to `/login`

**Error Messages:**
- "Username is already taken"
- "Email is already registered"
- "Passwords do not match"
- "Password must be at least 8 characters"
- "Invalid email format"

---

### 6.2 User Login Flow

**Steps:**
1. User navigates to `/login`
2. System displays login form
3. User submits form with username/email and password
4. Server validates:
   - Both fields present and not empty
   - Query database for user by username OR email
5. If user not found:
   - Display error: "Invalid credentials"
   - Re-display form
6. If user found:
   - Compare submitted password with stored hash using bcrypt
7. If password incorrect:
   - Display error: "Invalid credentials"
   - Re-display form
8. If password correct:
   - Create new session
   - Store user ID in session
   - Set session cookie
   - Set success flash message: "Welcome back, [username]!"
   - Redirect to `/dashboard`

**Security Notes:**
- Use generic "Invalid credentials" message (don't reveal if username or password was wrong)
- Implement rate limiting on login attempts (optional but recommended)
- Use secure session cookies (HTTP-only, Secure flag in production)

---

### 6.3 Session Management

**Session Creation (on login):**
- Generate unique session ID
- Store session data:
  - `user_id`: authenticated user's ID
  - `username`: for display purposes
  - `created_at`: session creation time
- Set session cookie with:
  - HTTP-only flag (prevent JavaScript access)
  - Secure flag (HTTPS only in production)
  - SameSite attribute
  - Max age: 24 hours (configurable)

**Authentication Middleware:**
- Function that wraps protected route handlers
- Checks if session exists
- Validates session data contains `user_id`
- If valid:
  - Allow request to proceed
  - Attach user info to request context
- If invalid:
  - Set flash message: "Please log in to continue"
  - Redirect to `/login`

**Session Destruction (on logout):**
- Clear session data
- Invalidate session ID
- Delete session cookie
- Redirect to home page

---

### 6.4 Article Creation Flow

**Steps:**
1. User clicks "New Article" button
2. System verifies authentication (middleware)
3. System displays article creation form
4. User fills title, content, and tags (optional)
5. User submits form
6. Server validates:
   - Title is not empty, max 255 characters
   - Content is not empty
   - Tags format (if provided)
7. If validation fails:
   - Re-display form with errors
   - Keep form data filled
8. If validation passes:
   - Begin database transaction
   - Insert article record with user_id, title, content, timestamps
   - Get new article ID
   - Process tags:
     - Split tag input by comma
     - Trim whitespace from each tag
     - For each tag:
       - Check if tag exists for this user
       - If not, create new tag record
       - Get tag ID
       - Insert record into article_tags junction table
   - Commit transaction
   - Set success flash message: "Article created successfully!"
   - Redirect to `/articles/:id` (view new article)

**Tag Processing Details:**
- Tag names should be case-insensitive (convert to lowercase)
- Maximum tag length: 50 characters
- Remove duplicate tags from input
- Limit maximum tags per article: 10 (optional)

---

### 6.5 Article View Flow

**Steps:**
1. User clicks on article or navigates to `/articles/:id`
2. System verifies authentication
3. Server queries database:
   - Fetch article by ID
   - Verify article.user_id matches authenticated user
4. If article not found OR doesn't belong to user:
   - Display 404 page OR 403 forbidden page
5. If article found and authorized:
   - Fetch associated tags from tags table via junction table
   - Render article view template with:
     - Article data (title, content, dates)
     - Tags array
     - Edit/delete action buttons

---

### 6.6 Article Update Flow

**Steps:**
1. User clicks "Edit" button on article
2. System verifies authentication
3. Server fetches article by ID and verifies ownership
4. If authorized:
   - Fetch article data and associated tags
   - Format tags as comma-separated string
   - Render edit form pre-filled with current data
5. User modifies title, content, and/or tags
6. User submits form
7. Server validates input (same as creation)
8. If validation passes:
   - Begin database transaction
   - Update article record: title, content, updated_at timestamp
   - Process tags:
     - Delete all existing entries from article_tags for this article
     - Parse new tags from input
     - Create/fetch tag records
     - Insert new article_tags records
   - Commit transaction
   - Set success flash message: "Article updated successfully!"
   - Redirect to `/articles/:id` (view updated article)

**Alternative Tag Update Strategy:**
- Compare old tags with new tags
- Only delete removed tags
- Only add new tags
- More efficient but more complex to implement

---

### 6.7 Article Deletion Flow

**Steps:**
1. User clicks "Delete" button on article
2. Browser shows JavaScript confirmation (optional, or server-side confirmation page)
3. Form submits POST to `/articles/:id/delete`
4. System verifies authentication
5. Server fetches article and verifies ownership
6. If authorized:
   - Begin database transaction
   - Delete related records from article_tags (cascade)
   - Delete article record
   - Commit transaction
   - Set success flash message: "Article deleted successfully"
   - Redirect to `/dashboard`
7. If not authorized:
   - Display 403 forbidden page

**Note:** If using ON DELETE CASCADE in foreign keys, article_tags records will be automatically deleted.

---

### 6.8 Dashboard & Article List Flow

**Steps:**
1. User navigates to `/dashboard`
2. System verifies authentication
3. Server queries database:
   - Fetch all articles for authenticated user
   - Order by created_at DESC (newest first)
   - Optionally: limit to recent N articles
4. For each article, fetch associated tags
5. Fetch all unique tags for this user (for filter section)
6. Render dashboard template with:
   - Articles array
   - Tags array (for filter)
   - Search query (if coming from search)
   - Active tag filter (if any)

**Optimization:**
- Use JOIN queries to fetch articles with tags in one query
- Cache user's tags list

---

### 6.9 Search Functionality

**Steps:**
1. User enters search term in search bar
2. Form submits GET request to `/search?q=searchterm`
3. System verifies authentication
4. Server processes search:
   - Get query parameter from URL
   - Query database for articles where:
     - user_id = authenticated user AND
     - (title LIKE '%searchterm%' OR content LIKE '%searchterm%')
   - Order by relevance or date
5. Fetch tags for matching articles
6. Render dashboard template with:
   - Filtered articles
   - Search query (to display in search bar)
   - Message: "Results for 'searchterm'" or "No results found"

**Advanced Search (Optional):**
- Full-text search using database features
- Search within tags
- Search operators (AND, OR, NOT)

---

### 6.10 Tag Filtering Flow

**Steps:**
1. User clicks on a tag in dashboard or article view
2. Browser navigates to `/tags/:tagname`
3. System verifies authentication
4. Server processes filter:
   - Get tag name from URL parameter
   - Query database for articles where:
     - user_id = authenticated user AND
     - article is associated with tag name (via article_tags join)
5. Render dashboard template with:
   - Filtered articles
   - Active tag highlighted in tag filter section
   - Clear filter option

---

## 7. Data Validation Rules

### 7.1 User Registration

| Field            | Validation Rules                                           |
|------------------|------------------------------------------------------------|
| Username         | Required, 3-50 chars, alphanumeric + underscores, unique   |
| Email            | Required, valid email format, max 100 chars, unique        |
| Password         | Required, minimum 8 chars, max 128 chars                   |
| Confirm Password | Required, must match password field                        |

---

### 7.2 Article Creation/Update

| Field    | Validation Rules                                           |
|----------|------------------------------------------------------------|
| Title    | Required, 1-255 characters, no leading/trailing whitespace |
| Content  | Required, minimum 1 character, max 100,000 characters      |
| Tags     | Optional, comma-separated, each tag max 50 chars           |

---

### 7.3 Tag Format

- Tags are comma-separated values
- Each tag is trimmed of whitespace
- Tags are converted to lowercase for consistency
- Empty tags are ignored
- Duplicate tags in input are removed
- Maximum 10 tags per article (optional limit)

---

## 8. Security Requirements

### 8.1 Password Security

**Hashing:**
- Use `bcrypt` algorithm
- Cost factor: 10-12 (balance security and performance)
- Never store plain text passwords
- Never log passwords

**Password Requirements:**
- Minimum length: 8 characters
- Recommended: require mix of letters and numbers (optional)
- No maximum length restriction (bcrypt handles this)

---

### 8.2 SQL Injection Prevention

**Required Practices:**
- ALWAYS use parameterized queries (prepared statements)
- NEVER concatenate user input into SQL queries
- Use database driver's parameter binding

**Example Pattern (Conceptual):**
```
BAD:  query := "SELECT * FROM users WHERE username = '" + username + "'"
GOOD: query := "SELECT * FROM users WHERE username = ?"
      db.Query(query, username)
```

---

### 8.3 XSS Prevention

**Template Rendering:**
- Use Go's `html/template` package (auto-escapes by default)
- Never use `template.HTML()` on user input
- Sanitize user input on storage (optional defense in depth)

**Content Security:**
- Escape HTML entities in user-generated content
- Be especially careful with article content and titles

---

### 8.4 CSRF Protection

**Implementation:**
- Generate CSRF token per session
- Include token in all forms as hidden input
- Validate token on all POST requests
- Use `gorilla/csrf` package (recommended)

**Form Pattern:**
```
Every form should include:
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
```

---

### 8.5 Session Security

**Cookie Settings:**
- `HttpOnly`: true (prevent JavaScript access)
- `Secure`: true (HTTPS only - in production)
- `SameSite`: "Lax" or "Strict"
- `MaxAge`: 86400 (24 hours)

**Session Data:**
- Store minimal data in session (user ID, username)
- Regenerate session ID on login
- Invalidate session on logout
- Implement session timeout

---

### 8.6 Authorization

**Article Access Control:**
- ALWAYS verify article ownership before:
  - Viewing article
  - Editing article
  - Deleting article
- Check: `article.user_id == authenticated_user_id`

**Pattern:**
```
1. Fetch article by ID
2. Check if article.user_id == session.user_id
3. If no: return 403 Forbidden
4. If yes: proceed with operation
```

---

### 8.7 Input Sanitization

**General Rules:**
- Trim whitespace from all string inputs
- Validate length limits
- Check for empty/null values
- Validate format (email, etc.)

**Article Content:**
- Allow reasonable HTML formatting (optional)
- If allowing HTML: use allowlist-based sanitizer
- Default: treat as plain text

---

## 9. Error Handling

### 9.1 HTTP Status Codes

| Code | Usage                                                      |
|------|------------------------------------------------------------|
| 200  | Successful GET request                                     |
| 302  | Redirect after successful POST                             |
| 400  | Bad request (validation errors)                            |
| 403  | Forbidden (authorization failed)                           |
| 404  | Resource not found                                         |
| 500  | Internal server error                                      |

---

### 9.2 User-Facing Error Messages

**Display Errors For:**
- Form validation failures
- Login failures
- Resource not found
- Permission denied
- Server errors (generic message)

**Error Message Guidelines:**
- Be specific for validation errors
- Be generic for security-sensitive errors (login)
- Be helpful with recovery suggestions
- Don't expose technical details to users

---

### 9.3 Flash Messages

**Implementation:**
- Store messages in session
- Display once on next page load
- Clear after display

**Message Types:**
- Success (green): "Article created successfully!"
- Error (red): "Failed to update article"
- Warning (yellow): "Your session will expire soon"
- Info (blue): "No articles found"

---

## 10. Project Structure

### 10.1 Recommended Directory Layout
```
knowledge-base/
│
├── cmd/
│   └── web/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── models/                  # Database models and queries
│   │   ├── users.go
│   │   ├── articles.go
│   │   └── tags.go
│   │
│   ├── handlers/                # HTTP handlers
│   │   ├── auth.go             # Login, register, logout
│   │   ├── articles.go         # Article CRUD
│   │   └── home.go             # Home page
│   │
│   └── middleware/              # HTTP middleware
│       ├── auth.go             # Authentication check
│       └── logging.go          # Request logging
│
├── templates/                   # HTML templates
│   ├── layout/
│   │   └── base.html
│   ├── auth/
│   │   ├── login.html
│   │   └── register.html
│   ├── articles/
│   │   ├── dashboard.html
│   │   ├── new.html
│   │   ├── view.html
│   │   └── edit.html
│   └── errors/
│       ├── 404.html
│       └── 403.html
│
├── static/                      # Static files
│   ├── css/
│   │   └── style.css
│   └── images/
│       └── logo.png
│
├── migrations/                  # Database migrations
│   └── 001_init.sql
│
├── go.mod                       # Go modules file
├── go.sum                       # Go modules checksums
├── .env                         # Environment variables (gitignored)
└── README.md                    # Project documentation
```

---

### 10.2 Configuration

**Environment Variables:**
```
DATABASE_URL=knowledge_base.db
SESSION_SECRET=your-secret-key-here-change-in-production
SERVER_PORT=8080
ENVIRONMENT=development
```

**Config File (optional):**
- Database connection settings
- Server port
- Session settings
- Security settings

---

## 11. Development Workflow

### 11.1 Setup Steps

1. **Initialize Project**
   - Create project directory
   - Initialize Go module: `go mod init knowledge-base`
   - Create directory structure

2. **Install Dependencies**
   - gorilla/mux or chi router
   - database driver (sqlite3 or postgresql)
   - gorilla/sessions
   - bcrypt

3. **Database Setup**
   - Create database
   - Run migration SQL to create tables
   - Create indexes

4. **Environment Configuration**
   - Create `.env` file
   - Set database path/URL
   - Generate session secret key

---

### 11.2 Build Order (Recommended)

**Phase 1: Foundation**
1. Set up project structure
2. Configure database connection
3. Create database schema and migrations
4. Set up basic routing and server

**Phase 2: Authentication**
5. Create user model with database methods
6. Implement registration handler and template
7. Implement login handler and template
8. Implement session management
9. Create authentication middleware
10. Implement logout handler

**Phase 3: Core Functionality**
11. Create article model with database methods
12. Implement dashboard (list articles)
13. Implement create article
14. Implement view article
15. Implement edit article
16. Implement delete article

**Phase 4: Tags & Search**
17. Create tag model and relationships
18. Add tag functionality to create/edit articles
19. Implement tag filtering
20. Implement search functionality

**Phase 5: Polish**
21. Improve CSS styling
22. Add error handling and validation
23. Implement flash messages
24. Add 404/403 error pages
25. Testing and bug fixes

---

### 11.3 Testing Checklist

**User Authentication:**
- [ ] Register with valid data
- [ ] Register with duplicate username
- [ ] Register with duplicate email
- [ ] Register with mismatched passwords
- [ ] Login with valid credentials
- [ ] Login with invalid credentials
- [ ] Logout successfully
- [ ] Access protected routes without login

**Article Management:**
- [ ] Create article with valid data
- [ ] Create article with empty title
- [ ] Create article with empty content
- [ ] View own article
- [ ] Edit own article
- [ ] Delete own article
- [ ] Try to view another user's article (should fail)
- [ ] Try to edit another user's article (should fail)
- [ ] Try to delete another user's article (should fail)

**Tags:**
- [ ] Create article with tags
- [ ] Create article without tags
- [ ] Edit article and add tags
- [ ] Edit article and remove tags
- [ ] Filter articles by tag
- [ ] Create multiple articles with same tag

**Search:**
- [ ] Search by title keyword
- [ ] Search by content keyword
- [ ] Search with no results
- [ ] Search with special characters

**Security:**
- [ ] Passwords are hashed in database
- [ ] Sessions persist across requests
- [ ] Sessions expire after logout
- [ ] CSRF protection works (if implemented)
- [ ] SQL injection attempts are prevented
- [ ] XSS attempts are escaped

---

## 12. Optional Enhancements

### Priority 1 (High Value)
- **Markdown Support**: Allow users to write articles in Markdown
  - Use library: `github.com/russross/blackfriday`
  - Syntax highlighting for code blocks
- **Article Categories**: Add a category field (dropdown: Tutorial, Command, Configuration, etc.)
- **Pagination**: Limit articles per page, add next/previous buttons

### Priority 2 (Medium Value)
- **Export Feature**: Download articles as .md or .txt files
- **Dark Mode**: CSS theme toggle stored in session/cookie
- **Rich Text Editor**: Simple formatting toolbar for article content
- **Article Stats**: Word count, character count, reading time estimate

### Priority 3 (Nice to Have)
- **Article Versioning**: Keep history of article edits
- **Favorite/Pin Articles**: Mark important articles to appear at top
- **Share Links**: Generate public read-only links for specific articles
- **Import Feature**: Upload markdown files to create articles
- **Tags Autocomplete**: Suggest existing tags while typing
- **Recent Activity**: Show recently viewed/edited articles

---

## 13. Performance Considerations

### Database Optimization
- Create indexes on frequently queried columns
- Use database connection pooling
- Implement prepared statement caching

### Template Optimization
- Parse templates once at startup
- Cache parsed templates
- Use template inheritance efficiently

### Query Optimization
- Fetch only needed columns (SELECT specific fields)
- Use JOINs to reduce query count
- Add LIMIT clauses where appropriate
- Consider pagination for large article lists

---

## 14. Deployment Considerations

### Production Checklist
- [ ] Use HTTPS (TLS certificates)
- [ ] Set secure session cookie flags
- [ ] Use environment variables for secrets
- [ ] Set proper CORS headers
- [ ] Enable request logging
- [ ] Set up database backups
- [ ] Use production database (PostgreSQL recommended)
- [ ] Compile binary with optimizations
- [ ] Set up process manager (systemd, supervisor)
- [ ] Configure reverse proxy (nginx, caddy)
- [ ] Set rate limiting on login endpoint
- [ ] Monitor disk space (for SQLite) or database connections (for PostgreSQL)

---

## 15. Useful SQL Queries

### Fetch Articles with Tags
```sql
SELECT 
    a.id, a.title, a.content, a.created_at, a.updated_at,
    GROUP_CONCAT(t.name) as tags
FROM articles a
LEFT JOIN article_tags at ON a.id = at.article_id
LEFT JOIN tags t ON at.tag_id = t.id
WHERE a.user_id = ?
GROUP BY a.id
ORDER BY a.created_at DESC
```

### Search Articles
```sql
SELECT * FROM articles
WHERE user_id = ?
AND (title LIKE ? OR content LIKE ?)
ORDER BY created_at DESC
```

### Filter Articles by Tag
```sql
SELECT a.* FROM articles a
INNER JOIN article_tags at ON a.id = at.article_id
INNER JOIN tags t ON at.tag_id = t.id
WHERE a.user_id = ? AND t.name = ?
ORDER BY a.created_at DESC
```

### Get All Tags for User
```sql
SELECT DISTINCT name FROM tags
WHERE user_id = ?
ORDER BY name ASC
```

### Get or Create Tag
```sql
-- Check if tag exists
SELECT id FROM tags
WHERE user_id = ? AND name = ?

-- If not exists, insert
INSERT INTO tags (user_id, name, created_at)
VALUES (?, ?, CURRENT_TIMESTAMP)
```

---

## 16. Common Pitfalls to Avoid

1. **Not checking article ownership** before edit/delete operations
2. **Storing passwords in plain text** (always use bcrypt)
3. **Concatenating SQL queries** instead of using prepared statements
4. **Not handling database errors** gracefully
5. **Exposing sensitive errors** to users (stack traces, SQL errors)
6. **Not validating user input** on the server side
7. **Forgetting to close database connections** (use defer)
8. **Not using HTTPS** in production
9. **Hardcoding secrets** in source code
10. **Not testing authorization** edge cases

---
