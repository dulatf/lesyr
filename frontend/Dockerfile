FROM node:18-alpine AS builder
WORKDIR /app

# Install pnpm globally
RUN npm install -g pnpm

# Install dependencies with pnpm
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Copy the rest of the app source code
COPY . .

# Build the Next.js application using pnpm
RUN pnpm run build

FROM node:18-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

# Install pnpm globally
RUN npm install -g pnpm

# Copy necessary files from the builder stage
# COPY --from=builder /app/next.config.js ./
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

EXPOSE 3000

# Start the Next.js application using pnpm
CMD ["pnpm", "start"]