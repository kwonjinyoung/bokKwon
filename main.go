package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/playwright-community/playwright-go"
	"github.com/robfig/cron/v3"
)

type UserAgent struct {
	UserAgent string
	Viewport  struct {
		Width  int
		Height int
	}
}

const (
	LOGIN_URL = "https://www.dhlottery.co.kr/login"
	LOTTO_URL = "https://el.dhlottery.co.kr/game/TotalGame.jsp?LottoId=LO40"
)

var (
	// 다양한 실제 사용자 에이전트
	userAgents = []UserAgent{
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			Viewport:  struct{ Width, Height int }{1920, 1080},
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			Viewport:  struct{ Width, Height int }{1366, 768},
		},
		{
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			Viewport:  struct{ Width, Height int }{1440, 900},
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
			Viewport:  struct{ Width, Height int }{1920, 1080},
		},
		{
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
			Viewport:  struct{ Width, Height int }{1680, 1050},
		},
	}
)

func main() {
	// 환경 변수 로드 (파일이 없어도 계속 진행)
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using system environment variables: %v", err)
	}

	// 필수 환경 변수 확인
	dhlotteryID := os.Getenv("DHLOTTERY_ID")
	dhlotteryPW := os.Getenv("DHLOTTERY_PW")
	if dhlotteryID == "" || dhlotteryPW == "" {
		log.Fatal("Error: DHLOTTERY_ID and DHLOTTERY_PW environment variables are required")
	}
	log.Printf("Environment variables loaded successfully for user: %s", dhlotteryID)

	// 스크린샷 디렉토리 생성
	err = os.MkdirAll("screenshots", 0755)
	if err != nil {
		log.Printf("스크린샷 디렉토리 생성 실패: %v", err)
	}

	// 서울 시간대 설정
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Printf("서울 시간대 로드 실패, UTC 사용: %v", err)
		location = time.UTC
	}

	// 명령행 인수 확인 (테스트 모드)
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Println("🚀 동행복권 로또 자동 구매 봇 수동 테스트를 시작합니다...")
		fmt.Printf("⏰ 현재 서울 시간: %s\n", time.Now().In(location).Format("2006-01-02 15:04:05"))
		fmt.Println("📸 스크린샷이 screenshots/ 폴더에 저장됩니다.")
		fmt.Println("")

		// 로또 구매 봇 실행
		err = runLottoBuyBot()
		if err != nil {
			log.Printf("❌ 로또 구매 중 오류 발생: %v", err)
			return
		}

		fmt.Println("✅ 테스트 완료!")
		fmt.Println("📸 screenshots/ 폴더에서 스크린샷을 확인하세요.")
		return
	}

	// 크론 스케줄러 설정 (서울 시간 기준, 매주 토요일 오후 8시)
	c := cron.New(cron.WithLocation(location))

	// 매주 일요일 오전 6시 10분에 실행
	c.AddFunc("10 6 * * SUN", func() {
		log.Println("자동 로또 구매 작업을 시작합니다...")
		err := runLottoBuyBot()
		if err != nil {
			log.Printf("로또 구매 중 오류 발생: %v", err)
		}
	})

	log.Println("동행복권 로또 자동 구매 봇이 시작되었습니다.")
	log.Printf("스케줄: 서울 시간 기준 매주 일요일 오전 6시 10분")
	log.Printf("현재 서울 시간: %s", time.Now().In(location).Format("2006-01-02 15:04:05"))

	// 시작 시 즉시 한 번 실행 (테스트용)
	log.Println("시작 시 첫 번째 로또 구매를 실행합니다...")
	err = runLottoBuyBot()
	if err != nil {
		log.Printf("첫 번째 로또 구매 중 오류 발생: %v", err)
	}

	c.Start()

	// 프로그램이 종료되지 않도록 대기
	select {}
}

