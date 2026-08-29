# ============================================================
#  Sportswear 本地开发模式一键启动脚本
#  用法：右键 -> "使用 PowerShell 运行"，或在终端执行：
#       powershell -ExecutionPolicy Bypass -File dev.ps1
#
#  会自动：
#    1. 启动 Docker 基础设施（PostgreSQL、Redis、MinIO）
#    2. 打开 3 个独立终端分别运行 Go 后端 / Admin / Portal
#    3. 改代码后各服务自动热重载，浏览器秒级刷新
# ============================================================

$ErrorActionPreference = "Stop"
$projectRoot = $PSScriptRoot
$composeFile  = "$projectRoot\docker-compose.yml"

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Sportswear 开发模式启动"              -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# ── 1. Docker 基础设施 ────────────────────────────────────────
Write-Host "[1/4] 启动 Docker 基础设施..." -ForegroundColor Yellow

# 确保 Docker 可用
$dockerOk = $false
docker compose version 2>$null | Out-Null
if ($LASTEXITCODE -eq 0) { $dockerOk = $true }
if (-not $dockerOk) {
    Write-Host "[ERROR] 未检测到 Docker，请先安装 Docker Desktop" -ForegroundColor Red
    pause; exit 1
}

# 只启动基础设施，不构建业务服务
docker compose -p sportswear -f $composeFile up -d postgres redis minio 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Docker 基础设施启动失败" -ForegroundColor Red
    pause; exit 1
}

Write-Host "  PostgreSQL : localhost:5432" -ForegroundColor Green
Write-Host "  Redis      : localhost:6379" -ForegroundColor Green
Write-Host "  MinIO      : localhost:9000 (console :9001)" -ForegroundColor Green

# 等待 PostgreSQL 就绪
Write-Host "  等待 PostgreSQL 就绪中..." -ForegroundColor Gray
$retry = 0
do {
    Start-Sleep -Seconds 2
    $healthy = docker inspect --format='{{.State.Health.Status}}' sportswear-postgres 2>$null
    $retry++
} while ($healthy -ne 'healthy' -and $retry -lt 20)
if ($healthy -eq 'healthy') {
    Write-Host "  PostgreSQL 已就绪" -ForegroundColor Green
} else {
    Write-Host "  [!] PostgreSQL 状态异常: $healthy" -ForegroundColor DarkYellow
}

# ── 2. 安装依赖（如需要）─────────────────────────────────────
Write-Host ""
Write-Host "[2/4] 检查依赖..." -ForegroundColor Yellow

$adminNM  = "$projectRoot\sportswear-admin\node_modules"
$portalNM = "$projectRoot\sportswear-portal\node_modules"

if (-not (Test-Path $adminNM)) {
    Write-Host "  安装 Admin 依赖..." -ForegroundColor Gray
    Set-Location "$projectRoot\sportswear-admin"
    npm install 2>&1 | Out-Null
    Write-Host "  Admin 依赖安装完成" -ForegroundColor Green
}
if (-not (Test-Path $portalNM)) {
    Write-Host "  安装 Portal 依赖..." -ForegroundColor Gray
    Set-Location "$projectRoot\sportswear-portal"
    npm install 2>&1 | Out-Null
    Write-Host "  Portal 依赖安装完成" -ForegroundColor Green
}

# ── 3. 启动业务服务（独立窗口）────────────────────────────
Write-Host ""
Write-Host "[3/4] 启动业务服务..." -ForegroundColor Yellow

# 3a. Go 后端（端口 8080）
Write-Host "  启动 Go 后端 :8080 ..." -ForegroundColor Gray
Start-Process -FilePath "powershell" -ArgumentList @(
    "-NoExit",
    "-Command",
    "Write-Host '=== Go Backend (sportswear-backend) ===' -ForegroundColor Cyan; cd '$projectRoot\sportswear-backend'; go run ./cmd/server"
) -WindowStyle Normal

Start-Sleep -Seconds 1

# 3b. Admin 管理后台（端口 5173，Vite 自动代理 /api → :8080）
Write-Host "  启动 Admin 管理后台 :5173 ..." -ForegroundColor Gray
Start-Process -FilePath "powershell" -ArgumentList @(
    "-NoExit",
    "-Command",
    "Write-Host '=== Admin (sportswear-admin) ===' -ForegroundColor Cyan; cd '$projectRoot\sportswear-admin'; npm run dev"
) -WindowStyle Normal

Start-Sleep -Seconds 1

# 3c. Portal 门户前端（端口 3000）
Write-Host "  启动 Portal 门户 :3000 ..." -ForegroundColor Gray
Start-Process -FilePath "powershell" -ArgumentList @(
    "-NoExit",
    "-Command",
    "Write-Host '=== Portal (sportswear-portal) ===' -ForegroundColor Cyan; cd '$projectRoot\sportswear-portal'; npm run dev"
) -WindowStyle Normal

# ── 4. 输出访问地址 ──────────────────────────────────────────
Write-Host ""
Write-Host "[4/4] 等待服务就绪后访问：" -ForegroundColor Yellow
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  管理后台 : http://localhost:5173" -ForegroundColor White
Write-Host "            (Vite HMR – 改代码自动热更新)" -ForegroundColor Gray
Write-Host ""
Write-Host "  门户前端 : http://localhost:3000" -ForegroundColor White
Write-Host "            (Nuxt HMR – 改代码自动热更新)" -ForegroundColor Gray
Write-Host ""
Write-Host "  后端 API : http://localhost:8080" -ForegroundColor White
Write-Host "            (Go 需要手动重启，或安装 air 热重载)" -ForegroundColor Gray
Write-Host ""
Write-Host "  MinIO    : http://localhost:9001" -ForegroundColor White
Write-Host "            账号: minioadmin / minioadmin" -ForegroundColor Gray
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "提示：关闭各窗口即可停止对应服务。" -ForegroundColor DarkYellow
Write-Host "      Docker 基础设施可用以下命令停止：" -ForegroundColor DarkYellow
Write-Host "      docker compose -p sportswear -f `"$composeFile`" stop" -ForegroundColor DarkGray
Write-Host ""

pause