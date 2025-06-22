#!/bin/bash

# Docker 실행 스크립트 for 동행복권 로또 자동 구매 봇

echo "🎲 동행복권 로또 자동 구매 봇 Docker 실행 스크립트"
echo "================================================"

# .env 파일 확인
if [ ! -f .env ]; then
    echo "❌ .env 파일이 없습니다."
    echo "📝 .env.example 파일을 참고하여 .env 파일을 생성해주세요."
    echo ""
    echo "cp .env.example .env"
    echo "nano .env  # 또는 다른 에디터로 편집"
    echo ""
    exit 1
fi

# 필요한 디렉토리 생성
echo "📁 필요한 디렉토리를 생성합니다..."
mkdir -p screenshots
mkdir -p logs

# Docker Compose 실행 방법 선택
echo ""
echo "실행 방법을 선택하세요:"
echo "1) Docker Compose로 백그라운드 실행 (권장)"
echo "2) Docker Compose로 포그라운드 실행 (로그 확인용)"
echo "3) 테스트 모드로 실행"
echo "4) 컨테이너 중지"
echo "5) 컨테이너 재시작"
echo "6) 로그 확인"
echo ""

read -p "선택 (1-6): " choice

case $choice in
    1)
        echo "🚀 백그라운드에서 실행합니다..."
        docker-compose up -d --build
        echo "✅ 컨테이너가 백그라운드에서 실행 중입니다."
        echo "📅 스케줄: 매주 토요일 오후 5시 30분 자동 실행"
        echo "📊 상태 확인: docker-compose ps"
        echo "📋 로그 확인: docker-compose logs -f"
        ;;
    2)
        echo "🚀 포그라운드에서 실행합니다... (Ctrl+C로 중지)"
        docker-compose up --build
        ;;
    3)
        echo "🧪 테스트 모드로 실행합니다..."
        docker-compose run --rm dhlottery-lotto-bot ./dhlottery-bot test
        ;;
    4)
        echo "⏹️ 컨테이너를 중지합니다..."
        docker-compose down
        echo "✅ 컨테이너가 중지되었습니다."
        ;;
    5)
        echo "🔄 컨테이너를 재시작합니다..."
        docker-compose restart
        echo "✅ 컨테이너가 재시작되었습니다."
        ;;
    6)
        echo "📋 로그를 확인합니다... (Ctrl+C로 종료)"
        docker-compose logs -f
        ;;
    *)
        echo "❌ 잘못된 선택입니다."
        exit 1
        ;;
esac

echo ""
echo "🔧 유용한 명령어들:"
echo "- 상태 확인: docker-compose ps"
echo "- 로그 확인: docker-compose logs -f"
echo "- 컨테이너 중지: docker-compose down"
echo "- 컨테이너 재시작: docker-compose restart"
echo "- 이미지 재빌드: docker-compose build --no-cache" 