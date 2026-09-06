package services

import (
	"bytes"
	"image"
	_ "image/gif"  // 注册 GIF 解码器
	_ "image/jpeg" // 注册 JPEG 解码器
	"image/jpeg"
	_ "image/png" // 注册 PNG 解码器
	"io"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // 注册 WebP 解码器
)

// ThumbnailMaxWidth 缩略图最长边像素（宽度）
const ThumbnailMaxWidth = 320

// GenerateThumbnail 生成缩略图（统一输出 JPEG，最长边不超过 maxWidth，质量 82）。
// 支持 JPEG / PNG / GIF / WebP 源图；非图片或解码失败返回 nil, error。
// 注：AVIF 无 Go 标准解码器，跳过缩略图生成（不影响上传）。
func GenerateThumbnail(reader io.Reader, maxWidth int) ([]byte, error) {
	if maxWidth <= 0 {
		maxWidth = ThumbnailMaxWidth
	}

	src, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w > maxWidth && w > 0 {
		nh := h * maxWidth / w
		if nh < 1 {
			nh = 1
		}
		dst := image.NewRGBA(image.Rect(0, 0, maxWidth, nh))
		draw.CatmullRom.Scale(dst, dst.Bounds(), src, b, draw.Over, nil)
		src = dst
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, src, &jpeg.Options{Quality: 82}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ThumbPath 根据原对象路径推导缩略图路径（统一 .jpg 后缀）
// 例如 products/2026-08/uuid.png -> products/2026-08/uuid_thumb.jpg
func ThumbPath(storagePath string) string {
	ext := filepath.Ext(storagePath)
	return strings.TrimSuffix(storagePath, ext) + "_thumb.jpg"
}
