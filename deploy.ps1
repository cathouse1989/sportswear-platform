# ============================================================
#  Sportswear 一键部署脚本
#  用法：右键 -> "使用 PowerShell 运行"，或在终端执行：
#       powershell -ExecutionPolicy Bypass -File deploy.ps1
# ============================================================

$ErrorActionPreference = "Stop"
$projectRoot = $PSScriptRoot
$composeFile  = "$projectRoot\docker-compose.yml"
$logFile      = "$projectRoot\deploy.log"

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Sportswear 一键部署"                   -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 确保 docker compose 子命令可用（老版本是 docker-compose） ----------
$composeCmd = $null
if (Get-Command "docker" -ErrorAction SilentlyContinue) {
    docker compose version 2>$null | Out-Null
    if ($LASTEXITCODE -eq 0) {
        $composeCmd = @("docker", "compose")
        Write-Host "[OK] docker compose 可用" -ForegroundColor Green
    }
}
if (-not $composeCmd) {
    if (Get-Command "docker-compose" -ErrorAction SilentlyContinue) {
        $composeCmd = @("docker-compose")
        Write-Host "[OK] docker-compose 可用" -ForegroundColor Green
    }
}
if (-not $composeCmd) {
    Write-Host "[ERROR] 未检测到 Docker，请先安装 Docker Desktop" -ForegroundColor Red
    pause
    exit 1
}

# 1. 拉取基础镜像 -------------------------------------------------
Write-Host ""
Write-Host "[1/4] 拉取基础镜像..." -ForegroundColor Yellow
& docker pull postgres:15-alpine  2>&1 | Out-Null
& docker pull redis:7-alpine      2>&1 | Out-Null
& docker pull minio/minio:latest  2>&1 | Out-Null
& docker pull node:22-alpine      2>&1 | Out-Null
& docker pull nginx:1.27-alpine   2>&1 | Out-Null
& docker pull golang:1.25-alpine  2>&1 | Out-Null
& docker pull alpine:3.20         2>&1 | Out-Null
Write-Host "  基础镜像拉取完成" -ForegroundColor Green

# 2. 构建并启动所有服务 -------------------------------------------
Write-Host ""
Write-Host "[2/4] 构建镜像并启动容器（--build）..." -ForegroundColor Yellow
Write-Host "  编排文件: $composeFile" -ForegroundColor Gray

$startTime = Get-Date
& $composeCmd[0] $composeCmd[1] -p sportswear -f $composeFile up -d --build 2>&1 | Tee-Object -FilePath $logFile

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "[ERROR] 部署失败，请查看日志：$logFile" -ForegroundColor Red
    pause
    exit 1
}

$elapsed = (Get-Date) - $startTime
Write-Host "  构建+启动耗时: $($elapsed.TotalSeconds.ToString('0'))s" -ForegroundColor Green

# 3. 等待健康检查 -------------------------------------------------
Write-Host ""
Write-Host "[3/4] 等待数据库就绪..." -ForegroundColor Yellow
$retry = 0
do {
    Start-Sleep -Seconds 2
    $healthy = & docker inspect --format='{{.State.Health.Status}}' sportswear-postgres 2>$null
    $retry++
} while ($healthy -ne 'healthy' -and $retry -lt 15)
if ($healthy -eq 'healthy') {
    Write-Host "  PostgreSQL 就绪" -ForegroundColor Green
} else {
    Write-Host "  PostgreSQL 状态: $healthy（继续等待）" -ForegroundColor DarkYellow
}

# 4. 显示运行状态 -------------------------------------------------
Write-Host ""
Write-Host "[4/4] 容器运行状态：" -ForegroundColor Yellow
Write-Host ""

& docker ps -a --filter "name=sportswear" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# 访问地址 --------------------------------------------------------
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  部署完成！访问地址："                  -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  后端 API   : http://localhost:8080"   -ForegroundColor White
Write-Host "  管理后台   : http://localhost:8090"   -ForegroundColor White
Write-Host "  门户前端   : http://localhost:3000"   -ForegroundColor White
Write-Host "  MinIO      : http://localhost:9001"   -ForegroundColor White
Write-Host ""
Write-Host "  默认管理员: admin@sportswear.com / Admin@123456" -ForegroundColor Gray
Write-Host "  MinIO账号 : minioadmin / minioadmin"             -ForegroundColor Gray
Write-Host "========================================" -ForegroundColor Cyan

pause