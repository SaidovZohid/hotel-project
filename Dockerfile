# Multistage dockerfile (it is for minimizing docker image size)

# -- Build stage -- 
FROM golang:1.19.1-alpine3.16 as builder

#workdir is the current working directory inside docker image 
# all dockerfile instractions will be executed inside workdir 
WORKDIR /hotel

# first dote means that copy everything from current folder (blog_db folder)
# second dot is the current working directory inside the image (/blog folder)
COPY . .

# to install migrate
RUN apk add curl
RUN go build -o main cmd/main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# -- Run stage -- 
FROM alpine:3.16

WORKDIR /hotel
RUN mkdir media

COPY --from=builder /hotel/main .
COPY --from=builder /hotel/migrate ./migrate
COPY migrations ./migrations
COPY templates ./templates

EXPOSE 8080

# running built stage
CMD [ "/hotel/main" ]

# building image and run image 
# docker build -t hotel-app:1.0.0 .
# docker run --name hotel-app -p 8080:8080 hotel-app:1.0.0

# docker network create "name-network"
# docker network connect "name-network " container-name (if it is more than 2, ishould write all containers like this)

# docker run --env-file ./.env --name blog-app --network "network-name" -p 8080:8080 -d blog:1.0.0  
