FROM node:20 AS build

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

COPY . .
RUN yarn build

FROM node:20-alpine AS app

WORKDIR /app

COPY --from=build /app/package.json /app/yarn.lock ./
RUN yarn install --frozen-lockfile --production

COPY --from=build /app/build ./build

RUN adduser -D appuser && chown -R appuser /app
USER appuser

EXPOSE 3000
CMD ["node", "build"]
