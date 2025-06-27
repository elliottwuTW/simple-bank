### Build Stage
# Define the base image for the application
FROM golang:1.24-alpine3.22 AS builder
# Declare the current working directory
WORKDIR /app
# First dot: Copy everything from current folder 
# Second dot: the current working directory inside the image
COPY . .
# Build the app to a single binary executable file
RUN go build -o main main.go
# 可以再執行一些像是要下載 binary 的指令

### Run stage
### (We want only the image just with the binary file, and without something like golang code)
FROM alpine:3.22
WORKDIR /app
# Copy the executable binary file from the builder stage to this run stage image
COPY --from=builder /app/main .
# 可以 COPY 像是下載完的 binary 到 /app/xxx(binary檔名) => COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY env.json .

# Best practice to EXPOSE instruction to inform Docker
# the container listens on the specified network port at runtime.
# i.e., the port we declare in the env.json file
# It doesn't actually publish the port, and only functions as
# document for the person about which ports are intended to be published
EXPOSE 8080

# Define the default command to run when the container starts
# => run the binary file
CMD [ "/app/main" ]
