FROM node:20-alpine
WORKDIR /app
COPY package.json ./
RUN npm install --omit=dev
COPY src ./src
EXPOSE 3000
# --import registers the OpenTelemetry SDK + auto-instrumentations before the app
# loads, so Express and outgoing fetch calls are traced with zero code changes.
CMD ["node", "--import", "@opentelemetry/auto-instrumentations-node/register", "src/index.js"]
