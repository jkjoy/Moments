# 前端构建阶段
FROM node:22.2.0-alpine AS front
WORKDIR /app
COPY front/package*.json ./
RUN npm install
COPY front/. .
RUN npm run generate

# 后端构建阶段
FROM golang:1.22.5-alpine AS backend
ENV CGO_ENABLED=1
RUN apk add --no-cache build-base
WORKDIR /app
COPY backend/go.mod .
COPY backend/go.sum .
RUN go mod tidy
RUN go mod download
RUN go list -m all  # 检查所有依赖项是否正确
COPY backend/. .
COPY --from=front /app/.output/public /app/public
RUN apk update --no-cache && apk add --no-cache tzdata
RUN go get -d -v ./...  # 获取所有依赖项
RUN set -x && go build -tags prod -ldflags="-s -w" -o /app/moments  # 启用详细日志输出

# 最终运行阶段
FROM alpine
ARG VERSION
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
COPY --from=backend /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai
WORKDIR /app/data
ENV VERSION=$VERSION
ENV FEISHU_WEBHOOK_URL="https://open.feishu.cn/open-apis/bot/v2/hook/"
COPY --from=backend /app/moments /app/moments
ENV PORT=3000
EXPOSE 3000
RUN chmod +x /app/moments
CMD ["/app/moments"]
