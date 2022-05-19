#####################################################
# Dockerfile
#
# Creates an image with the Godot headless app.
#
# Build Args:
#   - GODOT_VERSION: The version of Godot
#   - EXPORT_TEMPLATES: Included export templates
#       examples "all", "none", "win"
#

ARG EXPORT_TEMPLATES=all

#------------------------------
# Alias for the root image
FROM debian:stable-slim AS base

ARG GODOT_VERSION=3.4.2

#------------------------------
# Installs packages to use wget
FROM base as wget

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    unzip

#----------------
# Downloads Godot
FROM wget AS godot

RUN wget https://downloads.tuxfamily.org/godotengine/3.4.4/Godot_v3.4.4-stable_linux_headless.64.zip 
RUN unzip Godot_v3.4.4-stable_linux_headless.64.zip
RUN mv Godot_v3.4.4-stable_linux_headless.64 /usr/local/bin/godot

#----------------
# Build web-server

FROM golang:1.18.2 as web

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o ./server

#------------------------------
# Clean setup with no templates

FROM base AS export-none

ENV XDG_DATA_HOME /usr/local/share

RUN mkdir -p /root/.cache
RUN mkdir -p /root/.config/godot

WORKDIR /app
COPY --from=godot /usr/local/bin/godot /usr/local/bin/godot
COPY --from=web /app/server /usr/local/bin/server
COPY template ./
# EXPOSE 8080
ENV SCRIPTS=/app/scripts

ENTRYPOINT ["server"]
CMD ["--help"]