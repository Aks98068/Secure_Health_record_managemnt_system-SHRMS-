#1. Build stage
FROM golang:1.25-apline AS builder

# Set working dictory
WORKDIR /app

#Install git (if you fetch dependencies)
RUN apk add --no-cache git

# Copy go.sum and go.mod
COPY go.mod go.sum ./

# Download dependencies
RUN go mod Download

#copy the source code 
COPY . .

#Build the binary (disable CGO  for smaller image
RUN go build -o server main.go


# =========================
# 2. Final stages(smaller image)
# =====================

FROM alpibe:3.18

  
#set working directory
WORKDIR /app

#cOPY THE BINARY FROM BUILDER
    COPY --from=builder /app/server  .

# optional: copy TLS certificates it not mounting 
#COPY server.crt server.key ./

#Expose HTTPS PORT 
EXPOSE 8443

# run the server 
CMD ["./server"]