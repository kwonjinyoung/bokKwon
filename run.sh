#!/bin/bash

# Go 경로 설정
export PATH=$PATH:/usr/local/go/bin

echo "🚀 뽐뿌 자동 댓글 봇을 시작합니다..."
echo "📅 스케줄: 오전 9시~오후 9시, 1.5~2시간 간격"
echo "⏰ 현재 시간: $(date)"
echo ""

# 환경 변수 확인
if [ ! -f .env ]; then
    echo "❌ .env 파일이 없습니다. .env 파일을 생성하고 로그인 정보를 입력해주세요."
    echo "PPOMPPU_ID=your_actual_id"
    echo "PPOMPPU_PW=your_actual_password"
    exit 1
fi

# 프로그램 실행
go run main.go 