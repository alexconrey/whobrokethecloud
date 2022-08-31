FROM node:18-alpine as base
RUN apk update

FROM base as assets

WORKDIR /app
COPY package-lock.json package.json ./

RUN npm install

COPY . .
RUN npm run build

FROM base as app
RUN npm install --global serve
COPY --from=assets /app/build /app

ENV PORT 3000

WORKDIR /app
ENTRYPOINT [ "serve" ]