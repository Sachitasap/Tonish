# Frontend Dockerfile
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Final stage
FROM node:20-alpine

WORKDIR /app

# Copy built application and dependencies
COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules

ENV PORT=4173
ENV HOST=0.0.0.0

EXPOSE 4173

CMD ["node", "build"]
