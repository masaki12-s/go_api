FROM node:20.12.2

COPY package.json package-lock.json .
WORKDIR /usr/src/app

RUN npm install

COPY . .