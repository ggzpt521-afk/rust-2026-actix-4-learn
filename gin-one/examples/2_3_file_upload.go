// ============================================================================
// 2.3 文件上传与下载
// ============================================================================
// 运行方式: go run examples/2_3_file_upload.go
// 测试前先创建目录: mkdir -p uploads
// ============================================================================

package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// 核心概念：HTTP 文件上传原理
// ============================================================================
//
// 【multipart/form-data 格式】
//
// 文件上传使用 multipart/form-data 编码，格式如下：
//
// POST /upload HTTP/1.1
// Content-Type: multipart/form-data; boundary=----WebKitFormBoundary
//
// ------WebKitFormBoundary
// Content-Disposition: form-data; name="file"; filename="test.jpg"
// Content-Type: image/jpeg
//
// [文件二进制内容]
// ------WebKitFormBoundary--
//
// 【Gin 处理流程】
//
// 1. c.FormFile("file") - 获取文件头信息 (不读取内容)
// 2. file.Open() - 打开文件流
// 3. io.Copy() - 流式写入目标文件
//
// 【为什么用流式处理？】
//
// 大文件如果全部读入内存会导致 OOM
// 流式处理: 边读边写，内存占用恒定
//
// ============================================================================

const (
	// 上传目录
	UploadDir = "./uploads"
	// 单文件最大 10MB
	MaxFileSize = 10 << 20 // 10MB
	// 请求体最大 50MB (多文件上传)
	MaxBodySize = 50 << 20 // 50MB
)

// 允许的文件类型
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var AllowedDocTypes = map[string]bool{
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"text/plain": true,
}

