FROM node:alpine AS build

RUN npm install -g pnpm

WORKDIR /app
COPY . .

RUN pnpm i
RUN pnpm build

# ------------------------------------------------------------------------------

FROM nginx:alpine AS final

COPY --from=build /app/build/* ./usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]