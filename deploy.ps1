# ============================================================
#  Sportswear 一键部署脚本
#  用法：右键 -> "使用 PowerShell 运行"，或在终端执行：
#       powershell -ExecutionPolicy Bypass -File deploy.ps1
# ============================================================

$ErrorActionPreference = "Stop"
$projectRoot = $PSScriptRoot
$composeFile  = "$projectRoot\docker-compose.yml"
$logFile      = "$projectRoot\deploy.log"

# ------------------------------------------------------------
# 兼容 Windows PowerShell 5.1：
# 5.1 会把外部命令(docker 等)写入 stderr 的正常输出也看作
# "错误"，在 $ErrorActionPreference="Stop" 下会直接终止脚本并报
# NativeCommandError。下面函数保证外部命令正常运行：
#   - Invoke-Native  ：执行外部命令，期间临时放宽错误策略，
#                         把 stderr 作为普通文本显示，返回退出码
# ------------------------------------------------------------
function Invoke-Native {
    param([scriptblock]$ScriptBlock, [string]$LogFile = "")
    $savedEAP = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    $exitCode = -1
    try {
        $output = & $ScriptBlock 2>&1
        $exitCode = $LASTEXITCODE

        # 输出到控制台（Write-Host 不污染返回值），同时追加到日志（Add-Content 不输出到管道）
        if ($output) {
            $output | ForEach-Object { Write-Host "$_" }
            if ($LogFile) {
                $output | ForEach-Object { "$_" } | Add-Content -Path $LogFile -Encoding utf8
            }
        }
    }
    finally {
        $ErrorActionPreference = $savedEAP
    }
    return $exitCode
}

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

# 0. 配置 Docker 镜像加速器（国内网络环境必需） --------------------
function Configure-DockerMirrors {
    $daemonPath = "$env:USERPROFILE\.docker\daemon.json"
    $defaultMirrors = @(
        "https://docker.1ms.run",
        "https://docker.mproxy.org",
        "https://docker.1panel.live",
        "https://docker.m.daocloud.io"
    )

    $needRestart = $false
    $daemonConfig = $null

    if (Test-Path $daemonPath) {
        try {
            $daemonConfig = Get-Content $daemonPath -Raw | ConvertFrom-Json -ErrorAction SilentlyContinue
        } catch {
            $daemonConfig = $null
        }
    }

    if ($daemonConfig -and $daemonConfig.'registry-mirrors') {
        Write-Host "  镜像加速器已配置: $($daemonConfig.'registry-mirrors' -join ', ')" -ForegroundColor Green
        return $false  # 无需重启
    }

    # 需添加镜像加速器 -> 更新 daemon.json
    if (-not $daemonConfig) {
        $daemonConfig = New-Object PSObject -Property @{
            'builder'           = New-Object PSObject -Property @{
                'gc' = New-Object PSObject -Property @{
                    'defaultKeepStorage' = '20GB'
                    'enabled'            = $true
                }
            }
            'experimental'      = $false
        }
    }

    # 添加 registry-mirrors 属性
    if (-not $daemonConfig.'registry-mirrors') {
        Add-Member -InputObject $daemonConfig -NotePropertyName 'registry-mirrors' -NotePropertyValue $defaultMirrors -Force
    }

    $daemonConfig | ConvertTo-Json -Depth 10 | Set-Content $daemonPath
    Write-Host "  镜像加速器已添加到 daemon.json" -ForegroundColor Green
    return $true  # 需要重启
}

function Wait-DockerReady {
    $retry = 0
    do {
        Start-Sleep -Seconds 3
        $retry++
    } while ((docker version 2>$null) -and $LASTEXITCODE -ne 0 -and $retry -lt 30)
}

Write-Host ""
Write-Host "[0/4] 配置 Docker 镜像加速器..." -ForegroundColor Yellow
if (Configure-DockerMirrors) {
    Write-Host "  正在重启 Docker Desktop 以应用镜像加速器..." -ForegroundColor Yellow
    # 尝试多种重启方式
    $restarted = $false
    try { Restart-Service com.docker.service -ErrorAction Stop; $restarted = $true } catch {}
    if (-not $restarted) {
        try {
            Stop-Process -Name "Docker Desktop" -Force -ErrorAction SilentlyContinue
            Start-Sleep -Seconds 1
            Start-Process "Docker Desktop.exe" -ErrorAction SilentlyContinue
            $restarted = $true
        } catch {}
    }
    if (-not $restarted) {
        Write-Host "  [WARNING] 无法自动重启 Docker，请手动重启 Docker Desktop，然后按任意键继续" -ForegroundColor DarkYellow
        pause
    }
    Wait-DockerReady
    Write-Host "  Docker 已就绪" -ForegroundColor Green
}

