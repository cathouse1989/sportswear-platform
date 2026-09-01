$ErrorActionPreference = 'Stop'
$base = 'd:\code\sportswear-platform'
$tok = Get-Content "$base\_tok.txt"
$h = @{ Authorization = "Bearer $tok" }
$api = 'http://localhost:8080/api/v1'

Add-Type -AssemblyName System.Drawing
function New-Cover([string]$path, [int]$w, [int]$h, [string]$c1, [string]$c2, [string]$text) {
  $bmp = New-Object System.Drawing.Bitmap($w, $h)
  $g = [System.Drawing.Graphics]::FromImage($bmp)
  $rect = New-Object System.Drawing.Rectangle(0, 0, $w, $h)
  $b = New-Object System.Drawing.Drawing2D.LinearGradientBrush($rect, [System.Drawing.ColorTranslator]::FromHtml($c1), [System.Drawing.ColorTranslator]::FromHtml($c2), 45)
  $g.FillRectangle($b, $rect)
  $f = New-Object System.Drawing.Font('Arial', ($h / 9))
  $sf = New-Object System.Drawing.StringFormat
  $sf.Alignment = 'Center'; $sf.LineAlignment = 'Center'
  $g.DrawString($text, $f, [System.Drawing.Brushes]::White, (New-Object System.Drawing.RectangleF(0, 0, $w, $h)), $sf)
  $g.Dispose(); $bmp.Save($path, [System.Drawing.Imaging.ImageFormat]::Png); $bmp.Dispose()
}

function Upload-Media([string]$file, [string]$category, [string]$title) {
  $resp = curl.exe -s -H "Authorization: Bearer $tok" -F "file=@$file" "$api/admin/media/upload?category=$category" | ConvertFrom-Json
  if (-not $resp.success) { throw "upload failed: $file" }
  $id = $resp.data.id
  # 补充媒体元数据（alt/title 建立可维护的关联描述）
  curl.exe -s -X PUT -H "Authorization: Bearer $tok" -H 'Content-Type: application/json' `
    -d ("{`"alt`":`"" + $title + "`",`"title`":`"" + $title + "`"}") "$api/admin/media/$id" | Out-Null
  return $resp.data.url
}

function Set-JsonFile([string]$path, [hashtable]$body) {
  $json = $body | ConvertTo-Json -Depth 6
  Set-Content -Path $path -Value $json -Encoding UTF8
}

# ============ 1. 认证证书图片 ============
$certs = @(
  @{ id = '572cd93a-5de0-4058-8e9f-20f1e7189fc8'; name = 'BSCI Compliance';        c1 = '#1B4332'; c2 = '#95D5B2' },
  @{ id = '9154cddb-edd5-476d-8249-b61cf8930739'; name = 'ISO 9001:2015';          c1 = '#0D1B2A'; c2 = '#D4A853' },
  @{ id = '232c15af-a295-4153-b0e3-27e747bd69b9'; name = 'OEKO-TEX Standard 100';  c1 = '#3A0CA3'; c2 = '#F72585' },
  @{ id = 'c0627738-e055-458b-af3b-ba86d0c6e6b7'; name = 'SEDEX Registered';       c1 = '#134074'; c2 = '#8DA9C4' }
)
foreach ($c in $certs) {
  $f = "$base\_cert_$($c.id.Substring(0,8)).png"
  New-Cover $f 800 600 $c.c1 $c.c2 $c.name
  $url = Upload-Media $f 'certification' "Certificate - $($c.name)"
  $bodyPath = "$base\_cert_body.json"
  # 无单条 GET 接口，从列表取完整记录整体回传，避免清空其它字段
  $cur = (Invoke-RestMethod -Uri "$api/admin/certifications" -Headers $h).data | Where-Object { $_.id -eq $c.id }
  Set-JsonFile $bodyPath @{
    name        = $cur.name; issuer = $cur.issuer; number = $cur.number
    issue_date  = $cur.issue_date; expiry_date = $cur.expiry_date
    image       = $url; pdf = $cur.pdf; description = $cur.description
    status      = $cur.status
  }
  curl.exe -s -X PUT -H "Authorization: Bearer $tok" -H 'Content-Type: application/json' --data-binary "@$bodyPath" "$api/admin/certifications/$($c.id)" | Out-Null
  Write-Host "cert $($c.name) -> $url"
}

# ============ 2. 博客封面 ============
$blogs = @(
  @{ kw = 'bec94e6e'; title = 'Fabric Technology';   c1 = '#0D1B2A'; c2 = '#4CC9F0' },
  @{ kw = '95d00b9f'; title = 'Custom Design';       c1 = '#1B4332'; c2 = '#B7E4C7' },
  @{ kw = '512a069d'; title = 'Choose Manufacturer'; c1 = '#7B2CBF'; c2 = '#E0AAFF' },
  @{ kw = 'ac6a364c'; title = 'Sustainable 2025';    c1 = '#134074'; c2 = '#D4A853' }
)
$blogList = (Invoke-RestMethod -Uri "$api/admin/blogs?page_size=50" -Headers $h).data.items
foreach ($b in $blogs) {
  $blog = $blogList | Where-Object { $_.id.StartsWith($b.kw) }
  if (-not $blog) { Write-Host "skip blog $($b.kw)"; continue }
  $f = "$base\_blog_$($b.kw).png"
  New-Cover $f 1200 675 $b.c1 $b.c2 $b.title
  $url = Upload-Media $f 'blog' "Blog Cover - $($b.title)"
  $bodyPath = "$base\_blog_body.json"
  Set-JsonFile $bodyPath @{
    slug = $blog.slug; title = $blog.title; excerpt = $blog.excerpt
    category = $blog.category; tags = $blog.tags; author = $blog.author
    cover_image = $url; content = $blog.content; status = $blog.status
  }
  curl.exe -s -X PUT -H "Authorization: Bearer $tok" -H 'Content-Type: application/json' --data-binary "@$bodyPath" "$api/admin/blogs/$($blog.id)" | Out-Null
  Write-Host "blog $($blog.slug) -> $url"
}

Write-Host 'DATA-MAINT-DONE'
