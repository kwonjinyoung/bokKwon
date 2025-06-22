# 🎲 동행복권 로또 자동 구매 봇

매주 토요일 오후 5시 30분에 자동으로 로또 6/45를 5장씩 구매하는 봇입니다.

## 🚀 빠른 시작 (3단계로 완료!)

### 1단계: 소스 코드 다운로드
```bash
git clone <repository-url>
cd bokKwon
```

### 2단계: 동행복권 계정 정보 설정
```bash
# .env 파일 생성
cp .env.example .env

# .env 파일 편집 (메모장이나 에디터로 열어서 수정)
nano .env
```

`.env` 파일에서 아래 정보를 입력하세요:
```env
# 동행복권 로그인 정보 (필수)
DHLOTTERY_ID=여기에_동행복권_아이디_입력
DHLOTTERY_PW=여기에_동행복권_비밀번호_입력
```

### 3단계: Docker로 실행
```bash
# Docker 실행 스크립트 실행
./docker-run.sh
```

**끝!** 이제 매주 토요일 오후 5시 30분에 자동으로 로또를 구매합니다! 🎉

---

## ✨ 주요 기능

- 🔐 **자동 로그인**: 동행복권 사이트에 자동 로그인
- 🎲 **자동 구매**: 로또 6/45 자동번호 5장 구매
- ⏰ **스케줄 실행**: 매주 토요일 오후 5시 30분 자동 실행
- 📸 **과정 기록**: 모든 단계를 스크린샷으로 저장
- 🤖 **탐지 우회**: 브라우저 자동화 탐지 방지

## 📅 구매 스케줄

- **실행 시간**: 매주 토요일 오후 5시 30분
- **구매 수량**: 자동번호 5장 (5,000원)
- **마감 시간**: 토요일 오후 8시 (여유시간 2시간 30분)

## 🔧 상세 설치 가이드

### 사전 요구사항
- Docker 및 Docker Compose 설치
- 동행복권 계정 (회원가입 필요)
- 계좌에 충분한 잔액 (최소 5,000원 이상)

### Docker 설치 (Ubuntu/Debian)
```bash
# Docker 설치
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Docker Compose 설치
sudo apt-get update
sudo apt-get install docker-compose-plugin

# 현재 사용자를 docker 그룹에 추가
sudo usermod -aG docker $USER
```

### 환경 설정 상세
1. **저장소 클론**:
   ```bash
   git clone <repository-url>
   cd bokKwon
   ```

2. **환경 변수 파일 생성**:
   ```bash
   cp .env.example .env
   ```

3. **동행복권 계정 정보 입력**:
   ```bash
   # 에디터로 .env 파일 열기
   nano .env
   # 또는
   vi .env
   # 또는 Windows에서
   notepad .env
   ```

   **중요**: 실제 동행복권 아이디와 비밀번호를 정확히 입력하세요!

4. **Docker 실행**:
   ```bash
   # 대화형 스크립트 실행
   ./docker-run.sh
   
   # 또는 직접 명령어
   docker-compose up -d --build
   ```

## 📱 실행 방법

### 방법 1: 스크립트 사용 (권장)
```bash
./docker-run.sh
```
스크립트 실행 시 다음 옵션을 선택할 수 있습니다:
1. 백그라운드 실행 (권장) - 자동 스케줄 실행
2. 포그라운드 실행 - 로그 실시간 확인
3. 테스트 모드 - 즉시 구매 테스트
4. 컨테이너 중지
5. 컨테이너 재시작
6. 로그 확인

### 방법 2: 직접 명령어
```bash
# 백그라운드에서 실행 (스케줄 모드)
docker-compose up -d --build

# 즉시 테스트 실행
docker-compose run --rm dhlottery-lotto-bot ./dhlottery-bot test

# 로그 확인
docker-compose logs -f

# 중지
docker-compose down
```

## 📸 구매 과정 확인

실행 시 `screenshots/` 폴더에 각 단계별 스크린샷이 저장됩니다:

