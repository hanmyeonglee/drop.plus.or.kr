# PLUS Drop Server

미니멀하고 안전한 프라이빗 파일 저장소입니다. Go 언어를 활용한 빠르고 가벼운 백엔드와, JavaScript 없이 HTML/CSS만으로 구현된 세련되고 심플한 프론트엔드를 특징으로 합니다.

## 🚀 Features

- **단순하고 미니멀한 UI**: 불필요한 스크립트를 배제하고 순수 HTML/CSS 만으로 깔끔하고 반응성이 뛰어난 UI를 제공합니다.
- **빠른 업로드 및 다운로드**: Go의 강력한 성능을 활용하여 안정적인 파일 전송을 지원합니다.
- **SQLite3 메타데이터 관리**: 별도의 무거운 DB 서버 없이 SQLite3를 통해 파일 메타데이터(이름, 크기, 업로드 시간 등)를 효율적으로 관리합니다.
- **보안 파일 저장**: 업로드된 파일은 시스템에 원본 파일명이 아닌 고유 UUID로 저장되어 보안을 강화합니다.
- **가벼운 배포 환경**: Alpine Linux 기반의 Docker 및 멀티스테이지 빌드를 통해 매우 가볍게 컨테이너를 구동할 수 있으며 `nginx-proxy` 연동이 준비되어 있습니다.

## 🛠️ Tech Stack

- **Backend**: Go (net/http, html/template)
- **Frontend**: HTML5, Vanilla CSS (JS-Free)
- **Database**: SQLite3
- **Infrastructure**: Docker, docker-compose (nginx-proxy)

## 📦 Installation & Run

로컬 스탠드얼론(Standalone) 모드로 구동하거나, Nginx Proxy 기반으로 구동할 수 있습니다.

### 환경 변수 세팅

서버 실행을 위해 다음 환경 변수들을 설정할 수 있습니다. (현재는 테스트 모드로 동작 가능)
- `PORT`: 실행 포트 (기본값: `8080`)
- `DATA_DIR`: 파일과 DB가 저장될 경로 (기본값: `./data`)
- `MAX_UPLOAD_SIZE_MB`: 단일 파일 업로드 최대 용량 제한 (기본값: `50`)
- `ENTRA_CLIENT_ID`, `ENTRA_TENANT_ID`, `ENTRA_CLIENT_SECRET`: Microsoft Entra ID 연동용

### Standalone 모드로 실행 (nginx-proxy 없이)

```bash
docker-compose -f docker-compose.yml up --build
```
`http://localhost:8080` 으로 접속하여 사용할 수 있습니다.

### Nginx Proxy 연동 모드로 실행

기본 명령어 실행 시 `docker-compose.override.yml`이 병합되어 Nginx Proxy와 함께 동작합니다.
```bash
docker-compose up --build
```
이후 `/etc/hosts` 등에 `drop.plus.or.kr` 도메인을 로컬 IP로 매핑하여 접속합니다.

## 📋 TODO

- [ ] **파일 검색 기능**: 업로드한 파일을 파일명으로 빠르게 검색하는 기능
- [ ] **즐겨찾기 기능**: 자주 사용하는 파일에 즐겨찾기(Star) 지정
- [ ] **디렉토리 분할 기능**: 사용자가 폴더(디렉토리)를 생성하여 파일을 구조적으로 관리
- [ ] **PWA**: 앱으로 사용할 수 있도록 PWA 등록
