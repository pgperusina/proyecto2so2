FROM node:10 AS build-env
ADD . /app
WORKDIR /app

RUN npm ci --only=production
RUN npm run build

FROM gcr.io/distroless/nodejs:10
COPY --from=build-env /app /app
WORKDIR /app
CMD ["node_modules/serve/bin/serve.js", "-s", "-d",  "build"]