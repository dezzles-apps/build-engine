FROM node:latest AS production-stage
WORKDIR /
COPY ./src /src
COPY package*.json ./
RUN npm install --only=production
EXPOSE 8888
ENTRYPOINT [ "node", "./src/index.js" ]