# 멀티스테이지 빌드를 위한 Go 빌드 스테이지
FROM golang:1.21-alpine AS builder

# 작업 디렉토리 설정
WORKDIR /app

# Go 모듈 파일 복사 및 의존성 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 바이너리 빌드 (정적 링크)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dhlottery-bot .

# 실행 스테이지 - Ubuntu 기반 (Playwright 브라우저 지원)
FROM ubuntu:22.04

# 시스템 업데이트 및 필수 패키지 설치
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    wget \
    gnupg \
    lsb-release \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# 시간대 설정 (서울)
ENV TZ=Asia/Seoul
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Node.js 설치 (Playwright 의존성)
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - \
    && apt-get install -y nodejs

# Playwright 브라우저 의존성 설치
RUN apt-get update && apt-get install -y \
    libnss3 \
    libnspr4 \
    libatk-bridge2.0-0 \
    libdrm2 \
    libxkbcommon0 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    libgbm1 \
    libxss1 \
    libasound2 \
    libatspi2.0-0 \
    libgtk-3-0 \
    libgdk-pixbuf2.0-0 \
    && rm -rf /var/lib/apt/lists/*

# 작업 디렉토리 설정
WORKDIR /app

# 빌드된 바이너리 복사
COPY --from=builder /app/dhlottery-bot .

# 애플리케이션 사용자 생성 (홈 디렉토리 포함)
RUN useradd -m -s /bin/bash lottouser

# Playwright 브라우저 설치 (root 권한으로)
RUN npx playwright install chromium
RUN npx playwright install-deps chromium

# 스크린샷 디렉토리 생성
RUN mkdir -p screenshots

# Playwright 브라우저 디렉토리 생성 및 권한 설정
RUN mkdir -p /ms-playwright
RUN chown -R lottouser:lottouser /ms-playwright

# 권한 설정
RUN chown -R lottouser:lottouser /app
RUN chown -R lottouser:lottouser /home/lottouser

# 사용자 전환
USER lottouser

# Playwright 환경 변수 설정
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

# 포트 노출 (필요시)
# EXPOSE 8080

# 헬스체크 추가
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pgrep -f dhlottery-bot || exit 1

# 애플리케이션 실행
CMD ["./dhlottery-bot"] 