# Real Time Voting Application

## Setup Instructions
1. Clone the repo
    git clone <repo_url>
    cd voting

2. Build the docker image
    docker build . -t voting

3. Running the containers
    docker run -p 8080:8080 voting

4. Access the application by navigating to http://localhost:8080.