// generateID 生成随机 ID (替代 UUID)
func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func main() {
	r := gin.Default()

	// 设置请求体大小限制
	r.MaxMultipartMemory = MaxBodySize

	// 确保上传目录存在
	os.MkdirAll(UploadDir, 0755)

	// ========================================================================
	// 一、单文件上传 (基础版)
	// ========================================================================

	r.POST("/upload/simple", func(c *gin.Context) {
		// 获取文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "file_required",
				"message": "请选择要上传的文件",
			})
			return
		}

		// 保存文件 (使用原文件名，不推荐用于生产环境)
		dst := filepath.Join(UploadDir, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "save_failed",
				"message": "文件保存失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "上传成功",
			"filename": file.Filename,
			"size":     file.Size,
		})
	})

	// ========================================================================
	// 二、单文件上传 (生产级)
	// ========================================================================

	r.POST("/upload/image", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "file_required",
				"message": "请选择要上传的文件",
			})
			return
		}
		defer file.Close()

		// 1. 文件大小校验
		if header.Size > MaxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "file_too_large",
				"message": fmt.Sprintf("文件大小不能超过 %dMB", MaxFileSize/(1<<20)),
			})
			return
		}

		// 2. 文件类型校验 (读取文件头判断真实类型)
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "read_failed",
				"message": "无法读取文件",
			})
			return
		}

		// 检测真实 MIME 类型
		contentType := http.DetectContentType(buffer)
		if !AllowedImageTypes[contentType] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":        "invalid_type",
				"message":      "只允许上传图片文件 (jpg, png, gif, webp)",
				"content_type": contentType,
			})
			return
		}

		// 重置文件指针
		file.Seek(0, 0)

		// 3. 生成安全的文件名
		ext := filepath.Ext(header.Filename)
		if ext == "" {
			// 根据 MIME 类型推断扩展名
			switch contentType {
			case "image/jpeg":
				ext = ".jpg"
			case "image/png":
				ext = ".png"
			case "image/gif":
				ext = ".gif"
			case "image/webp":
				ext = ".webp"
			}
		}

		// 使用 UUID + 时间戳作为文件名
		newFilename := fmt.Sprintf("%s_%d%s",
			generateID(),
			time.Now().Unix(),
			strings.ToLower(ext),
		)

		// 4. 按日期分目录存储
		dateDir := time.Now().Format("2006/01/02")
		fullDir := filepath.Join(UploadDir, "images", dateDir)
		os.MkdirAll(fullDir, 0755)

		dst := filepath.Join(fullDir, newFilename)

		// 5. 流式保存文件
		out, err := os.Create(dst)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "create_failed",
				"message": "无法创建文件",
			})
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "write_failed",
				"message": "文件写入失败",
			})
			return
		}

		// 6. 返回文件访问路径
		relativePath := filepath.Join("images", dateDir, newFilename)

		c.JSON(http.StatusOK, gin.H{
			"message":       "上传成功",
			"original_name": header.Filename,
			"saved_name":    newFilename,
			"size":          header.Size,
			"content_type":  contentType,
			"path":          relativePath,
			"url":           fmt.Sprintf("/files/%s", relativePath),
		})
	})

	// ========================================================================
	// 三、多文件上传
	// ========================================================================

	r.POST("/upload/multiple", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "parse_failed",
				"message": "解析表单失败",
			})
			return
		}

		files := form.File["files"] // 字段名是 files
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "files_required",
				"message": "请选择要上传的文件",
			})
			return
		}

		// 限制上传数量
		if len(files) > 10 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "too_many_files",
				"message": "最多上传 10 个文件",
			})
			return
		}

		var results []gin.H
		var errors []gin.H

		for _, file := range files {
			// 校验单个文件大小
			if file.Size > MaxFileSize {
				errors = append(errors, gin.H{
					"filename": file.Filename,
					"error":    "文件过大",
				})
				continue
			}

			// 生成新文件名
			ext := filepath.Ext(file.Filename)
			newFilename := fmt.Sprintf("%s_%d%s",
				generateID(),
				time.Now().UnixNano(),
				ext,
			)

			dst := filepath.Join(UploadDir, "batch", newFilename)
			os.MkdirAll(filepath.Dir(dst), 0755)

			// 保存文件
			if err := c.SaveUploadedFile(file, dst); err != nil {
				errors = append(errors, gin.H{
					"filename": file.Filename,
					"error":    "保存失败",
				})
				continue
			}

			results = append(results, gin.H{
				"original_name": file.Filename,
				"saved_name":    newFilename,
				"size":          file.Size,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    fmt.Sprintf("成功上传 %d 个文件", len(results)),
			"success":    results,
			"failed":     errors,
			"total":      len(files),
			"successful": len(results),
		})
	})

	// ========================================================================
	// 四、大文件流式上传 (分块读取)
	// ========================================================================

	r.POST("/upload/stream", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no file"})
			return
		}
		defer file.Close()

		// 生成文件路径
		newFilename := fmt.Sprintf("%s_%d%s",
			generateID(),
			time.Now().Unix(),
			filepath.Ext(header.Filename),
		)
		dst := filepath.Join(UploadDir, "large", newFilename)
		os.MkdirAll(filepath.Dir(dst), 0755)

		// 创建目标文件
		out, err := os.Create(dst)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
			return
		}
		defer out.Close()

		// 流式写入 (32KB 缓冲区)
		buffer := make([]byte, 32*1024)
		var totalBytes int64

		for {
			n, err := file.Read(buffer)
			if n > 0 {
				written, writeErr := out.Write(buffer[:n])
				if writeErr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "write failed"})
					return
				}
				totalBytes += int64(written)
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "read failed"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "上传成功",
			"original_name": header.Filename,
			"saved_name":    newFilename,
			"size":          totalBytes,
		})
	})

	// ========================================================================
	// 五、文件下载
	// ========================================================================

	// 直接下载 (小文件)
	r.GET("/download/:filename", func(c *gin.Context) {
		filename := c.Param("filename")

		// 安全检查：防止路径遍历攻击
		filename = filepath.Base(filename)
		filePath := filepath.Join(UploadDir, filename)

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "文件不存在",
			})
			return
		}

		// 设置下载文件名
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.File(filePath)
	})

	// 流式下载 (大文件)
	r.GET("/download/stream/:filename", func(c *gin.Context) {
		filename := filepath.Base(c.Param("filename"))
		filePath := filepath.Join(UploadDir, filename)

		// 打开文件
		file, err := os.Open(filePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		defer file.Close()

		// 获取文件信息
		stat, _ := file.Stat()

		// 设置响应头
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))

		// 流式传输
		c.Stream(func(w io.Writer) bool {
			buffer := make([]byte, 32*1024)
			for {
				n, err := file.Read(buffer)
				if n > 0 {
					w.Write(buffer[:n])
				}
				if err == io.EOF {
					return false
				}
				if err != nil {
					return false
				}
			}
		})
	})

	// ========================================================================
	// 六、静态文件服务 (访问已上传的文件)
	// ========================================================================

	// 静态文件服务
	r.Static("/files", UploadDir)

	// 或者使用 StaticFS 自定义配置
	// r.StaticFS("/assets", http.Dir(UploadDir))

	// ========================================================================
	// 七、带表单数据的文件上传
	// ========================================================================

	type AvatarUploadForm struct {
		UserID      int    `form:"user_id" binding:"required"`
		Description string `form:"description"`
	}

	r.POST("/upload/avatar", func(c *gin.Context) {
		// 先绑定表单数据
		var form AvatarUploadForm
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_form",
				"message": err.Error(),
			})
			return
		}

		// 获取文件
		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "file_required",
				"message": "请选择头像文件",
			})
			return
		}

		// 保存文件
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("avatar_%d%s", form.UserID, ext)
		dst := filepath.Join(UploadDir, "avatars", newFilename)
		os.MkdirAll(filepath.Dir(dst), 0755)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     "头像上传成功",
			"user_id":     form.UserID,
			"description": form.Description,
			"avatar_url":  fmt.Sprintf("/files/avatars/%s", newFilename),
		})
	})

	r.Run(":8080")
}

