package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func testManualLotto() {
	fmt.Println("🎰 동행복권 로또 자동 구매 봇 수동 테스트")
	fmt.Println("===========================================")

	// 환경 변수 로드
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using system environment variables: %v", err)
	}

	// 필수 환경 변수 확인
	dhlotteryID := os.Getenv("DHLOTTERY_ID")
	dhlotteryPW := os.Getenv("DHLOTTERY_PW")

	fmt.Printf("📝 동행복권 ID: %s\n", dhlotteryID)
	if dhlotteryPW != "" {
		fmt.Printf("🔒 비밀번호: %s\n", "****")
	}

	if dhlotteryID == "" || dhlotteryPW == "" {
		fmt.Println("❌ 환경 변수 설정이 필요합니다:")
		fmt.Println("   DHLOTTERY_ID: 동행복권 아이디")
		fmt.Println("   DHLOTTERY_PW: 동행복권 비밀번호")
		fmt.Println("")
		fmt.Println("💡 .env 파일을 생성하거나 환경 변수를 설정하세요.")
		return
	}

	// 서울 시간대 설정
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		location = time.UTC
	}

	fmt.Printf("⏰ 현재 서울 시간: %s\n", time.Now().In(location).Format("2006-01-02 15:04:05"))
	fmt.Println("")

	fmt.Println("🚀 로또 자동 구매 봇을 실행합니다...")
	fmt.Println("📸 스크린샷이 screenshots/ 폴더에 저장됩니다.")
	fmt.Println("")

	// 로또 구매 봇 실행
	err = runLottoBuyBot()
	if err != nil {
		fmt.Printf("❌ 로또 구매 중 오류 발생: %v\n", err)
		return
	}

	fmt.Println("")
	fmt.Println("✅ 로또 구매가 완료되었습니다!")
	fmt.Println("📸 screenshots/ 폴더에서 스크린샷을 확인하세요.")

	// 구매 완료 후 정보 표시
	fmt.Println("")
	fmt.Println("🎫 구매 정보:")
	fmt.Println("   - 게임: 로또 6/45")
	fmt.Println("   - 수량: 5장")
	fmt.Println("   - 방식: 자동번호발급")
	fmt.Printf("   - 구매시간: %s\n", time.Now().In(location).Format("2006-01-02 15:04:05"))
}
