#збірка статичних файлів
FROM node:18.8.0-alpine as build
WORKDIR /app
COPY client/package.json ./
COPY client/package-lock.json ./
RUN npm install
COPY client/. ./
RUN npm run build

#запуск серверу
FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
# COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
