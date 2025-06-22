# 🎲 동행복권 로또 자동 구매 봇 Docker 가이드

이 가이드는 동행복권 로또 자동 구매 봇을 Docker 환경에서 실행하는 방법을 설명합니다.

## 📋 사전 요구사항

- Docker 및 Docker Compose 설치
- Ubuntu 20.04 이상 권장
- 최소 1GB RAM, 2GB 디스크 공간
- 동행복권 계정

## 🚀 빠른 시작

### 1. 환경 설정

```bash
# .env 파일 생성
cp .env.example .env

# .env 파일 편집 (동행복권 로그인 정보 입력)
nano .env
```

`.env` 파일 예시:
```
# 동행복권 로그인 정보 (필수)
DHLOTTERY_ID=your_actual_id
DHLOTTERY_PW=your_actual_password

# 로또 구매 설정
# 매주 토요일 오후 5시 30분에 자동 실행됩니다
# 구매 수량: 5장 (자동번호)
```

### 2. Docker 실행

#### 방법 1: 스크립트 사용 (권장)
```bash
./docker-run.sh
```

#### 방법 2: 직접 명령어 사용
```bash
# 백그라운드 실행
docker-compose up -d --build

# 포그라운드 실행 (로그 확인)
docker-compose up --build

# 테스트 모드 실행
docker-compose run --rm dhlottery-lotto-bot ./dhlottery-bot test
```

## 📊 컨테이너 관리

### 상태 확인
```bash
docker-compose ps
```

### 로그 확인
```bash
# 실시간 로그 확인
docker-compose logs -f

# 최근 로그만 확인
docker-compose logs --tail=100
```

### 컨테이너 중지/재시작
```bash
# 중지
docker-compose down

# 재시작
docker-compose restart

# 강제 재빌드 후 실행
docker-compose up -d --build --force-recreate
```

## 📁 파일 구조

```
dhlottery-lotto-bot/
├── Dockerfile              # Docker 이미지 빌드 설정
├── docker-compose.yml      # Docker Compose 설정
├── .dockerignore           # Docker 빌드 시 제외할 파일
├── .env.example            # 환경 변수 예시
├── docker-run.sh           # Docker 실행 스크립트
├── screenshots/            # 스크린샷 저장 폴더 (자동 생성)
└── logs/                   # 로그 저장 폴더 (자동 생성)
```

## 📅 스케줄링

- **자동 실행**: 매주 토요일 오후 5시 30분 (서울 시간 기준)
- **로또 마감**: 토요일 오후 8시까지이므로, 여유있게 5시 30분에 구매
- **구매 수량**: 자동번호 5장
- **스크린샷**: 각 단계별로 자동 저장

## 🔧 고급 설정

### 리소스 제한 조정

`docker-compose.yml`에서 리소스 제한을 조정할 수 있습니다:

```yaml
deploy:
  resources:
    limits:
      memory: 1G        # 메모리 제한
      cpus: '0.5'       # CPU 제한
    reservations:
      memory: 512M      # 최소 메모리
      cpus: '0.25'      # 최소 CPU
```

### 로그 설정 변경

```yaml
logging:
  driver: "json-file"
  options:
    max-size: "10m"     # 로그 파일 최대 크기
    max-file: "3"       # 보관할 로그 파일 수
```

### 시간대 변경

컨테이너는 기본적으로 Asia/Seoul로 설정되어 있습니다.

## 🐛 문제 해결

### 1. 브라우저 실행 오류
```bash
# 컨테이너 내부에서 디버깅
docker-compose exec dhlottery-lotto-bot /bin/bash

# 또는 새 컨테이너로 디버깅
docker-compose run --rm dhlottery-lotto-bot /bin/bash
```

### 2. 권한 오류
```bash
# 호스트에서 권한 수정
sudo chown -R $USER:$USER screenshots/
sudo chown -R $USER:$USER logs/
```

### 3. 메모리 부족
```bash
# 시스템 메모리 확인
free -h

# Docker 메모리 사용량 확인
docker stats
```

### 4. 네트워크 연결 문제
```bash
# 컨테이너 네트워크 확인
docker network ls
docker-compose exec dhlottery-lotto-bot ping dhlottery.co.kr
```

### 5. 로그인 실패
- 동행복권 계정 정보가 올바른지 확인
- 계정이 잠겨있지 않은지 확인
- 2단계 인증이 설정되어 있지 않은지 확인

## 📈 모니터링

### 컨테이너 상태 모니터링
```bash
# 실시간 리소스 사용량
docker stats dhlottery-auto-lotto-bot

# 헬스체크 상태 확인
docker inspect dhlottery-auto-lotto-bot | grep -A 10 Health
```

### 로그 분석
```bash
# 에러 로그만 필터링
docker-compose logs | grep -i error

# 성공 로그 확인
docker-compose logs | grep -i "로또 구매 완료"

# 특정 시간대 로그 확인
docker-compose logs --since="2024-01-01T09:00:00"
```

## 🔄 업데이트

### 코드 업데이트 후 재배포
```bash
# 이미지 재빌드 및 재시작
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### 의존성 업데이트
```bash
# Go 모듈 업데이트 후 재빌드
go mod tidy
docker-compose build --no-cache
```

## 🛡️ 보안 고려사항

1. **환경 변수 보안**: `.env` 파일을 Git에 커밋하지 마세요
2. **계정 보안**: 동행복권 계정 정보를 안전하게 관리하세요
3. **네트워크 보안**: 필요시 방화벽 설정을 확인하세요

## ⚖️ 법적 고지

이 봇은 교육 및 개인 사용 목적으로 제작되었습니다. 상업적 이용이나 대량 구매는 동행복권 이용약관에 위배될 수 있습니다. 사용자는 관련 법규와 이용약관을 준수할 책임이 있습니다.

```bash
# 시스템 정보 수집
docker --version
docker-compose --version
free -h
df -h
docker-compose logs --tail=50
``` 