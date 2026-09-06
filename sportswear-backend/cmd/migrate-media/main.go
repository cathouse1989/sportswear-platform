// 媒体本地 → MinIO 批量迁移工具
//
// 用法：
//   go run ./cmd/migrate-media          # 实际迁移（上传 + 回写 source=minio）
//   go run ./cmd/migrate-media -dry-run # 仅预览，不实际上传/更新
//   go run ./cmd/migrate-media -only-orphans # 仅扫描 uploads/ 目录中的孤儿文件（无媒体记录对应）
//
// 依赖环境变量：DB_*（连接 media 表）、MINIO_*（对象存储）、STORAGE_DRIVER 无关（本工具固定操作 MinIO）。
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"sportswear-backend/internal/config"
	"sportswear-backend/internal/database"
	"sportswear-backend/internal/models"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "仅预览，不实际上传/更新")
	onlyOrphans := flag.Bool("only-orphans", false, "仅扫描 uploads/ 目录中的孤儿文件（无媒体记录对应）")
	flag.Parse()

	_ = godotenv.Load()
	cfg := config.Load()

	client, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		log.Fatalf("初始化 MinIO 客户端失败: %v", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.MinIO.Bucket)
	if err != nil {
		log.Fatalf("检查桶失败: %v", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.MinIO.Bucket, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("创建桶失败: %v", err)
		}
	}

	if *onlyOrphans {
		scanOrphans()
		return
	}

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	var medias []models.Media
	if err := db.Where("source = ? AND path IS NOT NULL AND path <> ''", models.MediaSourceLocal).Find(&medias).Error; err != nil {
		log.Fatalf("查询媒体失败: %v", err)
	}

	var migrated, skipped, failed int
	for _, m := range medias {
		localPath := filepath.Join("uploads", filepath.FromSlash(m.Path))
		if _, err := os.Stat(localPath); err != nil {
			skipped++
			log.Printf("跳过 %s：本地文件不存在 %s", m.ID, localPath)
			continue
		}

		if *dryRun {
			log.Printf("[dry-run] 将迁移 %s -> minio://%s/%s", localPath, cfg.MinIO.Bucket, m.Path)
			continue
		}

		contentType := mime.TypeByExtension(filepath.Ext(m.Path))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		if _, err := client.FPutObject(ctx, cfg.MinIO.Bucket, m.Path, localPath,
			minio.PutObjectOptions{ContentType: contentType}); err != nil {
			failed++
			log.Printf("上传失败 %s: %v", m.Path, err)
			continue
		}

		if err := db.Model(&models.Media{}).Where("id = ?", m.ID).
			Update("source", models.MediaSourceMinIO).Error; err != nil {
			failed++
			log.Printf("回写 source 失败 %s: %v", m.ID, err)
			continue
		}
		migrated++
	}

	log.Printf("迁移完成：成功 %d，跳过 %d，失败 %d", migrated, skipped, failed)
}

// scanOrphans 扫描 uploads/ 目录下无 media 记录对应的孤儿文件（只读，不删除）
func scanOrphans() {
	cfg := config.Load()
	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 收集 DB 中已记录的 path（相对对象路径）
	known := map[string]bool{}
	var medias []models.Media
	if err := db.Find(&medias).Error; err != nil {
		log.Fatalf("查询媒体失败: %v", err)
	}
	for _, m := range medias {
		if m.Path != "" {
			known[m.Path] = true
		}
	}

	orphanCount := 0
	err = filepath.Walk("uploads", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel := filepath.ToSlash(path)
		rel = rel[len("uploads/"):]
		if !known[rel] {
			orphanCount++
			fmt.Printf("[orphan] %s\n", path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("扫描失败: %v", err)
	}
	log.Printf("孤儿文件扫描完成：共 %d 个", orphanCount)
}
