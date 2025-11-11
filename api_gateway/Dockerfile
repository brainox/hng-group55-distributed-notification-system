# FROM node:18-alpine AS builder

# WORKDIR /app

# # Copy package files
# COPY package*.json ./

# # Install dependencies
# RUN npm ci --only=production && npm cache clean --force

# # Copy source code
# COPY . .

# # Build the application
# RUN npm run build

# # Production image
# FROM node:18-alpine

# WORKDIR /app

# # Copy built application
# COPY --from=builder /app/dist ./dist
# COPY --from=builder /app/node_modules ./node_modules
# COPY --from=builder /app/package*.json ./

# EXPOSE 3000

# CMD ["node", "dist/main.js"]



# -------- Builder Stage --------
FROM node:18-alpine AS builder
WORKDIR /app

# Copy dependency files
COPY package*.json ./

# Install ALL dependencies (including devDeps)
RUN npm ci && npm cache clean --force

# Copy the source code
COPY . .

# Build the application
RUN npm run build

# -------- Production Stage --------
FROM node:18-alpine
WORKDIR /app

ENV NODE_ENV=production

# Copy only package files for production deps
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force

# Copy built app
COPY --from=builder /app/dist ./dist

EXPOSE 3000
CMD ["node", "dist/main.js"]
