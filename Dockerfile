FROM golang:latest



RUN wget https://github.com/knative/client/releases/download/knative-v1.11.0/kn-linux-amd64 && mv kn-linux-amd64 /usr/local/bin/kn && chmod +x /usr/local/bin/kn
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 && chmod 700 get_helm.sh && ./get_helm.sh
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && chmod +x kubectl && mv kubectl /usr/local/bin/kubectl
# Set the Current Working Directory inside the container
WORKDIR /go/src/app/

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go mod download -x

# Install compile daemon for hot reloading
RUN go install -mod=mod github.com/githubnemo/CompileDaemon

# Expose port 80 to the outside world
EXPOSE 8282

# Command to run the executable
ENTRYPOINT CompileDaemon --build="go build main.go" --command="./main"
