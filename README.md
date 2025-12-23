🚀 Cloud Task Manager (Full-Stack Go)
A high-performance, concurrent task management system built with Go, backed by PostgreSQL, and deployed on AWS EC2.

🌟 Features
Dual-Mode Interface: Manage tasks via a Terminal CLI or a modern Web Dashboard.

Persistent Cloud Storage: Integration with a managed PostgreSQL (Neon) database.

Concurrent Backend: Built with Go's high-performance standard library (no heavy frameworks).

Cloud Native: Cross-compiled for Linux and deployed on AWS infrastructure.

Automated Background Processing: Runs autonomously on EC2 using nohup.

🛠️ Tech Stack
Language: Go (Golang)

Database: PostgreSQL

Frontend: HTML5, CSS3, JavaScript (Vanilla)

Infrastructure: AWS EC2 (Amazon Linux)

Version Control: Git & GitHub


🚀 Deployment & Installation
Local Setup
Clone the repo: git clone https://github.com/Nelfander/go-cloud-tasks.git

Create a .env file with your DATABASE_URL.

Run the app: go run .

AWS Deployment Strategy
The app is cross-compiled for Linux using:

PowerShell

$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o myapp-linux
It is deployed to EC2 and kept alive as a background process to ensure 100% uptime.

📸 Preview
Live Demo: [http://16.171.16.175/]


Development Milestones :
## 🚀 Recent Milestones

### 🛠️ Infrastructure & DevOps
* **Fixed Critical Server Outage:** Diagnosed a "Disk Pressure" event on AWS (100% disk usage). Used `lsof` to identify "zombie" processes holding deleted 6.4GB log files and reclaimed system stability.
* **Containerized Go Backend:** Architected a `Dockerfile` to package the application, ensuring environment consistency between local development and production.
* **Streamlined Deployment Workflow:** Implemented a Docker Hub "Push/Pull" pipeline to bridge the local environment with AWS EC2.
* **Networking & Security:** Configured Docker Port Mapping (`80:8080`) and optimized AWS Security Groups to allow public HTTP traffic on Port 80.



📈 Roadmap
[1] Dockerization: Containerize the app for easier scaling.

[2] Authentication: Add JWT-based user login.

[3] Custom Domain: Link to a professional .com or .dev address.
