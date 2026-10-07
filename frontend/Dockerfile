# syntax=docker/dockerfile:1

FROM node:22-alpine AS build
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
ARG VITE_API_URL=/api/v1
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY nginx-security-headers.conf /etc/nginx/snippets/elabtrack-security-headers.conf
COPY nginx-api-proxy.conf /etc/nginx/snippets/elabtrack-api-proxy.conf
COPY --from=build /src/dist /usr/share/nginx/html
EXPOSE 80
