FROM node:20-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

EXPOSE 3000

# docker run -v ${PWD}:/app -w /app -p 3000:3000 fewt npm run dev
CMD ["npm", "run", "dev"]