1. `01_login_page_[시간].png` - 로그인 페이지
2. `02_login_filled_[시간].png` - 로그인 정보 입력
3. `03_login_complete_[시간].png` - 로그인 완료
4. `04_lotto_page_[시간].png` - 로또 게임 페이지
5. `05_auto_tab_clicked_[시간].png` - 자동번호발급 선택
6. `06_quantity_set_[시간].png` - 구매 수량 설정 (5장)
7. `07_confirm_clicked_[시간].png` - 확인 버튼 클릭
8. `08_buy_clicked_[시간].png` - 구매하기 버튼 클릭
9. `09_purchase_complete_[시간].png` - 구매 완료

## 🔍 상태 확인

### 컨테이너 상태 확인
```bash
docker-compose ps
```

### 로그 확인
```bash
# 실시간 로그
docker-compose logs -f

# 최근 로그만
docker-compose logs --tail=50

# 에러 로그만
docker-compose logs | grep -i error
```

### 다음 실행 시간 확인
로그에서 다음과 같은 메시지를 확인할 수 있습니다:
```
동행복권 로또 자동 구매 봇이 시작되었습니다.
스케줄: 서울 시간 기준 매주 토요일 오후 5시 30분
현재 서울 시간: 2024-XX-XX XX:XX:XX
```

## 🛠️ 문제 해결

### 자주 발생하는 문제

#### 1. 로그인 실패
**증상**: "로그인 실패" 메시지
**해결책**:
- `.env` 파일의 아이디/비밀번호 확인
- 동행복권 사이트에서 직접 로그인 테스트
- 계정 잠금 상태 확인

#### 2. 구매 실패
**증상**: "구매 실패" 메시지
**해결책**:
- 계좌 잔액 확인 (최소 5,000원 필요)
- 구매 가능 시간 확인 (토요일 오후 8시 마감)
- 동행복권 사이트 점검 상태 확인

#### 3. Docker 실행 오류
**증상**: "Permission denied" 또는 "Docker not found"
**해결책**:
```bash
# Docker 서비스 시작
sudo systemctl start docker

# 권한 문제 해결
sudo usermod -aG docker $USER
# 로그아웃 후 다시 로그인

# Docker Compose 설치 확인
docker-compose --version
```

#### 4. 환경 변수 오류
**증상**: "환경 변수를 찾을 수 없습니다"
**해결책**:
```bash
# .env 파일 존재 확인
ls -la .env

# .env 파일 내용 확인
cat .env

# 파일이 없으면 다시 생성
cp .env.example .env
nano .env
```

### 디버깅 방법

#### 테스트 모드로 즉시 실행
```bash
# 스케줄을 기다리지 않고 즉시 테스트
docker-compose run --rm dhlottery-lotto-bot ./dhlottery-bot test
```

#### 스크린샷 확인
구매 과정에서 문제가 발생하면 `screenshots/` 폴더의 스크린샷을 확인하여 어느 단계에서 실패했는지 파악할 수 있습니다.

## ⚠️ 중요 주의사항

### 보안
- **절대로** `.env` 파일을 다른 사람과 공유하지 마세요
- GitHub 등에 업로드할 때 `.env` 파일은 제외하세요
- 동행복권 계정 정보를 안전하게 관리하세요

### 사용 제한
- 개인 사용 목적으로만 사용하세요
- 상업적 이용은 금지됩니다
- 동행복권 이용약관을 준수하세요

### 구매 관련
- 매주 5,000원씩 자동 결제됩니다
- 계좌 잔액을 항상 확인하세요
- 구매 결과에 대한 책임은 사용자에게 있습니다

## 📞 지원

### 로그 수집
문제 발생 시 다음 정보를 수집해 주세요:
```bash
# 시스템 정보
uname -a
docker --version
docker-compose --version

# 컨테이너 상태
docker-compose ps

# 최근 로그
docker-compose logs --tail=100 > logs.txt
```

### 이슈 리포트
버그나 문제 발생 시 GitHub Issues에 다음 정보와 함께 등록해 주세요:
- 실행 환경 (OS, Docker 버전)
- 에러 메시지
- 로그 파일
- 스크린샷

---

## 📝 라이선스

MIT License

**⚖️ 법적 고지**: 이 봇은 교육 및 개인 사용 목적으로 제작되었습니다. 상업적 이용이나 대량 구매는 동행복권 이용약관에 위배될 수 있습니다. 사용자는 관련 법규와 이용약관을 준수할 책임이 있습니다. 