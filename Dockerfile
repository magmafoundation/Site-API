FROM alpine AS build

RUN apk add --no-cache nodejs nodejs-npm

WORKDIR /app

COPY package.json .
COPY tsconfig.json .
COPY ./src src

RUN npm install
RUN npm run build
RUN ls


FROM alpine

RUN apk add --no-cache nodejs nodejs-npm
WORKDIR /app
COPY --from=build /app/dist/ .
COPY package.json .
COPY tsconfig.json .

RUN npm install
CMD ["npm", "start"]
EXPOSE 8080
