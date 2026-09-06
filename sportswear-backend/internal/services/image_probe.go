package services

import (
	"image"
	_ "image/gif"  // 注册 GIF 解码器
	_ "image/jpeg" // 注册 JPEG 解码器
	_ "image/png"  // 注册 PNG 解码器
	"mime/multipart"

	_ "golang.org/x/image/webp" // 注册 WebP 解码器
)

// ProbeImageDimensions 解析图片宽高（仅解码头部，开销小）。
// 支持 JPEG / PNG / GIF / WebP；返回宽、高，非图片或解析失败返回 ok=false。
// 注：AVIF 目前无 Go 标准解码器，跳过尺寸解析（不影响上传）。
func ProbeImageDimensions(fileHeader *multipart.FileHeader) (int, int, bool) {
	src, err := fileHeader.Open()
	if err != nil {
		return 0, 0, false
	}
	defer src.Close()

	cfg, _, err := image.DecodeConfig(src)
	if err != nil {
		return 0, 0, false
	}
	return cfg.Width, cfg.Height, true
}
