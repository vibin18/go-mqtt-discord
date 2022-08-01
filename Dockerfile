FROM golang:1.18-alpine as build
RUN apk upgrade --no-cache --force
RUN apk add --update build-base make git
WORKDIR /go/src/github.com/vibin18/go-mqtt-discord

# Compile
COPY ./ /go/src/github.com/vibin18/go-mqtt-discord
RUN make dependencies
RUN make build
RUN ./go-mqtt-discord --help

# Final Image
FROM gcr.io/distroless/static AS export-stage
ENV FRIGATE_SERVER="http://192.168.68.126:5000" \
    FRIGATE_MQTT_SERVER="192.168.68.126:1883" \
    DISCORD_TOKEN="YOUR_DISCORD_TOKEN" \
    DISCORD_CHANNEL_ID="YOUR_DISCORD_CHANNEL_ID" \
    SNAPSHOT_PIXEL="400"

COPY --from=build /go/src/github.com/vibin18/go-mqtt-discord/go-mqtt-discord /
USER 1000:1000
ENTRYPOINT ["/go-mqtt-discord"]