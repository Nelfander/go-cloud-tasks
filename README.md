# Go Cloud Task Manager (V2)

A secure, Dockerized To-Do application built with Go, PostgreSQL, and JWT Authentication, deployed on AWS EC2.

## 🚀 Key Achievements
- **Full-Stack Security:** Implemented user registration and login using **Bcrypt** for password hashing and **JWT (JSON Web Tokens)** for secure session management.
- **Cloud Database:** Integrated with **Neon PostgreSQL**, featuring a relational schema with `users` and `tasks` linked via foreign keys.
- **DevOps & Deployment:** - Containerized the application using **Docker** (Multi-stage builds for optimization).
    - Deployed to an **AWS EC2** instance.
    - Implemented secure environment variable injection to keep credentials out of source code.

## 🛠️ Tech Stack
- **Backend:** Go (Golang)
- **Database:** PostgreSQL (Neon)
- **Containerization:** Docker
- **Cloud:** AWS (EC2)
- **Auth:** JWT & Bcrypt

## 🛠️ Lessons Learned & Problem Solving

### 💾 Incident: Disk Space Exhaustion & Ghost Files
During deployment the EC2 instance ran out of disk space causing the application to crash.
- **The Issue:** Deleted files were still occupying space because their file handles were being held open by zombie processes.
- **The Fix:** I used `lsof +L1` to identify the open handles and killed the associated processes to reclaim **6GB** of space instantly.
- **Prevention:** Implemented **Docker Log Rotation** in the configuration to prevent container logs from growing indefinitely in the future.

### 🔐 Incident: Environment Variable Mismatches
Faced a "Restart Loop" in Docker where the container crashed due to a failed DB connection.
- **The Issue:** The Go binary was looking for `DB_URL` while the environment was providing `DATABASE_URL`.
- **The Fix:** Synchronized naming conventions and added `db.Ping()` to the initialization logic to provide clearer error logging during the startup phase.

## ⚙️ Local Development
1. Clone the repo: `git clone https://github.com/Nelfander/go-cloud-tasks.git`
2. Create a `.env` file with `DB_URL` and `JWT_SECRET`.
3. Run `go run .` or use Docker:
   ```bash
   docker build -t task-manager .
   docker run -p 8080:8080 --env-file .env task-manager
