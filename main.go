package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver"
)

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		// Получаем файл из формы запроса по ключу "file"
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Создаем директорию для сохранения файлов
		os.MkdirAll("uploads", os.ModePerm)

		// Генерируем полный путь к файлу на сервере
		dst := filepath.Join("uploads", file.Filename)

		// Сохраняем zip-файл на сервере
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Извлекаем файл из архива
		archiver.Unarchive(dst, "arch")

		// Открываем и читаем файл
		archivedFile, err := os.Open("arch/text.txt")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer archivedFile.Close()

		// Читаем содержимое файла
		data, err := io.ReadAll(archivedFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Преобразуем данные в строку
		content := string(data)

		// Выводим содержимое файла в консоль
		fmt.Println(content)

		c.JSON(http.StatusOK, gin.H{"message": "Файл успешно сохранен"})
	})

	r.Run(":8080")
}
