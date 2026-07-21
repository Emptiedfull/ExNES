FROM tinygo/tinygo:0.39.0 AS build
USER root 
WORKDIR /src
COPY . . 
RUN cp -r cmd/wasm/server/static /build
RUN tinygo build -o /build/nes.wasm -target wasm -opt z cmd/wasm/wasm.go
RUN cp "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" /build/wasm_exec.js
RUN go run ./cmd/wasm/bundler /build 

FROM caddy:alpine 
EXPOSE 5000
COPY --from=build /build /srv 
COPY Caddyfile /etc/caddy/Caddyfile
