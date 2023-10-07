
# Use an official Go runtime as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 7069 for your application (adjust as needed)
EXPOSE 7069

# Define environment variables
ENV MC_PASS=
ENV MC_ADDRPORT=

# Run the Go application
CMD ["./main"]