# 1. 拉取基础镜像 -------------------------------------------------
Write-Host ""
Write-Host "[1/4] 拉取基础镜像..." -ForegroundColor Yellow

$baseImages = @(
    "postgres:15-alpine",
    "redis:7-alpine",
    "minio/minio:latest",
    "node:22-alpine",
    "nginx:1.27-alpine",
    "golang:1.25-alpine",
    "alpine:3.20"
)

$savedEAP = $ErrorActionPreference
$ErrorActionPreference = "Continue"

foreach ($img in $baseImages) {
    $ok = $false
    for ($attempt = 1; $attempt -le 3; $attempt++) {
        & docker pull $img 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) { $ok = $true; break }
        Start-Sleep -Seconds 2
    }
    if ($ok) {
        Write-Host "  [OK] $img" -ForegroundColor Green
    } else {
        Write-Host "  [WARN] $img 拉取失败（将在构建时重试）" -ForegroundColor DarkYellow
    }
}
$ErrorActionPreference = $savedEAP
Write-Host "  基础镜像拉取完成" -ForegroundColor Green

# 2. 构建并启动所有服务 -------------------------------------------
Write-Host ""
Write-Host "[2/4] 构建镜像并启动容器（--build）..." -ForegroundColor Yellow
Write-Host "  编排文件: $composeFile" -ForegroundColor Gray

$startTime = Get-Date
$buildRetryMax = 3
$buildOk = $false

for ($buildAttempt = 1; $buildAttempt -le $buildRetryMax; $buildAttempt++) {
    Write-Host "  尝试构建启动 ($buildAttempt/$buildRetryMax)..." -ForegroundColor Gray
    # 使用 Invoke-Native 避免 PowerShell 5.1 的 NativeCommandError 终止问题；输出实时打印并追加到 deploy.log
    $exitCode = Invoke-Native { & $composeCmd[0] $composeCmd[1] -p sportswear -f $composeFile up -d --build } -LogFile $logFile

    if ($exitCode -eq 0) {
        $buildOk = $true
        break
    }

    Write-Host "  [WARN] 构建失败，正在重试..." -ForegroundColor DarkYellow
    Start-Sleep -Seconds 3
}

if (-not $buildOk) {
    Write-Host ""
    Write-Host "[ERROR] 部署失败，请查看日志：$logFile" -ForegroundColor Red
    Write-Host "  可能原因：" -ForegroundColor DarkYellow
    Write-Host "    1. Docker Hub 镜像无法访问（检查 daemon.json registry-mirrors）" -ForegroundColor Gray
    Write-Host "    2. Go 模块下载失败（检查 Dockerfile 中的 GOPROXY）" -ForegroundColor Gray
    Write-Host "    3. npm 依赖安装失败（检查 Dockerfile 中的 npm registry）" -ForegroundColor Gray
    Write-Host "  解决方案：" -ForegroundColor DarkYellow
    Write-Host "    - 在 ~/.docker/daemon.json 添加 registry-mirrors" -ForegroundColor Gray
    Write-Host "    - 或设置 BASE_REGISTRY 到 .env，例如：" -ForegroundColor Gray
    Write-Host "        echo BASE_REGISTRY=docker.mirrors.ustc.edu.cn/library/ > .env" -ForegroundColor Gray
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

$savedEAP = $ErrorActionPreference
$ErrorActionPreference = "Continue"
& docker ps -a --filter "name=sportswear" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
$ErrorActionPreference = $savedEAP

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
Write-Host ""
Write-Host "  SEO 环境变量配置（如需）:" -ForegroundColor Yellow
Write-Host "    请查看 sportswear-portal/.env.example 并设置 .env 文件" -ForegroundColor Gray
Write-Host "    NUXT_PUBLIC_GA_MEASUREMENT_ID     # Google Analytics 4" -ForegroundColor Gray
Write-Host "    NUXT_PUBLIC_GOOGLE_SITE_VERIFICATION # Search Console" -ForegroundColor Gray
Write-Host "    NUXT_PUBLIC_SITE_URL               # 生产域名" -ForegroundColor Gray
Write-Host "========================================" -ForegroundColor Cyan

pause