func runLottoBuyBot() error {
	// 랜덤 시드 설정
	rand.Seed(time.Now().UnixNano())

	// 랜덤한 사용자 에이전트 선택
	selectedUA := userAgents[rand.Intn(len(userAgents))]

	// Playwright 초기화 (드라이버 자동 설치)
	err := playwright.Install()
	if err != nil {
		log.Printf("Playwright 설치 실패, 계속 진행: %v", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("playwright 실행 실패: %v", err)
	}
	defer pw.Stop()

	// 브라우저 실행 (PC 버전을 위한 추가 플래그)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true), // headless 모드로 변경
		Args: []string{
			"--no-first-run",
			"--no-default-browser-check",
			"--disable-blink-features=AutomationControlled",
			"--disable-web-security",
			"--disable-features=VizDisplayCompositor",
			"--disable-mobile-emulation",    // 모바일 에뮬레이션 비활성화
			"--force-device-scale-factor=1", // 데스크톱 스케일 강제
			"--disable-touch-events",        // 터치 이벤트 비활성화
			"--disable-touch-drag-drop",     // 터치 드래그 드롭 비활성화
		},
	})
	if err != nil {
		return fmt.Errorf("브라우저 실행 실패: %v", err)
	}
	defer browser.Close()

	// 새 페이지 생성 (랜덤한 뷰포트와 사용자 에이전트 설정)
	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: &selectedUA.UserAgent,
		Viewport: &playwright.Size{
			Width:  selectedUA.Viewport.Width,
			Height: selectedUA.Viewport.Height,
		},
	})
	if err != nil {
		return fmt.Errorf("페이지 생성 실패: %v", err)
	}

	// 웹드라이버 탐지 방지 및 PC 브라우저 설정 스크립트 추가
	err = page.AddInitScript(playwright.Script{
		Content: playwright.String(`
			Object.defineProperty(navigator, 'webdriver', {
				get: () => undefined,
			});
			
			// Chrome에서 자동화 탐지 방지
			window.chrome = {
				runtime: {},
			};
			
			// 플러그인 정보 추가
			Object.defineProperty(navigator, 'plugins', {
				get: () => [1, 2, 3, 4, 5],
			});
			
			// 언어 설정
			Object.defineProperty(navigator, 'languages', {
				get: () => ['ko-KR', 'ko', 'en-US', 'en'],
			});
			
			// 모바일 감지 방지 - 데스크톱으로 강제 설정
			Object.defineProperty(navigator, 'maxTouchPoints', {
				get: () => 0,
			});
			
			Object.defineProperty(navigator, 'platform', {
				get: () => 'Win32',
			});
			
			// 모바일 User-Agent 패턴 제거
			Object.defineProperty(navigator, 'userAgent', {
				get: () => 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
			});
			
			// 터치 이벤트 지원 제거 (데스크톱으로 인식되도록)
			if ('ontouchstart' in window) {
				delete window.ontouchstart;
			}
		`),
	})
	if err != nil {
		log.Printf("초기화 스크립트 추가 실패: %v", err)
	}

	// 로그인 수행
	err = login(page)
	if err != nil {
		return fmt.Errorf("로그인 실패: %v", err)
	}

	// 로또 구매 수행
	err = buyLotto(page)
	if err != nil {
		return fmt.Errorf("로또 구매 실패: %v", err)
	}

	return nil
}

func login(page playwright.Page) error {
	log.Println("동행복권 로그인 페이지로 이동 중...")

	// 로그인 페이지로 이동
	_, err := page.Goto(LOGIN_URL)
	if err != nil {
		return err
	}

	// 페이지 로드 대기 (랜덤한 시간)
	randomDelay(2000, 4000)

	// 스크린샷 저장
	timestamp := time.Now().Format("20060102_150405")
	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/01_login_page_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("스크린샷 저장 실패: %v", err)
	}

	// 사람처럼 천천히 입력 (새로운 사이트 구조에 맞게 수정)
	err = humanLikeType(page, "#inpUserId", os.Getenv("DHLOTTERY_ID"))
	if err != nil {
		return fmt.Errorf("아이디 입력 실패: %v", err)
	}

	// 잠시 대기
	randomDelay(500, 1500)

	err = humanLikeType(page, "#inpUserPswdEncn", os.Getenv("DHLOTTERY_PW"))
	if err != nil {
		return fmt.Errorf("비밀번호 입력 실패: %v", err)
	}

	// 입력 완료 후 스크린샷
	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/02_login_filled_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("스크린샷 저장 실패: %v", err)
	}

	// 잠시 대기 후 로그인 버튼 클릭
	randomDelay(1000, 2000)

	// 로그인 버튼이 나타날 때까지 대기
	_, err = page.WaitForSelector("#btnLogin", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		log.Printf("로그인 버튼 대기 실패: %v", err)
	}

	err = page.Click("#btnLogin")
	if err != nil {
		return fmt.Errorf("로그인 버튼 클릭 실패: %v", err)
	}

	// 로그인 완료 대기
	randomDelay(3000, 5000)

	// 로그인 완료 후 스크린샷
	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/03_login_complete_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("스크린샷 저장 실패: %v", err)
	}

	log.Println("로그인 완료")
	return nil
}

