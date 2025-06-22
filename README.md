# 🎰 동행복권 로또 자동 구매 봇

동행복권 사이트에서 로또 6/45를 자동으로 구매하는 Go 언어 기반 봇입니다.

## ✨ 주요 기능

- 🔐 동행복권 사이트 자동 로그인
- 🎲 로또 6/45 자동번호발급으로 5장 구매
- 📸 각 과정별 스크린샷 자동 저장
- ⏰ 매주 토요일 오후 5시 30분 자동 실행 (크론 스케줄)
- 🤖 브라우저 자동화 탐지 우회
- 🔄 사람처럼 자연스러운 동작 구현

## 🎯 구매 프로세스

1. **로그인**: https://dhlottery.co.kr/user.do?method=login&returnUrl=
2. **게임 페이지 이동**: https://el.dhlottery.co.kr/game/TotalGame.jsp?LottoId=LO40
3. **자동번호발급 선택**: 자동번호발급 탭 클릭
4. **구매 수량 설정**: 5장으로 설정
5. **번호 확인**: 확인 버튼 클릭
6. **구매 실행**: 구매하기 버튼 클릭
7. **구매 확정**: 팝업 확인 버튼 클릭

## 🛠️ 설치 및 설정

### 1. 저장소 클론
```bash
git clone <repository-url>
cd bokKwon
```

### 2. Go 모듈 설치
```bash
go mod download
```

### 3. 환경 변수 설정
`.env.example` 파일을 `.env`로 복사하고 실제 계정 정보를 입력하세요:

```bash
cp .env.example .env
```

`.env` 파일 내용:
```env
# 동행복권 계정 정보
DHLOTTERY_ID=your_dhlottery_id
DHLOTTERY_PW=your_dhlottery_password
```

### 4. Playwright 브라우저 설치
```bash
# Playwright 브라우저 드라이버는 첫 실행 시 자동으로 설치됩니다
```

## 🚀 사용법

### 수동 테스트 실행
```bash
go run . test
```

### 자동 스케줄 실행
```bash
go run .
```
또는
```bash
./run.sh
```

## 📸 스크린샷

실행 중 각 단계별로 스크린샷이 `screenshots/` 폴더에 자동 저장됩니다:

- `01_login_page_[timestamp].png` - 로그인 페이지
- `02_login_filled_[timestamp].png` - 로그인 정보 입력 완료
- `03_login_complete_[timestamp].png` - 로그인 완료
- `04_lotto_page_[timestamp].png` - 로또 게임 페이지
- `05_auto_tab_clicked_[timestamp].png` - 자동번호발급 탭 선택
- `06_quantity_set_[timestamp].png` - 구매 수량 설정
- `07_confirm_clicked_[timestamp].png` - 확인 버튼 클릭
- `08_buy_clicked_[timestamp].png` - 구매하기 버튼 클릭
- `09_purchase_complete_[timestamp].png` - 구매 완료

## ⏰ 스케줄

- **자동 실행**: 매주 토요일 오후 5시 30분 (로또 마감 전 여유시간 확보)
- **시간대**: 한국 표준시 (Asia/Seoul)

## 🐳 Docker 실행

### Docker Compose 사용
```bash
docker-compose up -d
```

### 일반 Docker 사용
```bash
docker build -t dhlottery-bot .
docker run -d --name dhlottery-bot \
  -e DHLOTTERY_ID=your_id \
  -e DHLOTTERY_PW=your_password \
  -v $(pwd)/screenshots:/app/screenshots \
  dhlottery-bot
```

## 📋 시스템 요구사항

- Go 1.19 이상
- Chrome/Chromium 브라우저 (자동 설치)
- 네트워크 연결
- 동행복권 계정

## ⚠️ 주의사항

1. **계정 보안**: `.env` 파일을 버전 관리에 포함하지 마세요
2. **이용 약관**: 동행복권 이용약관을 준수하여 사용하세요
3. **구매 한도**: 개인별 구매 한도를 확인하고 사용하세요
4. **네트워크**: 안정적인 인터넷 연결이 필요합니다
5. **책임**: 자동 구매 결과에 대한 책임은 사용자에게 있습니다

## 🔧 문제 해결

### 로그인 실패
- 계정 정보 확인
- 동행복권 사이트 점검 상태 확인
- VPN/프록시 사용 시 해제

### 구매 실패
- 계좌 잔액 확인
- 구매 가능 시간 확인 (마감시간 전)
- 사이트 점검 상태 확인

### 스크린샷 확인
실행 과정에서 문제가 발생하면 `screenshots/` 폴더의 스크린샷을 확인하여 어느 단계에서 실패했는지 파악할 수 있습니다.

## 📝 라이선스

MIT License

## 🤝 기여하기

버그 리포트나 기능 제안은 이슈로 등록해 주세요.

---

**⚖️ 법적 고지**: 이 봇은 교육 및 개인 사용 목적으로 제작되었습니다. 상업적 이용이나 대량 구매는 동행복권 이용약관에 위배될 수 있습니다. 사용자는 관련 법규와 이용약관을 준수할 책임이 있습니다. 