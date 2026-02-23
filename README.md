# logMessageCheck-lint

Custom Go linter that validates log messages for style and security rules. Supports `log/slog` and `go.uber.org/zap`.

## Правила проверки

- Сообщение должно быть только на английском языке
- Сообщение должно начинаться с маленькой буквы
- Запрещены спецсимволы и emoji
- Проверка на утечку чувствительных данных (token, password, secret и т.д.)

---

## Интеграция с golangci-lint

Линтер интегрируется с golangci-lint двумя способами.

### Способ 1: Module Plugin System (рекомендуется)

Современный способ без CGO. Собирает кастомный бинарник golangci-lint с плагином.

#### Требования

- Go 1.21+
- git
- golangci-lint v2.x

#### Установка

**Вариант A: из локального исходного кода (разработка)**

1. Клонируйте репозиторий:

```bash
git clone git@github.com:BahadirAhmedov/logMessageCheck-lint.git
cd logMessageCheck-lint
```

2. Создайте `.custom-gcl.yml` в корне вашего проекта:

```yaml
version: v2.10.1
plugins:
  - module: "github.com/BahadirAhmedov/LogMessageCheck/lint"
    import: "github.com/BahadirAhmedov/LogMessageCheck/lint/plugin"
    path: /путь/к/папке/lint   # абсолютный путь к клону репозитория
```

3. Соберите кастомный бинарник:

```bash
golangci-lint custom -v
```

Создастся файл `custom-gcl` в текущей директории.

**Вариант B: из Go proxy (релизная версия)**

1. Создайте `.custom-gcl.yml`:

```yaml
version: v2.10.1
plugins:
  - module: "github.com/BahadirAhmedov/LogMessageCheck/lint"
    import: "github.com/BahadirAhmedov/LogMessageCheck/lint/plugin"
    version: v1.0.0   # после публикации модуля
```

2. Соберите:

```bash
golangci-lint custom -v
```

#### Конфигурация .golangci.yml

```yaml
version: "2"
linters:
  default: none
  enable:
    - mycustomlinter
  settings:
    custom:
      mycustomlinter:
        type: "module"
        description: "Checks log messages for style and security rules"
        settings: {}
```

#### Использование

```bash
./custom-gcl run ./...
```

---

### Способ 2: Go Plugin System (.so плагин)

Классический способ через shared library. Требует CGO.

#### Требования

- CGO_ENABLED=1
- golangci-lint должен быть собран с той же версией Go и для той же платформы

#### Сборка плагина

```bash
cd /путь/к/LogMessageCheck/lint
go build -buildmode=plugin -tags=plugin -o mycustomlinter.so ./goplugin
```

#### Конфигурация .golangci.yml

```yaml
version: "2"
linters:
  enable:
    - mycustomlinter
  settings:
    custom:
      mycustomlinter:
        path: /путь/к/mycustomlinter.so
        description: "Checks log messages for style and security rules"
```

#### Использование

```bash
golangci-lint run -E mycustomlinter ./...
```

---

## Запуск линтера отдельно (без golangci-lint)

```bash
go run ./cmd/log-message-check ./...
```

---

## Структура проекта

```
.
├── cmd/log-message-check/   # CLI для standalone запуска
├── pkg/analyzer/            # Логика анализатора
├── plugin/                  # Адаптер для Module Plugin System
├── goplugin/                # Адаптер для Go Plugin System (.so)
├── .custom-gcl.yml          # Пример конфигурации для golangci-lint custom
├── .golangci.yml            # Пример конфигурации golangci-lint
└── example.go               # Пример кода для тестирования
```