func buyLotto(page playwright.Page) error {
	log.Println("로또 구매 프로세스를 시작합니다...")

	timestamp := time.Now().Format("20060102_150405")

	// 로그인 후 메인 페이지로 이동
	_, err := page.Goto("https://dhlottery.co.kr/common.do?method=main")
	if err != nil {
		log.Printf("메인 페이지 이동 실패: %v", err)
		return err
	}

	// 페이지 로드 대기
	randomDelay(3000, 5000)

	// 메인 페이지 스크린샷
	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/04_main_page_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("스크린샷 저장 실패: %v", err)
	}

	// 현재 페이지 정보 확인
	currentURL := page.URL()
	title, _ := page.Title()
	log.Printf("메인 페이지 URL: %s", currentURL)
	log.Printf("메인 페이지 제목: %s", title)

	// JavaScript로 로또 구매 팝업 열기 - 여러 방법 시도
	log.Println("JavaScript로 로또 구매 팝업을 엽니다...")

	// 다양한 JavaScript 함수들을 순서대로 시도
	jsFunctions := []string{
		`goLottoBuy(2)`,
		`goLottoBuy('2')`,
		`goGame('LO40', 'Y')`,
		`goGame('LO40')`,
		`window.open('https://el.dhlottery.co.kr/game/TotalGame.jsp?LottoId=LO40', '_blank', 'width=1024,height=768,scrollbars=yes,resizable=yes')`,
		`window.open('https://el.dhlottery.co.kr/game/TotalGame.jsp?LottoId=LO40', 'lotto_game', 'width=1024,height=768,scrollbars=yes,resizable=yes')`,
		`var popup = window.open('https://el.dhlottery.co.kr/game/TotalGame.jsp?LottoId=LO40', '_blank'); popup.focus();`,
	}

	var jsSuccess bool
	for i, jsCode := range jsFunctions {
		log.Printf("JavaScript 시도 %d: %s", i+1, jsCode)
		_, err = page.Evaluate(jsCode)
		if err == nil {
			log.Printf("JavaScript 성공: %s", jsCode)
			jsSuccess = true
			break
		}
		log.Printf("JavaScript 실패: %s - %v", jsCode, err)
		randomDelay(1000, 2000) // 각 시도 사이에 잠시 대기
	}

	if !jsSuccess {
		log.Println("모든 JavaScript 함수 호출이 실패했습니다. 직접 팝업을 엽니다...")
	}

	// 팝업이 열릴 때까지 대기
	randomDelay(3000, 5000)

	// 새로 열린 팝업/탭 확인
	context := page.Context()
	pages := context.Pages()

	var lottoBuyPage playwright.Page = page
	var foundLottoPage bool

	// 모든 페이지를 확인하여 올바른 로또 구매 페이지 찾기
	for i, p := range pages {
		url := p.URL()
		log.Printf("페이지 %d URL: %s", i, url)

		// 이벤트 팝업인지 확인하고 닫기
		if strings.Contains(url, "popupOne") {
			log.Println("이벤트 팝업을 감지했습니다. 팝업을 닫습니다.")
			err = p.Close()
			if err != nil {
				log.Printf("이벤트 팝업 닫기 실패: %v", err)
			}
			continue
		}

		// 로또 구매 페이지인지 확인
		if strings.Contains(url, "el.dhlottery.co.kr") && strings.Contains(url, "TotalGame.jsp") {
			lottoBuyPage = p
			foundLottoPage = true
			log.Printf("올바른 로또 구매 페이지를 찾았습니다: %s", url)
			break
		}
	}

	// 로또 구매 페이지를 찾지 못했다면 새 컨텍스트에서 직접 접근
	if !foundLottoPage {
		log.Println("로또 구매 페이지를 찾지 못했습니다. 새 브라우저 컨텍스트에서 직접 접근합니다...")

		// 브라우저에서 새 컨텍스트 생성 (PC 환경으로 강제)
		browser := page.Context().Browser()
		newContext, err := browser.NewContext(playwright.BrowserNewContextOptions{
			UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
			Viewport: &playwright.Size{
				Width:  1920,
				Height: 1080,
			},
			ExtraHttpHeaders: map[string]string{
				"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
				"Accept-Language": "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7",
				"Cache-Control":   "no-cache",
				"Pragma":          "no-cache",
			},
		})
		if err != nil {
			log.Printf("새 컨텍스트 생성 실패: %v", err)
			return fmt.Errorf("새 브라우저 컨텍스트 생성 실패")
		}

		// 기존 페이지의 쿠키를 새 컨텍스트로 복사 (로그인 세션 유지)
		log.Println("로그인 세션을 새 컨텍스트로 복사 중...")
		cookies, err := page.Context().Cookies()
		if err != nil {
			log.Printf("쿠키 가져오기 실패: %v", err)
		} else {
			// Cookie를 OptionalCookie로 변환
			optionalCookies := make([]playwright.OptionalCookie, len(cookies))
			for i, cookie := range cookies {
				optionalCookies[i] = playwright.OptionalCookie{
					Name:     cookie.Name,
					Value:    cookie.Value,
					Domain:   &cookie.Domain,
					Path:     &cookie.Path,
					Expires:  &cookie.Expires,
					HttpOnly: &cookie.HttpOnly,
					Secure:   &cookie.Secure,
					SameSite: cookie.SameSite,
				}
			}

			err = newContext.AddCookies(optionalCookies)
			if err != nil {
				log.Printf("쿠키 추가 실패: %v", err)
			} else {
				log.Printf("로그인 쿠키 %d개를 새 컨텍스트에 추가했습니다.", len(cookies))
			}
		}

		// 새 페이지 생성
		newPage, err := newContext.NewPage()
		if err != nil {
			log.Printf("새 페이지 생성 실패: %v", err)
			newContext.Close()
			return fmt.Errorf("새 페이지 생성 실패")
		}

		// PC 브라우저 설정 스크립트 추가 (강화된 버전)
		err = newPage.AddInitScript(playwright.Script{
			Content: playwright.String(`
				// 데스크톱으로 강제 설정
				Object.defineProperty(navigator, 'maxTouchPoints', {
					get: () => 0,
				});
				Object.defineProperty(navigator, 'platform', {
					get: () => 'Win32',
				});
				Object.defineProperty(navigator, 'userAgent', {
					get: () => 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
				});
				
				// 터치 이벤트 지원 제거
				if ('ontouchstart' in window) {
					delete window.ontouchstart;
				}
				
				// 모바일 리다이렉트 방지 - 페이지 로드 전에 실행
				window.addEventListener('DOMContentLoaded', function() {
					// 모든 location.href 변경 시도를 차단
					var originalLocation = window.location;
					Object.defineProperty(window, 'location', {
						get: function() { return originalLocation; },
						set: function(val) {
							// 모바일 사이트로의 리다이렉트 차단
							if (typeof val === 'string' && val.includes('m.dhlottery.co.kr')) {
								console.log('모바일 리다이렉트 차단:', val);
								return;
							}
							originalLocation.href = val;
						}
					});
				});
				
				// 즉시 실행으로 모바일 감지 스크립트 무력화
				if (typeof filter !== 'undefined') {
					filter = 'win16|win32|win64|macintel|linux x86_64|linux i686';
				}
			`),
		})
		if err != nil {
			log.Printf("초기화 스크립트 추가 실패: %v", err)
		}

		// 로또 구매 페이지로 직접 이동
		log.Println("새 컨텍스트에서 로또 구매 페이지로 이동 중...")
		_, err = newPage.Goto("https://el.dhlottery.co.kr/game/TotalGame.jsp?LottoId=LO40")
		if err != nil {
			log.Printf("새 컨텍스트에서 로또 페이지 이동 실패: %v", err)
			newContext.Close()
			return fmt.Errorf("로또 구매 페이지에 접근할 수 없습니다")
		}

		// 페이지 로드 대기
		err = newPage.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State: playwright.LoadStateDomcontentloaded,
		})
		if err != nil {
			log.Printf("새 페이지 로드 대기 실패: %v", err)
		}

		randomDelay(3000, 5000)

		// 새 페이지 정보 확인
		newURL := newPage.URL()
		newTitle, _ := newPage.Title()
		log.Printf("새 컨텍스트 페이지 URL: %s", newURL)
		log.Printf("새 컨텍스트 페이지 제목: %s", newTitle)

		// 세션 만료 메시지 확인
		pageContent, _ := newPage.Content()
		if strings.Contains(pageContent, "시간 초과로 세션이 해제되었습니다") ||
			strings.Contains(pageContent, "로그인해 주시기 바랍니다") {
			log.Println("세션이 만료되었습니다. 새 컨텍스트에서 다시 로그인합니다...")

			// 로그인 페이지로 이동
			_, err = newPage.Goto(LOGIN_URL)
			if err != nil {
				log.Printf("로그인 페이지 이동 실패: %v", err)
				newContext.Close()
				return fmt.Errorf("로그인 페이지 이동 실패")
			}

			randomDelay(2000, 3000)

			// 로그인 수행 (새로운 사이트 구조에 맞게 수정)
			err = humanLikeType(newPage, "#inpUserId", os.Getenv("DHLOTTERY_ID"))
			if err != nil {
				log.Printf("아이디 입력 실패: %v", err)
				newContext.Close()
				return fmt.Errorf("아이디 입력 실패")
			}

			randomDelay(500, 1500)

			err = humanLikeType(newPage, "#inpUserPswdEncn", os.Getenv("DHLOTTERY_PW"))
			if err != nil {
				log.Printf("비밀번호 입력 실패: %v", err)
				newContext.Close()
				return fmt.Errorf("비밀번호 입력 실패")
			}

			randomDelay(1000, 2000)

			// 로그인 버튼이 나타날 때까지 대기
			_, err = newPage.WaitForSelector("#btnLogin", playwright.PageWaitForSelectorOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(10000),
			})
			if err != nil {
				log.Printf("로그인 버튼 대기 실패: %v", err)
			}

			err = newPage.Click("#btnLogin")
			if err != nil {
				log.Printf("로그인 버튼 클릭 실패: %v", err)
				newContext.Close()
				return fmt.Errorf("로그인 버튼 클릭 실패")
			}

			randomDelay(3000, 5000)
			log.Println("새 컨텍스트에서 로그인 완료. 다시 로또 구매 페이지로 이동합니다...")

			// 다시 로또 구매 페이지로 이동
			_, err = newPage.Goto("https://el.dhlottery.co.kr/game/TotalGame.jsp?LottoId=LO40")
			if err != nil {
				log.Printf("로그인 후 로또 페이지 이동 실패: %v", err)
				newContext.Close()
				return fmt.Errorf("로그인 후 로또 페이지 이동 실패")
			}

			randomDelay(3000, 5000)

			// 페이지 로드 후 모바일 감지 스크립트 무력화
			_, err = newPage.Evaluate(`
				// 모바일 감지 변수 재정의
				if (typeof filter !== 'undefined') {
					filter = 'win16|win32|win64|macintel|linux x86_64|linux i686';
				}
				
				// navigator.platform을 다시 강제 설정
				Object.defineProperty(navigator, 'platform', {
					get: () => 'Win32',
					configurable: false
				});
			`)
			if err != nil {
				log.Printf("모바일 감지 무력화 스크립트 실행 실패: %v", err)
			}

			newURL = newPage.URL()
			newTitle, _ = newPage.Title()
			log.Printf("로그인 후 페이지 URL: %s", newURL)
			log.Printf("로그인 후 페이지 제목: %s", newTitle)
		}

		// 올바른 로또 구매 페이지인지 확인
		if strings.Contains(newURL, "el.dhlottery.co.kr") && strings.Contains(newURL, "TotalGame.jsp") {
			lottoBuyPage = newPage
			foundLottoPage = true
			log.Println("새 컨텍스트에서 올바른 로또 구매 페이지를 열었습니다!")
		} else {
			log.Printf("새 컨텍스트에서도 올바른 페이지가 열리지 않았습니다: %s", newURL)
			newContext.Close()
			return fmt.Errorf("로또 구매 페이지 로드 실패")
		}
	}

	// 팝업으로 포커스 이동
	if foundLottoPage && lottoBuyPage != page {
		err = lottoBuyPage.BringToFront()
		if err != nil {
			log.Printf("팝업 포커스 이동 실패: %v", err)
		}
	}

	// 팝업 페이지 로드 대기
	err = lottoBuyPage.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})
	if err != nil {
		log.Printf("팝업 페이지 로드 대기 실패: %v", err)
	}

	randomDelay(3000, 5000)

	// 팝업 페이지 정보 확인
	popupURL := lottoBuyPage.URL()
	popupTitle, _ := lottoBuyPage.Title()
	log.Printf("로또 구매 팝업 URL: %s", popupURL)
	log.Printf("로또 구매 팝업 제목: %s", popupTitle)

	// 팝업 페이지 스크린샷
	_, err = lottoBuyPage.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/05_lotto_popup_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("스크린샷 저장 실패: %v", err)
	}

	// 올바른 로또 구매 페이지가 아니라면 에러 반환
	if !strings.Contains(popupURL, "el.dhlottery.co.kr") || !strings.Contains(popupURL, "TotalGame.jsp") {
		log.Printf("올바른 로또 구매 페이지가 아닙니다. 현재 URL: %s", popupURL)
		return fmt.Errorf("로또 구매 페이지 로드 실패: 잘못된 페이지가 열렸습니다")
	}

	// iframe으로 전환 - 실제 로또 구매 인터페이스는 iframe 안에 있음
	log.Println("iframe으로 전환하여 실제 로또 구매 페이지에 접근합니다...")

	// iframe이 로드될 때까지 대기
	_, err = lottoBuyPage.WaitForSelector("#ifrm_tab", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		log.Printf("iframe 로드 대기 실패: %v", err)
		return fmt.Errorf("iframe을 찾을 수 없습니다")
	}

	// iframe 요소 가져오기
	iframe, err := lottoBuyPage.QuerySelector("#ifrm_tab")
	if err != nil {
		log.Printf("iframe 요소 찾기 실패: %v", err)
		return fmt.Errorf("iframe 요소를 찾을 수 없습니다")
	}

	// iframe의 content frame 가져오기
	iframeContent, err := iframe.ContentFrame()
	if err != nil {
		log.Printf("iframe content frame 가져오기 실패: %v", err)
		return fmt.Errorf("iframe content frame을 가져올 수 없습니다")
	}

	// iframe이 완전히 로드될 때까지 대기
	err = iframeContent.WaitForLoadState(playwright.FrameWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})
	if err != nil {
		log.Printf("iframe 로드 상태 대기 실패: %v", err)
	}

	randomDelay(3000, 5000)

	// iframe 내용 확인
	iframeURL := iframeContent.URL()
	log.Printf("iframe URL: %s", iframeURL)

	// iframe 스크린샷 (전체 페이지로)
	_, err = lottoBuyPage.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/05_iframe_loaded_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("iframe 스크린샷 저장 실패: %v", err)
	}

	// 자동번호발급 탭 클릭 - iframe 내에서 진행
	log.Println("iframe 내에서 자동번호발급 탭 클릭 중...")

	selectors := []string{
		"#num2",
		"a[href='#divWay2Buy1']",
		"a[onclick*='selectWayTab(1)']",
		"a:has-text('자동번호발급')",
		".tab_menu a:nth-child(2)",
		"#divWay2Buy1",
		".tab02",
		"input[name='genType'][value='1']",
	}

	var clickSuccess bool
	for _, selector := range selectors {
		err = iframeContent.Click(selector, playwright.FrameClickOptions{
			Timeout: playwright.Float(10000), // 10초 타임아웃
		})
		if err == nil {
			log.Printf("iframe에서 자동번호발급 탭 클릭 성공: %s", selector)
			clickSuccess = true
			break
		}
		log.Printf("iframe에서 셀렉터 %s로 클릭 실패: %v", selector, err)
	}

	if !clickSuccess {
		// iframe 내의 모든 요소를 더 자세히 분석
		links, _ := iframeContent.QuerySelectorAll("a")
		log.Printf("iframe에서 발견된 링크 수: %d", len(links))

		// 모든 input 요소 확인
		inputs, _ := iframeContent.QuerySelectorAll("input")
		log.Printf("iframe에서 발견된 input 요소 수: %d", len(inputs))

		// 모든 div 요소 확인
		divs, _ := iframeContent.QuerySelectorAll("div")
		log.Printf("iframe에서 발견된 div 요소 수: %d", len(divs))

		// iframe에서 '자동' 또는 '로또' 관련 텍스트 검색
		pageText, _ := iframeContent.TextContent("body")
		if strings.Contains(pageText, "자동") {
			log.Println("iframe에서 '자동' 텍스트를 발견했습니다.")
		}
		if strings.Contains(pageText, "번호발급") {
			log.Println("iframe에서 '번호발급' 텍스트를 발견했습니다.")
		}
		if strings.Contains(pageText, "로또") {
			log.Println("iframe에서 '로또' 텍스트를 발견했습니다.")
		}

		// iframe HTML 소스 덤프 (더 많은 내용)
		pageContent, err := iframeContent.Content()
		if err == nil {
			// 중요한 키워드가 포함된 부분을 찾아서 출력
			lines := strings.Split(pageContent, "\n")
			for i, line := range lines {
				if strings.Contains(line, "자동") || strings.Contains(line, "번호발급") ||
					strings.Contains(line, "tab") || strings.Contains(line, "num2") ||
					strings.Contains(line, "divWay2Buy") {
					start := i - 2
					if start < 0 {
						start = 0
					}
					end := i + 3
					if end > len(lines) {
						end = len(lines)
					}
					log.Printf("iframe 관련 HTML 라인 %d-%d:", start, end)
					for j := start; j < end; j++ {
						log.Printf("  %d: %s", j, strings.TrimSpace(lines[j]))
					}
				}
			}

			// iframe HTML도 파일로 저장
			htmlFile := fmt.Sprintf("screenshots/iframe_source_%s.html", timestamp)
			err = os.WriteFile(htmlFile, []byte(pageContent), 0644)
			if err == nil {
				log.Printf("iframe HTML 소스를 %s에 저장했습니다.", htmlFile)
			}
		}

		// 현재 페이지 상태 스크린샷
		_, err = lottoBuyPage.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String(fmt.Sprintf("screenshots/06_click_failed_%s.png", timestamp)),
		})
		if err != nil {
			log.Printf("실패 스크린샷 저장 실패: %v", err)
		}

		return fmt.Errorf("iframe에서 자동번호발급 탭을 찾을 수 없습니다. iframe이 올바르게 로드되지 않았을 수 있습니다")
	}

	randomDelay(2000, 4000)

	// 자동번호발급 탭 클릭 후 스크린샷
	_, err = lottoBuyPage.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/07_auto_tab_clicked_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("스크린샷 저장 실패: %v", err)
	}

	// 구매 수량을 5로 설정
	log.Println("iframe에서 구매 수량을 5로 설정 중...")

	// 구매 수량 셀렉터도 여러개 시도
	quantitySelectors := []string{
		"#amoundApply",
		"select[name='amoundApply']",
		"select[onchange*='paperTextChange']",
		"#buyAmountApply",
		"select[name='buyAmountApply']",
	}

	var quantitySuccess bool
	for _, selector := range quantitySelectors {
		_, err = iframeContent.SelectOption(selector, playwright.SelectOptionValues{
			Values: &[]string{"5"},
		})
		if err == nil {
			log.Printf("iframe에서 구매 수량 설정 성공: %s", selector)
			quantitySuccess = true
			break
		}
		log.Printf("iframe에서 구매 수량 셀렉터 %s 실패: %v", selector, err)
	}

	if !quantitySuccess {
		return fmt.Errorf("iframe에서 구매 수량 설정 실패: 모든 셀렉터가 실패했습니다")
	}

	randomDelay(1000, 2000)

	// 구매 수량 설정 후 스크린샷
	_, err = lottoBuyPage.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/08_quantity_set_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("스크린샷 저장 실패: %v", err)
	}

	// 확인 버튼 클릭
	log.Println("iframe에서 확인 버튼 클릭 중...")

	confirmSelectors := []string{
		"#btnSelectNum",
		"input[name='btnSelectNum']",
		"input[value='확인'].button.lrg.confirm",
		"input.button.lrg.confirm",
		"#btnAutoSelect",
		"input[onclick*='createAutoLotto']",
	}

	var confirmSuccess bool
	for _, selector := range confirmSelectors {
		err = iframeContent.Click(selector, playwright.FrameClickOptions{
			Timeout: playwright.Float(10000),
		})
		if err == nil {
			log.Printf("iframe에서 확인 버튼 클릭 성공: %s", selector)
			confirmSuccess = true
			break
		}
		log.Printf("iframe에서 확인 버튼 셀렉터 %s 실패: %v", selector, err)
	}

	if !confirmSuccess {
		return fmt.Errorf("iframe에서 확인 버튼 클릭 실패: 모든 셀렉터가 실패했습니다")
	}

	randomDelay(2000, 3000)

	// 확인 버튼 클릭 후 스크린샷
	_, err = lottoBuyPage.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/09_confirm_clicked_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("스크린샷 저장 실패: %v", err)
	}

	// 구매하기 버튼 클릭
	log.Println("iframe에서 구매하기 버튼 클릭 중...")

	buySelectors := []string{
		"#btnBuy",
		"input[name='btnBuy']",
		"input[value='구매하기'].button.buy",
		"input.button.buy",
		"#lottoPayment",
		"input[onclick*='payment']",
	}

	var buySuccess bool
	for _, selector := range buySelectors {
		err = iframeContent.Click(selector, playwright.FrameClickOptions{
			Timeout: playwright.Float(10000),
		})
		if err == nil {
			log.Printf("iframe에서 구매하기 버튼 클릭 성공: %s", selector)
			buySuccess = true
			break
		}
		log.Printf("iframe에서 구매하기 버튼 셀렉터 %s 실패: %v", selector, err)
	}

	if !buySuccess {
		return fmt.Errorf("iframe에서 구매하기 버튼 클릭 실패: 모든 셀렉터가 실패했습니다")
	}

	randomDelay(2000, 3000)

	// 구매하기 버튼 클릭 후 스크린샷
	_, err = lottoBuyPage.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/10_buy_clicked_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("스크린샷 저장 실패: %v", err)
	}

	// 팝업 확인 버튼 클릭
	log.Println("iframe에서 팝업 확인 버튼 클릭 중...")

	popupSelectors := []string{
		"input[value='확인'][onclick*='closepopupLayerConfirm']",
		"input.button.lrg.confirm[value='확인']",
		"input[onclick*='closepopup']",
		"button:has-text('확인')",
		".popup input[value='확인']",
		"#confirm_button",
		".layerPop input[value='확인']",
	}

	var popupSuccess bool
	for _, selector := range popupSelectors {
		err = iframeContent.Click(selector, playwright.FrameClickOptions{
			Timeout: playwright.Float(5000),
		})
		if err == nil {
			log.Printf("iframe에서 팝업 확인 버튼 클릭 성공: %s", selector)
			popupSuccess = true
			break
		}
		log.Printf("iframe에서 팝업 확인 버튼 셀렉터 %s 실패: %v", selector, err)
	}

	if !popupSuccess {
		log.Printf("iframe에서 팝업 확인 버튼을 찾을 수 없습니다. 팝업이 없을 수도 있습니다.")
	}

	randomDelay(2000, 3000)

	// 최종 완료 후 스크린샷
	_, err = lottoBuyPage.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("screenshots/11_purchase_complete_%s.png", timestamp)),
	})
	if err != nil {
		log.Printf("스크린샷 저장 실패: %v", err)
	}

	log.Println("로또 구매 완료!")
	return nil
}

// 사람처럼 천천히 타이핑하는 함수
func humanLikeType(page playwright.Page, selector, text string) error {
	err := page.Click(selector)
	if err != nil {
		return err
	}

	// 기존 텍스트 지우기
	err = page.Fill(selector, "")
	if err != nil {
		return err
	}

	// 한 글자씩 천천히 입력
	for _, char := range text {
		err = page.Type(selector, string(char), playwright.PageTypeOptions{
			Delay: playwright.Float(float64(rand.Intn(100) + 50)), // 50-150ms 랜덤 딜레이
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// 랜덤한 딜레이 함수 (밀리초)
func randomDelay(min, max int) {
	delay := rand.Intn(max-min) + min
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
