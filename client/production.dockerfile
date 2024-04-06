FROM node:21-alpine3.18 as build
WORKDIR /client

COPY package*.json ./
RUN npm ci
COPY . ./

RUN npm run build
RUN npm prune --production

FROM node:21-alpine3.18
WORKDIR /client

COPY --from=build /client/build build/
COPY --from=build /client/node_modules node_modules/
COPY --from=build /client/package.json ./

EXPOSE 80
ENV NODE_ENV=production
CMD ["node", "build"]
