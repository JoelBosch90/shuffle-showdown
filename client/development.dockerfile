FROM node:21-alpine3.18
WORKDIR /client

COPY package.json package-lock.json ./
RUN npm install
COPY . .

EXPOSE 80

CMD ["npm", "run", "dev", "--", "--host", "--port", "80"]