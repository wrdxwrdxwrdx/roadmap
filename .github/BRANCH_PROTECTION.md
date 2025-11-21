# Настройка защиты веток (Branch Protection)

Для запрета влития PR при неуспешном workflow необходимо настроить Branch Protection Rules в настройках репозитория GitHub.

## Инструкция по настройке

1. Перейдите в настройки репозитория:
   - Settings → Branches → Add rule

2. Для каждой защищаемой ветки (`main`, `develop`) настройте:

### Обязательные настройки:

- ✅ **Require a pull request before merging**
  - Require approvals: 1 (или больше, по необходимости)
  - Dismiss stale pull request approvals when new commits are pushed

- ✅ **Require status checks to pass before merging**
  - Require branches to be up to date before merging
  - Добавьте следующие статусы:
    - `Lint` (из job `lint`)
    - `Test` (из job `test`)
    - `Build` (из job `build`)
    - `Check CI Status` (из job `check-status` в merge-blocker.yml, опционально)

- ✅ **Require conversation resolution before merging** (опционально)

- ✅ **Do not allow bypassing the above settings** (рекомендуется)

### Пример настройки:

```
Branch name pattern: main
  ☑ Require a pull request before merging
    ☑ Require approvals: 1
    ☑ Dismiss stale pull request approvals when new commits are pushed
  ☑ Require status checks to pass before merging
    ☑ Require branches to be up to date before merging
    Status checks:
      ☑ Lint
      ☑ Test
      ☑ Build
  ☑ Do not allow bypassing the above settings
```

## Альтернативный способ (через GitHub CLI)

Если у вас установлен GitHub CLI, можно настроить через команды:

```bash
# Для ветки main
gh api repos/:owner/:repo/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["Lint","Test","Build"]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"required_approving_review_count":1}' \
  --field restrictions=null

# Для ветки develop
gh api repos/:owner/:repo/branches/develop/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["Lint","Test","Build"]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"required_approving_review_count":1}' \
  --field restrictions=null
```

## Проверка работы

После настройки:
1. Создайте PR с намеренной ошибкой (например, неправильное форматирование)
2. Убедитесь, что кнопка "Merge" заблокирована
3. Исправьте ошибку
4. После успешного прохождения всех проверок кнопка "Merge" должна стать доступной

## Важно

- Настройки защиты веток применяются только для указанных веток
- Администраторы репозитория могут обойти защиту (если не включена опция "Do not allow bypassing")
- Статусы проверок должны совпадать с именами jobs в workflow файлах
- После настройки Branch Protection Rules, кнопка "Merge" будет заблокирована до тех пор, пока все проверки не пройдут успешно

## Дополнительный workflow

Файл `.github/workflows/merge-blocker.yml` содержит дополнительную проверку статусов CI.
Он не является обязательным, так как основная защита обеспечивается через Branch Protection Rules.
Однако он может быть полезен для дополнительной валидации перед merge.

## Проверка работы защиты

После настройки Branch Protection Rules:

1. Создайте тестовый PR с намеренной ошибкой (например, неправильное форматирование кода)
2. Убедитесь, что:
   - Кнопка "Merge" заблокирована (серая)
   - Отображается сообщение: "Merging is blocked: Required status check 'Lint' is expected"
3. Исправьте ошибку и дождитесь успешного прохождения всех проверок
4. После успешного прохождения всех проверок кнопка "Merge" должна стать доступной

