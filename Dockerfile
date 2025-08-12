#---------- Build stage ----------
    FROM golang:1.22 AS builder

    # Рабочая директория для сборки
    WORKDIR /app
    
    # Кэшируем зависимости
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Копируем исходный код
    COPY . .
    
    # Собираем бинарник из папки cmd/web
    RUN go build -o myapp ./cmd/web
    
    
    # ---------- Run stage ----------
    FROM debian:bookworm-slim
    
    # Рабочая директория контейнера
    WORKDIR /root/
    
    # Копируем бинарник из builder-образа
    COPY --from=builder /app/myapp .
    
    # Указываем порт, на котором будет слушать приложение
    EXPOSE 8080
    
    # Запускаем приложение
    CMD ["./myapp"]