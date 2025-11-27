# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Instalar yarn
RUN corepack enable

# Copiar arquivos de dependências
COPY package.json yarn.lock ./

# Instalar dependências
RUN yarn install --frozen-lockfile

# Copiar código fonte
COPY . .

# Build
RUN yarn build

# Runtime stage
FROM nginx:alpine

# Copiar configuração do nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copiar arquivos de build
COPY --from=builder /app/dist /usr/share/nginx/html

# Expor porta
EXPOSE 3000

# Executar
CMD ["nginx", "-g", "daemon off;"]

