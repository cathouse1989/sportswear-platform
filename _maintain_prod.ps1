$ErrorActionPreference = 'Stop'
$base = 'd:\code\sportswear-platform'
$tok = Get-Content "$base\_tok.txt"
$h = @{ Authorization = "Bearer $tok" }
$api = 'http://localhost:8080/api/v1'

$urls = Get-Content "$base\_prod_urls.txt" | Where-Object { $_.Trim() -ne '' }
if ($urls.Count -lt 4) { throw 'product urls missing' }

$prods = (Invoke-RestMethod -Uri "$api/admin/products?page_size=50" -Headers $h).data.items
$i = 0
$out = @()
foreach ($p in $prods) {
  if ($p.cover_image -notmatch '^https://images\.unsplash\.com') { continue }
  $detail = (Invoke-RestMethod -Uri "$api/admin/products/$($p.id)" -Headers $h).data
  $url = $urls[$i % $urls.Count]; $i++
  $body = [ordered]@{
    sku = $detail.sku; slug = $detail.slug; category_id = $detail.category_id
    type = $detail.type; gender = $detail.gender; status = $detail.status
    is_featured = $detail.is_featured; is_new = $detail.is_new; sort_order = $detail.sort_order
    cover_image = $url
    brief = $detail.brief; description = $detail.description
    features = $detail.features; usage = $detail.usage
    material = $detail.material; composition = $detail.composition; weight = $detail.weight
    elasticity = $detail.elasticity; fit = $detail.fit; support_level = $detail.support_level
    season = $detail.season; size_range = $detail.size_range
    sample_moq = $detail.sample_moq; production_moq = $detail.production_moq
    color_moq = $detail.color_moq; size_moq = $detail.size_moq
    translations = @($detail.translations | ForEach-Object {
      [ordered]@{ language = $_.language; name = $_.name; brief = $_.brief
                  description = $_.description; features = $_.features; usage = $_.usage }
    })
  }
  $bp = "$base\_prod_body.json"
  Set-Content -Path $bp -Value ($body | ConvertTo-Json -Depth 6) -Encoding UTF8
  curl.exe -s -X PUT -H "Authorization: Bearer $tok" -H 'Content-Type: application/json' --data-binary "@$bp" "$api/admin/products/$($p.id)" | Out-Null
  $out += "PROD $($detail.sku) -> $url"
}
$out | Out-File -Encoding utf8 "$base\_prod_result.txt"
Write-Host 'PROD-MAINT-DONE'