// ============================================================================
// 测试命令
// ============================================================================
//
// # 先创建上传目录
// mkdir -p uploads
//
// # 创建测试文件
// echo "test content" > test.txt
// # 或使用真实图片文件
//
// # 单文件上传 (简单版)
// curl -X POST http://localhost:8080/upload/simple \
//   -F "file=@test.txt"
//
// # 单文件上传 (生产版，上传图片)
// curl -X POST http://localhost:8080/upload/image \
//   -F "file=@/path/to/image.jpg"
//
// # 多文件上传
// curl -X POST http://localhost:8080/upload/multiple \
//   -F "files=@test.txt" \
//   -F "files=@test2.txt"
//
// # 带表单数据的文件上传
// curl -X POST http://localhost:8080/upload/avatar \
//   -F "user_id=123" \
//   -F "description=我的头像" \
//   -F "avatar=@/path/to/avatar.jpg"
//
// # 文件下载
// curl -O http://localhost:8080/download/test.txt
//
// # 访问静态文件
// curl http://localhost:8080/files/test.txt
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【文件名安全】
//    永远不要直接使用用户上传的文件名
//    - 可能包含路径遍历攻击: ../../../etc/passwd
//    - 可能包含特殊字符导致问题
//    - 可能与已有文件冲突
//    解决: 使用 UUID 或 Hash 生成新文件名
//
// 2. 【文件类型校验】
//    不要仅依赖文件扩展名或 Content-Type
//    - 扩展名可以伪造
//    - Content-Type 由客户端设置，不可信
//    解决: 读取文件头，使用 http.DetectContentType()
//
// 3. 【文件大小限制】
//    必须设置限制，否则可能被大文件攻击
//    - r.MaxMultipartMemory: 内存限制
//    - 在代码中校验 header.Size
//
// 4. 【路径遍历攻击】
//    下载文件时要验证路径
//    filename = filepath.Base(filename) // 只保留文件名
//
// 5. 【大文件处理】
//    大文件应该用流式处理，不要一次性读入内存
//    使用 io.Copy 或分块读取
//
// 6. 【并发写入】
//    多用户同时上传同名文件会冲突
//    使用 UUID 或时间戳确保文件名唯一
//
// 7. 【磁盘空间】
//    生产环境要监控磁盘空间
//    定期清理过期文件
//    考虑使用对象存储 (S3, OSS)
//
// ============================================================================

// ============================================================================
// 生产环境建议
// ============================================================================
//
// 1. 【使用对象存储】
//    - 阿里云 OSS
//    - AWS S3
//    - MinIO (自建)
//    优势: 高可用、CDN 加速、自动扩容
//
// 2. 【文件元数据存储】
//    将文件信息存入数据库:
//    - 原始文件名
//    - 存储路径
//    - 文件大小
//    - MIME 类型
//    - 上传者
//    - 上传时间
//    - MD5/SHA256 (去重)
//
// 3. 【病毒扫描】
//    接入 ClamAV 等病毒扫描服务
//
// 4. 【图片处理】
//    - 自动生成缩略图
//    - 去除 EXIF 信息 (隐私)
//    - 压缩优化
//
// 5. 【断点续传】
//    大文件支持分块上传和断点续传
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现文件去重:
//    - 计算文件 MD5
//    - 相同 MD5 的文件不重复存储
//    - 返回已有文件的路径
//
// 2. 实现图片上传自动缩略图:
//    - 上传图片后自动生成 100x100、300x300 缩略图
//    - 返回原图和缩略图的 URL
//
// 3. 实现文件上传进度:
//    - 使用 SSE 或 WebSocket 推送上传进度
//
// ============================================================================
